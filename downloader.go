package download

import "fmt"

type downloader struct {
	Threads int
	Urls    []string
}

func New() *downloader {
	d := &downloader{}
	d.Threads = 3
	return d
}

func (d *downloader) Start() {

	cUrls := make(chan string, d.Threads)

	for _, url := range d.Urls {
		go func(url string) {
			fmt.Println("added", url)
			cUrls <- url
		}(url)
	}

	for i := 0; i < len(d.Urls); i++ {
		url := <-cUrls
		fmt.Printf("got %q\n", url)
	}
}
