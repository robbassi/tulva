// Copyright 2013 Jari Takkala and Brian Dignan. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"code.google.com/p/go.net/websocket"
	"log"
	"fmt"
	"net/http"
)

type ChannelUpdate struct {
	Channel interface{}
}

type AddNodes struct {
	AddNodes []string
}

func AddNodeMessage(node string) *GraphStateChange {
	gsc := new(GraphStateChange)
	gsc.Operation = AddNode
	gsc.Node = node
	return gsc
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

type RemoveEdges struct {
	RemoveEdges []Edge
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

// Used to store the entire state of the graph
type DirectedGraph struct {
	Nodes map[string]struct{}
	Edges map[string]Edge  // edge-id to edge, where the edge-id is SourceNode-TargetNode-Name
}

func NewDirectedGraph() *DirectedGraph {
	dg := new(DirectedGraph)
	dg.Nodes = make(map[string]struct{})
	dg.Edges = make(map[string]Edge)
	return dg
}

const (
	AddNode int = iota
	RemoveNode 
	AddEdge
	UpdateEdge
	RemoveEdge
)

type GraphStateChange struct {
	Operation int // Either AddNode, RemoveNode, AddEdge, UpdateEdge or RemoveEdge
	Node    string // Is nil is Edge is set
	Edge 	Edge // Is nil is Node is set
}

type Dashboard struct {
	pieces        []ReceivedPiece
	totalPieces   int
	pieceChan     chan chan ReceivedPiece
	statsCh       chan CurrentStats
	websocketChan chan *websocket.Conn
	websockets    map[string]*websocket.Conn
	finishedPieces FinishedPieces
	graphCh    	  chan GraphStateChange
	directedGraph DirectedGraph
}

func NewDashboard(statsCh chan CurrentStats, graphCh chan GraphStateChange, totalPieces int) *Dashboard {
	return &Dashboard{
		pieces:        make([]ReceivedPiece, 0),
		pieceChan:     make(chan chan ReceivedPiece),
		websocketChan: make(chan *websocket.Conn),
		websockets:    make(map[string]*websocket.Conn),
		statsCh:       statsCh,
		graphCh: 	   graphCh,
		totalPieces:   totalPieces,
		directedGraph: *NewDirectedGraph(),
	}
}

func (ds *Dashboard) wsHandler(ws *websocket.Conn) {
	log.Printf("New websocket connection: %#v", ws.Config)
	// init and send the list of pieces

	ds.websocketChan <- ws
	// FIXME: Do something else here?
	for {
		select {}
	}
}

func (ds *Dashboard) removeEdge(edgeId string) {

	re := new(RemoveEdges)
	re.RemoveEdges = []Edge{ds.directedGraph.Edges[edgeId]}

	cu := new(ChannelUpdate)
	cu.Channel = re

	for _, ws := range ds.websockets {
		websocket.JSON.Send(ws, cu)
	}

	delete(ds.directedGraph.Edges, edgeId)
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
		case innerChan := <-ds.pieceChan:

			justFinished := make([]ReceivedPiece,0)
			for piece := range innerChan {
				justFinished = append(justFinished, piece)
			}
			log.Println("Dashboard : Run : Received %d finished pieces from controller.", len(justFinished))

			ds.finishedPieces.FinishedPieces = append(ds.finishedPieces.FinishedPieces, justFinished...)
			pieceUpdate := &PieceUpdate{Piece: &FinishedPieces{FinishedPieces: justFinished}}
			// tell websockets we have a piece
			for _, ws := range ds.websockets {
				websocket.JSON.Send(ws, pieceUpdate)
			}
		case ws := <-ds.websocketChan:
			ds.websockets[ws.Request().RemoteAddr] = ws
			var totalPieces TotalPieces
			totalPieces.TotalPieces = ds.totalPieces
			pieceTotal := &PieceUpdate{Piece: totalPieces}
			websocket.JSON.Send(ws, pieceTotal)

			pieceUpdate := &PieceUpdate{Piece: ds.finishedPieces}
			websocket.JSON.Send(ws, pieceUpdate)

		case stats := <-ds.statsCh:
			statsUpdate := &StatsUpdate{Stats: stats}
			for _, ws := range ds.websockets {
				websocket.JSON.Send(ws, statsUpdate)
			}

		case graphChange := <-ds.graphCh:

			switch graphChange.Operation {
			case AddNode:
				node := graphChange.Node
				_, exists := ds.directedGraph.Nodes[node]
				if exists {
					log.Fatalf("Dashboard : Run : Attempted to add node %s, but it already exists in the graph", node)
				} else {
					log.Printf("Dashboard : Run : Adding node %s to the graph", node)
				}

				an := new(AddNodes)
				an.AddNodes = []string{node}

				cu := new(ChannelUpdate)
				cu.Channel = an

				for _, ws := range ds.websockets {
					websocket.JSON.Send(ws, cu)
				}

				ds.directedGraph.Nodes[node] = struct{}{}


			case RemoveNode:
				node := graphChange.Node
				_, exists := ds.directedGraph.Nodes[node]
				if !exists {
					log.Fatalf("Dashboard : Run : Attempted to remove node %s, but it doesn't exist in the graph", node)
				} else {
					log.Printf("Dashboard : Run : Removing node %s from the graph", node)
				}

				// Check to see if there are any edges linked to this graph. 
				for edgeId, edge := range ds.directedGraph.Edges {
					if edge.SourceNode == node || edge.TargetNode == node {
						log.Printf("Dashboard : Run : Edge %s must be removed because it depends on %s", edgeId, node)
						ds.removeEdge(edgeId)
					}
				}

				rn := new(RemoveNodes)
				rn.RemoveNodes = []string{node}

				cu := new(ChannelUpdate)
				cu.Channel = rn

				for _, ws := range ds.websockets {
					websocket.JSON.Send(ws, cu)
				}

				delete(ds.directedGraph.Nodes, node)


			case AddEdge:
				edge := graphChange.Edge
				edgeId := fmt.Sprintf("%s-%s-%s",edge.SourceNode,edge.TargetNode,edge.Name) 

				_, exists := ds.directedGraph.Edges[edgeId]
				if exists {
					log.Fatalf("Dashboard : Run : Attempted to add edge %s, but it already exists in the graph", edgeId)
				} else {
					log.Printf("Dashboard : Run : Adding edge %s to the graph", edgeId)
				}

				ae := new(AddEdges)
				ae.AddEdges = []Edge{edge}

				cu := new(ChannelUpdate)
				cu.Channel = ae

				for _, ws := range ds.websockets {
					websocket.JSON.Send(ws, cu)
				}

				ds.directedGraph.Edges[edgeId] = edge


			case UpdateEdge:
				edge := graphChange.Edge
				edgeId := fmt.Sprintf("%s-%s-%s",edge.SourceNode,edge.TargetNode,edge.Name) 
				_, exists := ds.directedGraph.Edges[edgeId]
				if !exists {
					log.Fatalf("Dashboard : Run : Attempted to update edge %s, but it doesn't exist in the graph", edgeId)
				} else {
					log.Printf("Dashboard : Run : Updating edge %s. Changing intensity from %d to %d", edgeId, ds.directedGraph.Edges[edgeId].Intensity, edge.Intensity)
				}	

				existingEdge := ds.directedGraph.Edges[edgeId]
				existingEdge.Intensity = edge.Intensity

				ue := new(UpdateEdges)
				ue.UpdateEdges = []Edge{existingEdge}

				cu := new(ChannelUpdate)
				cu.Channel = ue

				for _, ws := range ds.websockets {
					websocket.JSON.Send(ws, cu)
				}


			case RemoveEdge:
				edge := graphChange.Edge
				edgeId := fmt.Sprintf("%s-%s-%s",edge.SourceNode,edge.TargetNode,edge.Name) 

				_, exists := ds.directedGraph.Edges[edgeId]
				if !exists {
					log.Fatalf("Dashboard : Run : Attempted to remove edge %s, but it doesn't exist in the graph", edgeId)
				} else {
					log.Printf("Dashboard : Run : Removing edge %s from the graph", edgeId)
				}		

				ds.removeEdge(edgeId)
			}
		}
	}
}
