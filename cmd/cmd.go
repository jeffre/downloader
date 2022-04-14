package main

import (
	"github.com/jeffre/download"
)

func main() {
	d := download.New()
	d.Urls = []string{"cmd"}
	d.Start()
}
