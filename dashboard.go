// Copyright 2013 Jari Takkala and Brian Dignan. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"code.google.com/p/go.net/websocket"
	"log"
	"net/http"
)

type ChannelUpdate struct {
	Channel interface{}
}

type AddNodes struct {
	AddNodes []string
}

type RemoveNodes struct {
	RemoveNodes []string
}

type AddEdges struct {
	AddEdges []Edge
}

type UpdateEdges struct {
	UpdateEdges []Edge
}

type Edge struct {
	SourceNode string
	TargetNode string
	Name string
	Intensity int  // an intensity scale from 1 to 100
}

type PieceUpdate struct {
	Piece interface{}
}

type TotalPieces struct {
	TotalPieces int
}

type FinishedPieces struct {
	FinishedPieces []ReceivedPiece
}

type StatsUpdate struct {
	Stats interface{}
}

type Dashboard struct {
	pieces        []ReceivedPiece
	totalPieces   int
	pieceChan     chan ReceivedPiece
	statsCh       chan CurrentStats
	websocketChan chan *websocket.Conn
	websockets    map[string]*websocket.Conn
}

func NewDashboard(statsCh chan CurrentStats) *Dashboard {
	return &Dashboard{
		pieces:        make([]ReceivedPiece, 0),
		pieceChan:     make(chan ReceivedPiece),
		websocketChan: make(chan *websocket.Conn),
		websockets:    make(map[string]*websocket.Conn),
		statsCh:       statsCh,
	}
}

func (ds *Dashboard) wsHandler(ws *websocket.Conn) {
	log.Printf("New websocket connection: %#v", ws.Config)
	// init and send the list of pieces

	var totalPieces TotalPieces
	totalPieces.TotalPieces = ds.totalPieces
	pieceUpdate := &PieceUpdate{Piece: totalPieces}
	websocket.JSON.Send(ws, pieceUpdate)

	ds.websocketChan <- ws
	// FIXME: Do something else here?
	for {
		select {}
	}
}

func (ds *Dashboard) Run() {
	log.Println("Dashboard : Run : Started")
	defer log.Println("Dashboard : Run : Stopped")

	http.Handle("/", http.FileServer(http.Dir("content")))
	http.Handle("/ws", websocket.Handler(ds.wsHandler))
	go func() {
		err := http.ListenAndServe(":8080", nil)
		if err != nil {
			log.Fatal("ListenAndServe:", err)
		}
	}()
	for {
		select {
		case piece := <-ds.pieceChan:
			var finishedPieces FinishedPieces
			finishedPieces.FinishedPieces = append(finishedPieces.FinishedPieces, piece)
			pieceUpdate := &PieceUpdate{Piece: finishedPieces}
			// tell websockets we have a piece
			for _, ws := range ds.websockets {
				go websocket.JSON.Send(ws, pieceUpdate)
			}
		case ws := <-ds.websocketChan:
			ds.websockets[ws.Request().RemoteAddr] = ws
		case stats := <-ds.statsCh:
			statsUpdate := &StatsUpdate{Stats: stats}
			for _, ws := range ds.websockets {
				go websocket.JSON.Send(ws, statsUpdate)
			}
		}
	}
}
