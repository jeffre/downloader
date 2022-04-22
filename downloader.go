package download

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"sync"
)

type Downloader struct {
	Threads     int
	destDir     string
	jobsChan    chan job
	resultsChan chan result
	done        chan bool
	jobs        []job
}

type job struct {
	url      string
	filename string
}

type result struct {
	job
	size int64
	err  error
}

var ErrDuplicateFilename = errors.New("duplicate filename")

// New provides a *Downloader which defaults to using 3 threads for downloading
// and the current directory as DestDir.
func New() *Downloader {
	d := &Downloader{}
	d.Threads = 3
	d.DestDir(".")
	return d
}

func (d *Downloader) queueJobs() {
	for _, job := range d.jobs {
		d.jobsChan <- job
	}
	close(d.jobsChan)
}

func (d *Downloader) receiveResults() {
	for result := range d.resultsChan {
		if result.err != nil {
			fmt.Fprintf(os.Stderr, "Error downloading from %q: %q\n", result.url, result.err)
			continue
		}

		fmt.Printf("Downloaded %v from %q\n", byteCountSI(result.size), result.url)
	}
	d.done <- true
}

func (d *Downloader) startWorkers(workers int, jobs chan job) {
	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go d.worker(&wg, jobs, d.resultsChan)
	}
	wg.Wait()
	close(d.resultsChan)
}

func (d *Downloader) worker(wg *sync.WaitGroup, jobs chan job, results chan result) {
	for job := range jobs {
		size, err := d.download(job.url, job.filename)
		output := result{job, size, err}
		results <- output
	}
	wg.Done()
}

func (d *Downloader) download(url, filename string) (int64, error) {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	// Fail if response was not 200
	if resp.StatusCode != http.StatusOK {
		msg := fmt.Sprintf("got http status code: %v", resp.StatusCode)
		return 0, errors.New(msg)
	}

	// Create the file
	fullPath := path.Join(d.destDir, filename)
	file, err := os.Create(fullPath)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	// Write the body to file
	bytes, err := io.Copy(file, resp.Body)
	return bytes, err
}

// Run starts the concurrent downloads and returns once all downloads are
// finished.
func (d *Downloader) Run() {

	d.jobsChan = make(chan job, d.Threads)
	d.resultsChan = make(chan result, d.Threads)
	d.done = make(chan bool)

	// Queue jobs
	go d.queueJobs()

	// Queue reception of job results
	go d.receiveResults()

	// Start workers
	go d.startWorkers(d.Threads, d.jobsChan)

	// Wait for workers to finish
	<-d.done
}

// Add queues a url to be downloaded to filename. It can be called multiple
// times. An error (ErrDuplicateFilename) is returned if the same filename is
// provided more than once.
func (d *Downloader) Add(url, filename string) error {
	for _, job := range d.jobs {
		if job.filename == filename {
			return ErrDuplicateFilename
		}
	}
	d.jobs = append(d.jobs, job{url, filename})
	return nil
}

// DestDir sets the dir where downloads will be written. If dir does not exist
// it will return a PathError.
func (d *Downloader) DestDir(dir string) error {
	if _, err := os.Stat(dir); err != nil {
		return err
	}

	d.destDir = dir
	return nil
}

func byteCountSI(b int64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB",
		float64(b)/float64(div), "kMGTPE"[exp])
}
