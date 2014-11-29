package mktorrent

import (
	"crypto/sha1"
	"github.com/zeebo/bencode"
	"io"
	"time"
)

const piece_len = 512000

type InfoDict struct {
	Name        string `bencode:"name"`
	Length      int    `bencode:"length"`
	PieceLength int    `bencode:"piece length"`
	Pieces      string `bencode:"pieces"`
}

type Torrent struct {
	Info         InfoDict   `bencode:"info"`
	AnnounceList [][]string `bencode:"announce-list,omitempty"`
	Announce     string     `bencode:"announce,omitempty"`
	CreationDate int64      `bencode:"creation date,omitempty"`
	Comment      string     `bencode:"comment,omitempty"`
	CreatedBy    string     `bencode:"created by,omitempty"`
	UrlList      string     `bencode:"url-list,omitempty"`
}

func (t *Torrent) Save(w io.Writer) error {
	enc := bencode.NewEncoder(w)
	return enc.Encode(t)
}

func hashPiece(b []byte) []byte {
	h := sha1.New()
	h.Write(b)
	return h.Sum(nil)
}
func MakeTorrent(r io.Reader, name string, url string, ann ...string) (*Torrent, error) {
	t := &Torrent{
		AnnounceList: make([][]string, 0),
		CreationDate: time.Now().Unix(),
		CreatedBy:    "mktorrent.go",
		Info: InfoDict{
			Name:        name,
			PieceLength: piece_len,
		},
		UrlList: url,
	}
	// the outer list is tiers
	for _, a := range ann {
		t.AnnounceList = append(t.AnnounceList, []string{a})
	}

	b := make([]byte, piece_len)
	for {
		n, err := io.ReadFull(r, b)
		if err != nil && err != io.ErrUnexpectedEOF {
			return nil, err
		}
		if err == io.ErrUnexpectedEOF {
			b = b[:n]
			t.Info.Pieces += string(hashPiece(b))
			t.Info.Length += n
			break
		} else if n == piece_len {
			t.Info.Pieces += string(hashPiece(b))
			t.Info.Length += n
		} else {
			panic("short read!")
		}
	}
	return t, nil
}
