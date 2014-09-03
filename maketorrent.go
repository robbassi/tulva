package main

import (
	"log"
	"os"

	"github.com/jackpal/bencode-go"
)

func MakeTorrentFile(filename string) *os.File {
	file, err := os.Open(filename)
	defer file.Close()
	checkError(err)

	info := &Info{
		PieceLength: 1,
	}
	metaInfo := &MetaInfo{
		Info:         info,
		Announce:     "announce",
		AnnounceList: [][]string{[]string{"announcelist1"}},
		CreationDate: 1,
		Comment:      "test",
		CreatedBy:    "Ramsey",
		Encoding:     "utf-8",
	}

	torrentFile, err := os.Create(filename + ".torrent")
	bencode := bencode.Marshal(torrentFile, metaInfo)
	log.Println(bencode)
	defer torrentFile.Close()
	checkError(err)

	return torrentFile
}
