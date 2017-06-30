package graph

import (
	"container/heap"
	"math"
)

type distanceNode struct {
	node      *Node
	dist      float64
	realDist  float64
	direction bool // true: forward, false: reverse
}

// priorityQueue implementation, fulfills heap interface
type priorityQueue []*distanceNode

func (q priorityQueue) Len() int { return len(q) }

func (q priorityQueue) Less(i, j int) bool {
	return q[i].dist < q[j].dist
}

func (q priorityQueue) Swap(i, j int) {
	q[i], q[j] = q[j], q[i]
}

func (q *priorityQueue) Push(n interface{}) {
	*q = append(*q, n.(*distanceNode))
}

func (q *priorityQueue) Pop() interface{} {
	node := (*q)[len(*q)-1]
	*q = (*q)[:len(*q)-1]
	return node
}

// Dijkstra returns a a shortest path from u to v and the distance
// in graph g.
func (g *DirectedGraph) Dijkstra(u, v *Node) ([]*Node, float64) {

	forwardDist := make(map[*Node]float64)
	forwardDist[u] = 0.0
	next := make(map[*Node]*Node)

	Q := priorityQueue{{node: u, dist: 0}}
	var mid *distanceNode
	heap.Init(&Q)

	for len(Q) > 0 {

		mid = heap.Pop(&Q).(*distanceNode)

		// terminates when final node is found
		if mid.node.ID == v.ID {
			break
		}

		for _, e := range mid.node.EdgeStart {
			n := e.To

			// total distance travelled so far
			acc_dist := forwardDist[mid.node] + e.Weight

			// update shortest paths
			if dist, ok := forwardDist[n]; !ok || acc_dist < dist {
				heap.Push(&Q, &distanceNode{node: n, dist: acc_dist})
				forwardDist[n] = acc_dist
				next[n] = mid.node
			}

		}
	}

	// no path found
	if mid == nil || mid.node.ID != v.ID {
		return nil, math.Inf(1)
	}

	// retrieve path from results
	n := v
	path := []*Node{v}

	for n != u {
		n = next[n]
		path = append(path, n)
	}

	// reverse the path
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}

	return path, forwardDist[v]
}

// DijkstraBi returns a shortest path from u to v
// in the graph g. Bidirectional variant of dijkstra
func (g *DirectedGraph) DijkstraBi(u, v *Node) ([]*Node, float64) {

	forwardDist := make(map[*Node]float64)
	backwardDist := make(map[*Node]float64)
	next := make(map[*Node]*Node)
	back := make(map[*Node]*Node)

	forwardDist[u] = 0.0
	backwardDist[v] = 0.0

	Q := priorityQueue{{node: u, dist: 0, direction: true}, {node: v, dist: 0,
		direction: false}}
	heap.Init(&Q)

	lengthBestPath := math.Inf(1)
	var midPathNode *Node
	mid, mid_backward := &distanceNode{node: nil}, &distanceNode{node: nil}

	for len(Q) > 0 {

		// terminates when no shorter paths can be found
		if mid.dist+mid_backward.dist >= lengthBestPath {
			break
		}

		bestNode := heap.Pop(&Q).(*distanceNode)

		if bestNode.direction {
			/* Forward Search */

			// if the next source has traversed a greater distance than recorded,
			// skip it
			mid = bestNode
			if mid.dist <= forwardDist[mid.node] {
				for _, e := range mid.node.EdgeStart {
					n := e.To
					new_distance := e.Weight

					// total distance travelled so far
					acc_dist := forwardDist[mid.node] + new_distance

					// update shortest paths
					if dist, ok := forwardDist[n]; !ok || acc_dist < dist {
						heap.Push(&Q, &distanceNode{node: n, dist: acc_dist, direction: true})
						forwardDist[n] = acc_dist
						next[n] = mid.node

						// update length of best path if it exists
						_, ok = backwardDist[n]
						if newLength := backwardDist[n] + forwardDist[n]; ok && newLength < lengthBestPath {
							lengthBestPath = newLength
							midPathNode = n
						}
					}
				}
			}
		}

		if !bestNode.direction {

			/* Reverse Search */
			mid_backward = bestNode
			if mid_backward.dist <= backwardDist[mid_backward.node] {
				for _, e := range mid_backward.node.EdgeEnd {
					n := e.From
					new_distance := e.Weight

					// total distance travelled so far
					acc_dist := backwardDist[mid_backward.node] + new_distance

					// update shortest paths
					if dist, ok := backwardDist[n]; !ok || acc_dist < dist {
						heap.Push(&Q, &distanceNode{node: n, dist: acc_dist, direction: false})
						backwardDist[n] = acc_dist
						back[n] = mid_backward.node

						// update length of best path if it exists
						_, ok = forwardDist[n]
						if newLength := backwardDist[n] + forwardDist[n]; ok && newLength < lengthBestPath {
							lengthBestPath = newLength
							midPathNode = n
						}
					}

				}
			}
		}
	}

	/* Get shortest path */

	// no path found
	if midPathNode == nil {
		return nil, math.Inf(1)
	}

	path := []*Node{midPathNode}

	n := next[midPathNode]

	for n != u {
		path = append(path, n)
		n = next[n]
	}

	path = append(path, u)

	// reverse the path
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}

	n = midPathNode
	for n != v {
		path = append(path, n)
		n = back[n]
	}

	path = append(path, v)

	return path, lengthBestPath
}
