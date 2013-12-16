// Copyright 2013 Jari Takkala and Brian Dignan. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"log"
	"time"
)

type StatsCollector struct {
	Left          int
	Uploaded      int
	Downloaded    int
	Errors        int
	lastUploads   []int
	lastDownloads []int
	peerCh        chan PeerStats
	trackerCh     chan CurrentStats
	dashboardCh   chan CurrentStats
	ticker        <-chan time.Time
}

type CurrentStats struct {
	UploadRate   float64
	DownloadRate float64
	Uploaded     int
	Downloaded   int
}

func NewStats() *StatsCollector {
	return &StatsCollector{
		peerCh:        make(chan PeerStats),
		ticker:        make(chan time.Time),
		lastUploads:   make([]int, 5),
		lastDownloads: make([]int, 5),
		trackerCh:     make(chan CurrentStats),
		dashboardCh:   make(chan CurrentStats),
	}
}

func (s *StatsCollector) Run() {
	log.Println("StatsCollector : Run : Started")
	defer log.Println("StatsCollector : Run : Stopped")

	s.ticker = time.Tick(time.Second * 1)
	i := 0

	for {
		select {
		case stat := <-s.peerCh:
			s.Downloaded += stat.read
			s.Uploaded += stat.write
			s.Errors += stat.errors
		case <-s.ticker:
			s.lastDownloads[i] = s.Downloaded
			s.lastUploads[i] = s.Uploaded
			total := 0
			for _, octets := range s.lastDownloads {
				total += octets
			}
			downloadRate := float64(total) / 5
			total = 0
			for _, octets := range s.lastUploads {
				total += octets
			}
			uploadRate := float64(total) / 5
			currentStats := CurrentStats{DownloadRate: downloadRate, UploadRate: uploadRate, Downloaded: s.Downloaded, Uploaded: s.Uploaded}
			go func() { s.trackerCh <- currentStats }()
			go func() { s.dashboardCh <- currentStats }()
			i++
			if i%5 == 0 {
				i = 0
			}
		}
	}
}
