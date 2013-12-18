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

type DiGraphUpdate struct {
	DiGraph interface{}
}

type AddNodes struct {
	AddNodes []Node
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

type Node struct {
	NodeName string
	NodeID string
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
	Nodes map[string]Node  // nodeName to node
	Edges map[string]Edge  // edge-id to edge, where the edge-id is SourceNode-TargetNode-Name
}

func NewDirectedGraph() *DirectedGraph {
	dg := new(DirectedGraph)
	dg.Nodes = make(map[string]Node)
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
	Node    Node // Is nil is Edge is set
	Edge 	Edge // Is nil is Node is set
}

func AddNodeMessage(nodeName string, nodeID string) GraphStateChange {
	gsc := new(GraphStateChange)
	gsc.Operation = AddNode
	gsc.Node = Node{NodeName: nodeName, NodeID: nodeID}
	return *gsc
}

func RemoveNodeMessage(nodeName string) GraphStateChange {
	gsc := new(GraphStateChange)
	gsc.Operation = RemoveNode
	gsc.Node = Node{NodeName: nodeName}
	return *gsc
}

func AddEdgeMessage(sourceNode string, targetNode string, name string, intensity int) GraphStateChange {
	gsc := new(GraphStateChange)
	gsc.Operation = AddEdge
	gsc.Edge = Edge{
		SourceNode: sourceNode,
		TargetNode: targetNode,
		Name: name,
		Intensity: intensity,
	}
	return *gsc
}

func UpdateEdgeMessage(sourceNode string, targetNode string, name string, intensity int) GraphStateChange {
	if intensity < 1 || intensity > 100 {
		log.Fatalf("Edge intensity must be between 1 and 100")
	}
	gsc := new(GraphStateChange)
	gsc.Operation = UpdateEdge
	gsc.Edge = Edge{
		SourceNode: sourceNode,
		TargetNode: targetNode,
		Name: name,
		Intensity: intensity,
	}
	return *gsc
}

func RemoveEdgeMessage(sourceNode string, targetNode string, name string) GraphStateChange {
	gsc := new(GraphStateChange)
	gsc.Operation = RemoveEdge
	gsc.Edge = Edge{
		SourceNode: sourceNode,
		TargetNode: targetNode,
		Name: name,
	}
	return *gsc
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

	gu := new(DiGraphUpdate)
	gu.DiGraph = re

	for _, ws := range ds.websockets {
		websocket.JSON.Send(ws, gu)
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
			log.Printf("Dashboard : Run : Received %d finished pieces from controller.", len(justFinished))

			if (len(justFinished)) > 0 {
				ds.finishedPieces.FinishedPieces = append(ds.finishedPieces.FinishedPieces, justFinished...)
				pieceUpdate := &PieceUpdate{Piece: &FinishedPieces{FinishedPieces: justFinished}}
				// tell websockets we have a piece
				for _, ws := range ds.websockets {
					websocket.JSON.Send(ws, pieceUpdate)
				}
			}
		case ws := <-ds.websocketChan:
			ds.websockets[ws.Request().RemoteAddr] = ws
			var totalPieces TotalPieces
			totalPieces.TotalPieces = ds.totalPieces
			pieceTotal := &PieceUpdate{Piece: totalPieces}
			websocket.JSON.Send(ws, pieceTotal)

			if len(ds.finishedPieces.FinishedPieces) > 0 {
				pieceUpdate := &PieceUpdate{Piece: ds.finishedPieces}
				websocket.JSON.Send(ws, pieceUpdate)
			}

			nodes := make([]Node, 0)
			for _, node := range ds.directedGraph.Nodes {
				nodes = append(nodes, node)
			}

			nodeUpdate := DiGraphUpdate{DiGraph: AddNodes{AddNodes:nodes}}
			websocket.JSON.Send(ws, nodeUpdate)

			edges := make([]Edge, 0)
			for _, edge := range ds.directedGraph.Edges {
				edges = append(edges, edge)
			} 

			edgeUpdate := DiGraphUpdate{DiGraph: AddEdges{AddEdges: edges}}
			websocket.JSON.Send(ws, edgeUpdate)


		case stats := <-ds.statsCh:
			statsUpdate := &StatsUpdate{Stats: stats}
			for _, ws := range ds.websockets {
				websocket.JSON.Send(ws, statsUpdate)
			}

		case graphChange := <-ds.graphCh:

			switch graphChange.Operation {
			case AddNode:
				nodeName := graphChange.Node.NodeName
				_, exists := ds.directedGraph.Nodes[nodeName]
				if exists {
					log.Fatalf("Dashboard : Run : Attempted to add node %s, but it already exists in the graph", nodeName)
				} else {
					log.Printf("Dashboard : Run : Adding node %s to the graph", nodeName)
				}

				an := new(AddNodes)
				an.AddNodes = []Node{graphChange.Node}

				gu := new(DiGraphUpdate)
				gu.DiGraph = an

				for _, ws := range ds.websockets {
					websocket.JSON.Send(ws, gu)
				}

				ds.directedGraph.Nodes[nodeName] = graphChange.Node


			case RemoveNode:
				nodeName := graphChange.Node.NodeName
				_, exists := ds.directedGraph.Nodes[nodeName]
				if !exists {
					log.Fatalf("Dashboard : Run : Attempted to remove node %s, but it doesn't exist in the graph", nodeName)
				} else {
					log.Printf("Dashboard : Run : Removing node %s from the graph", nodeName)
				}

				// Check to see if there are any edges linked to this graph. 
				for edgeId, edge := range ds.directedGraph.Edges {
					if edge.SourceNode == nodeName || edge.TargetNode == nodeName {
						log.Printf("Dashboard : Run : Edge %s must be removed because it depends on %s", edgeId, nodeName)
						ds.removeEdge(edgeId)
					}
				}

				rn := new(RemoveNodes)
				rn.RemoveNodes = []string{nodeName}

				gu := new(DiGraphUpdate)
				gu.DiGraph = rn

				for _, ws := range ds.websockets {
					websocket.JSON.Send(ws, gu)
				}

				delete(ds.directedGraph.Nodes, nodeName)


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

				gu := new(DiGraphUpdate)
				gu.DiGraph = ae

				for _, ws := range ds.websockets {
					websocket.JSON.Send(ws, gu)
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

				gu := new(DiGraphUpdate)
				gu.DiGraph = ue

				for _, ws := range ds.websockets {
					websocket.JSON.Send(ws, gu)
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
