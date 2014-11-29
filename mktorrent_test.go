package mktorrent

import (
	"bytes"
	"os"
	"testing"
)

func TestMktorrent(t *testing.T) {
	b := bytes.NewBufferString("test")
	tor, err := MakeTorrent(b, "1.txt", "udp://tracker.openbittorrent.com:80/announce", "")
	if err != nil {
		t.Fatal(err)
	}
	if tor.Info.Name != "1.txt" {
		t.Fatal("did not get right name")
	}
	if tor.Info.Pieces == "" {
		t.Fatal("did not hash correctly")
	}
	f, err := os.Create("test.torrent")
	if err != nil {
		t.Fatal(err)
	}
	tor.Save(f)
}
