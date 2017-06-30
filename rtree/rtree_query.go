package rtree

import (
	"container/heap"
)

/* Querying */

// SearchIntersect returns all spatial objects that intersect the specified bounding box.
func (tree *Rtree) SearchIntersect(bb *Rect) []Spatial {
	results := []Spatial{}
	return tree.searchIntersect(tree.Root, bb, results)
}

func (tree *Rtree) searchIntersect(n *rTreeNode, bb *Rect, results []Spatial) []Spatial {
	for _, e := range n.entries {
		if intersect(e.bb, bb) {
			if n.isLeaf {
				results = append(results, e.obj)
			} else {
				results = tree.searchIntersect(e.child, bb, results)
			}
		}
	}

	return results
}

/* Priority queue for knn */

type distRTreeNode struct {
	rEntry entry
	dist   float64
}

type priorityRQueue []*distRTreeNode

func (q priorityRQueue) Len() int { return len(q) }

func (q priorityRQueue) Less(i, j int) bool {
	return q[i].dist < q[j].dist
}

func (q priorityRQueue) Swap(i, j int) {
	q[i], q[j] = q[j], q[i]
}

func (q *priorityRQueue) Push(n interface{}) {
	*q = append(*q, n.(*distRTreeNode))
}

func (q *priorityRQueue) Pop() interface{} {
	old := *q
	n := len(old)
	node := old[n-1]
	*q = old[0 : n-1]
	return node
}

// KNearestNeighbours returns k nearest spatial objects and their distances
func (tree *Rtree) KNN(k int, point Spatial) []Spatial {

	nearestNeighbours := make([]Spatial, 0, k)

	Q := priorityRQueue{}
	for _, e := range tree.Root.entries {
		newQNode := &distRTreeNode{e, point.SquaredDist(e.bb)}
		heap.Push(&Q, newQNode)
	}

	for len(Q) > 0 && len(nearestNeighbours) < k {
		mid := heap.Pop(&Q).(*distRTreeNode)

		if mid.rEntry.obj != nil {
			nearestNeighbours = append(nearestNeighbours, mid.rEntry.obj)
		} else {
			for _, e := range mid.rEntry.child.entries {
				newQNode := &distRTreeNode{e, point.SquaredDist(e.bb)}
				heap.Push(&Q, newQNode)
			}
		}
	}

	return nearestNeighbours
}
