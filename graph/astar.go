package graph

import (
	"container/heap"
	"math"
)

// AStar returns a shortest path from u to v and the distance
// in graph g. Heuristic: great circle distance, Time complexity: O(|E| * log |V|)
func (g *DirectedGraph) AStar(u, v *Node) ([]*Node, float64) {

	forwardDist := make(map[*Node]float64)
	forwardDist[u] = 0.0

	Q := priorityQueue{{node: u, dist: 0 + Dist(u, v)}}
	var mid *distanceNode
	heap.Init(&Q)

	next := make(map[*Node]*Node)

	for len(Q) > 0 {
		mid = heap.Pop(&Q).(*distanceNode)

		// terminates when final node is found
		if mid.node == v {
			break
		}

		for _, e := range mid.node.EdgeStart {
			n := e.To

			new_distance := e.Weight

			// total distance travelled so far
			acc_dist := forwardDist[mid.node] + new_distance

			// update shortest paths
			if dist, ok := forwardDist[n]; !ok || acc_dist < dist {
				heap.Push(&Q, &distanceNode{node: n, dist: acc_dist + Dist(n, v)})
				forwardDist[n] = acc_dist
				next[n] = mid.node
			}

		}
	}

	// no path found
	if mid == nil || mid.node != v {
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

// AStarBi returns a shortest path from u to all nodes
// in the graph g. Time complexity: O(|E| * log |V|)
func (g *DirectedGraph) AStarBi(u, v *Node) ([]*Node, float64) {

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
		if mid.realDist+mid_backward.realDist >= lengthBestPath {
			break
		}

		bestNode := heap.Pop(&Q).(*distanceNode)

		if bestNode.direction {

			/* Forward Search */

			// if the next source has traversed a greater distance than recorded,
			// skip it
			mid = bestNode
			if mid.realDist <= forwardDist[mid.node] {
				for _, e := range mid.node.EdgeStart {
					n := e.To

					// acc_dist = total distance travelled so far
					acc_dist := forwardDist[mid.node] + e.Weight

					// update shortest paths
					if dist, ok := forwardDist[n]; !ok || acc_dist < dist {
						heap.Push(&Q, &distanceNode{node: n, dist: acc_dist + Dist(n, v),
							realDist: acc_dist, direction: true})
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
		} else {

			/* Reverse Search */
			mid_backward = bestNode
			if mid_backward.realDist <= backwardDist[mid_backward.node] {
				for _, e := range mid_backward.node.EdgeEnd {
					n := e.From

					// total distance travelled so far
					acc_dist := backwardDist[mid_backward.node] + e.Weight

					// update shortest paths
					if dist, ok := backwardDist[n]; !ok || acc_dist < dist {
						heap.Push(&Q, &distanceNode{node: n, dist: acc_dist + Dist(n, u),
							realDist: acc_dist, direction: false})
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
