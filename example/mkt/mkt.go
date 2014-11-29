package main

import (
	"flag"
	"fmt"
	"github.com/polvi/mktorrent"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

func main() {
	var u = flag.String("url", "", "url to make a torrent of")
	flag.Parse()
	if *u == "" {
		flag.Usage()
		os.Exit(1)
	}
	res, err := http.Get(*u)
	defer res.Body.Close()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	ur, err := url.Parse(*u)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	filename := filepath.Base(ur.Path)
	t, err := mktorrent.MakeTorrent(res.Body, filename, *u, "udp://tracker.openbittorrent.com:80/announce", "udp://tracker.publicbt.com:80")
	if err != nil {
		fmt.Println(err)
	}
	f, err := os.Create(filename + ".torrent")
	defer f.Close()
	if err != nil {
		fmt.Println(err)
	}
	t.Save(f)
	fmt.Println(filename + ".torrent created")
}
