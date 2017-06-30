/*
Package graph implements a directed graph with node and edge structure
Note that edge query is inefficient if graph has a high degree.
*/

package graph

import (
	"math"
	"time"
)

type Node struct {
	ID        int
	X         float64
	Y         float64
	EdgeEnd   []*Edge
	EdgeStart []*Edge
}

type Edge struct {
	ID       [2]int
	From, To *Node
	Weight   float64
}

type DirectedGraph struct {
	Nodes []*Node
}

// NewDirectedGraph initialises an empty graph
func NewDirectedGraph() *DirectedGraph {
	return &DirectedGraph{
		Nodes: make([]*Node, 0),
	}
}

// HasNode checks if node exists in a graph
func (g *DirectedGraph) HasNode(n *Node) bool {
	return n.ID < len(g.Nodes)
}

// Node returns the corresponding node, given an id
// otherwise returns a nil pointer
func (g *DirectedGraph) Node(id int) *Node {
	if len(g.Nodes) > id {
		return g.Nodes[id]
	}

	return nil
}

// HasEdge checks if edge exists in a graph
func (g *DirectedGraph) Edge(id [2]int) bool {
	for _, e := range g.Nodes[id[0]].EdgeStart {
		if e.ID == id {
			return true
		}
	}

	return false
}

// Edge returns the corresponding edge, given an id
// otherwise returns a nil pointer
func (g *DirectedGraph) Edge(id [2]int) *Edge {
	for _, e := range g.Nodes[id[0]].EdgeStart {
		if e.ID == id {
			return e
		}
	}

	return nil
}

// AddNode adds node n to the graph
func (g *DirectedGraph) AddNode(n *Node) {
	n.ID = len(g.Nodes)
	g.Nodes = append(g.Nodes, n)
}

// AddDirectedEdge adds directed edge e to the graph
func (g *DirectedGraph) AddDirectedEdge(e *Edge) {
	from, to := e.From, e.To

	if from.ID == to.ID {
		panic("Self edge detected.")
	}

	if !g.HasNode(from) {
		panic("Unable to add edge. Node does not exist. ")
		return
	}

	if !g.HasNode(to) {
		panic("Unable to add edge. Node does not exist. ")
		return
	}

	g.Nodes[e.From.ID].EdgeStart = append(g.Nodes[e.From.ID].EdgeStart, e)
	g.Nodes[e.To.ID].EdgeEnd = append(g.Nodes[e.To.ID].EdgeEnd, e)
}

// Weight returns weight of directed edge from u to v
func (g *DirectedGraph) Weight(u, v *Node) (w float64, ok bool) {

	//self-edge
	if u.ID == v.ID {
		return 0, false
	}

	for _, e := range g.Nodes[u.ID].EdgeStart {
		if e.To == v {
			return e.Weight, true
		}
	}

	//no edge
	return 0, false
}

// Dist returns the Euclidean distance between two nodes.
func Dist(u, v *Node) float64 {
	return math.Sqrt(SquaredDist(u, v))
}

// SquaredDist returns the squared Euclidean distance between two nodes.
func SquaredDist(u, v *Node) float64 {
	Xdiff := (u.X - v.X)
	Ydiff := (u.Y - v.Y)
	return (Ydiff*Ydiff + Xdiff*Xdiff)
}

// DistFromEdge returns the perpendicular distance between a node and an edge.
func DistFromEdge(n *Node, e *Edge) float64 {
	return math.Sqrt(SquaredDistFromEdge(n, e))
}

// SquaredDistFromEdge returns the squared perpendicular distance between a node and an edge.
func SquaredDistFromEdge(n *Node, e *Edge) float64 {
	x := e.From.Y
	y := e.From.X
	dx := e.To.Y - x
	dy := e.To.X - y

	if dx != 0 || dy != 0 {

		// numerator: projection distance of point from e.From
		// t = fraction of distance : length of line
		t := ((n.Y-x)*dx + (n.X-y)*dy) / (dx*dx + dy*dy)

		if t > 1 {
			x = e.To.Y
			y = e.To.X
		} else if t > 0 {
			// projection of n to the line
			x += dx * t
			y += dy * t
		}
	}

	dx = n.Y - x
	dy = n.X - y

	return (dx*dx + dy*dy)
}
