// Copyright 2013 Jari Takkala. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"log"
	"math/rand"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Unique client ID, encoded as '-' + 'TV' + <version number> + random digits
var PeerID = [20]byte{'-', 'T', 'V', '0', '0', '0', '1'}

// init initializes a random PeerID for this client
func init() {
	// Initialize PeerID
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 7; i < 20; i++ {
		PeerID[i] = byte(r.Intn(256))
	}
}

func main() {
	switch {
	case len(os.Args) == 2:
		quit := make(chan struct{})
		t, err := NewTorrent(os.Args[1], quit)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("main : main : Started")
		defer log.Println("main : main : Exiting")

		go func() {
			log.Println(http.ListenAndServe("localhost:6060", nil))
		}()

		// Signal handler to catch Ctrl-C and SIGTERM from 'kill' command
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		go func() {
			// Block waiting for a signal
			<-c
			// Unregister the signal handler
			signal.Stop(c)
			log.Println("Received Interrupt. Shutting down...")
			close(t.quit)
		}()

		// Launch the torrent
		t.Run()

	case len(os.Args) == 3 && os.Args[1] == "maketorrentfile":
		t := MakeTorrentFile(os.Args[2])
		log.Println("Made torrent File:", t.Name())
	default:
		log.Fatalf("Usage: %s: <torrent file>\n", os.Args[0])
	}
}
