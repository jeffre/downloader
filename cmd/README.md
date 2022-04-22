## Installation
1. Install [go](https://go.dev/dl/)
1. Run `go get github.com/jeffre/downloader`

## Usage
Downloader takes a list of urls as arguments and downloads them. The default 
mode is to download up to 3 concurrently at a time. This can be changed using
the `-t` flag.

    downloader -t 4 "http://example.com/1" "http://example.com/2" "http://example.com/3" "http://example.com/4"

In the example above, all 4 links would be start downloading at the same time.
