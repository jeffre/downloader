package main

import (
	"flag"
	"fmt"
	"os"
	"path"

	download "github.com/jeffre/downloader"
)

func main() {

	threadPtr := flag.Int("t", 3, "Maximum number of simulatenous downloads to allow")
	destDir := flag.String("d", ".", "Directory to save all downloads")
	flag.Parse()

	d := download.New()
	d.Threads = *threadPtr
	err := d.DestDir(*destDir)
	if err != nil {
		if pathError, ok := err.(*os.PathError); ok {
			fmt.Fprintf(os.Stderr, "%v: %v\n", pathError.Path, pathError.Err)
		} else {
			fmt.Fprintln(os.Stderr, err)
		}
		os.Exit(3)
	}

	for _, url := range flag.Args() {
		d.Add(url, path.Base(url))
	}
	d.Run()
}
