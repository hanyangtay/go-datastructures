/*
Package rtree implements this space partitioning datastructurein 2D setting
https://en.wikipedia.org/wiki/R*_tree

Useful for querying nearest neighbours of an object.
*/

package rtree

import (
	"math"
)

// Rtree represents the balanced search tree for storing and querying 2D data.
type Rtree struct {
	MinBranch int
	MaxBranch int
	Root      *rTreeNode
	Size      int
	Height    int
}

// node represents a tree node of an R tree, which contains multiple entries
type rTreeNode struct {
	parent  *rTreeNode
	isLeaf  bool
	entries []entry
	level   int
}

// entry represents a spatial index record stored in a tree node
type entry struct {
	bb    *Rect // bounding-box of all children of this entry
	child *rTreeNode
	obj   Spatial
}

// any spatial object can fulfill this interface - e.g. point, line, rectangle
// ToRect returns the min bounding box of an object
// Dist returns the distance of an object to a bounding box
// Refer to type RTreePoint for an example of a type that fulfills this interface
type Spatial interface {
	ToRect() *Rect
	SquaredDist(*Rect) float64
}

// NewTree initialises a new R-tree with a specified min and max number of branches.
func NewTree(MinBranch, MaxBranch int) *Rtree {
	return &Rtree{
		MinBranch: MinBranch,
		MaxBranch: MaxBranch,
		Root: &rTreeNode{
			entries: make([]entry, 0, MaxBranch),
			isLeaf:  true,
			level:   1,
		},
		Size:   0,
		Height: 1,
	}
}

/* Insertion */

// Insert inserts a spatial object into the tree.
// Tree is rebalanced if a leaf node overflows.
func (tree *Rtree) Insert(obj Spatial) {
	e := entry{obj.ToRect(), nil, obj}
	tree.insert(e, 1)
	tree.Size++
}

// insert adds specified entry to the tree at the specified level
func (tree *Rtree) insert(e entry, level int) {
	leaf := tree.chooseNode(tree.Root, e, level)
	leaf.entries = append(leaf.entries, e)

	// update parent pointer if necessary
	if e.child != nil {
		e.child.parent = leaf
	}

	// split leaf if it overflows
	var split *rTreeNode
	if len(leaf.entries) > tree.MaxBranch {
		leaf, split = leaf.split(tree.MinBranch)
	}

	// adjusts the tree and rebalances if necessary
	_, _ = tree.adjustTree(leaf, split)
}

// chooseNode finds the node at the specified level to which e should be added
func (tree *Rtree) chooseNode(n *rTreeNode, e entry, level int) *rTreeNode {

	if n.isLeaf || n.level == level {
		return n
	}

	// find the entry whose bb needs least enlargement to include obj
	leastDiff := math.MaxFloat64
	var chosen entry
	var bb Rect
	for _, e2 := range n.entries {
		initBoundingBox(&bb, e2.bb, e.bb)
		diff := bb.size - e2.bb.size

		// choose the smallest difference, or the smaller bounding box size in a tie
		if diff < leastDiff || (diff == leastDiff && e2.bb.size < chosen.bb.size) {
			leastDiff = diff
			chosen = e2
		}
	}

	return tree.chooseNode(chosen.child, e, level)
}

// adjustTree splits overflowing nodes and propagates the changes upwards
func (tree *Rtree) adjustTree(leaf, split *rTreeNode) (*rTreeNode, *rTreeNode) {

	// edge case: handle Root adjustments
	if leaf == tree.Root {
		if split != nil {
			tree.Height++
			tree.Root = &rTreeNode{
				parent: nil,
				isLeaf: false,
				level:  tree.Height,
				entries: []entry{
					entry{bb: leaf.computeBoundingBox(), child: leaf},
					entry{bb: split.computeBoundingBox(), child: split},
				},
			}
			leaf.parent = tree.Root
			split.parent = tree.Root

			return leaf, split
		} else {
			return nil, nil
		}
	}

	// resize the bounding box of n from lower level changes
	e := leaf.getEntry()
	e.bb = leaf.computeBoundingBox()

	// if no split, just propagate changes upwards
	if split == nil {
		return tree.adjustTree(leaf.parent, nil)
	}

	// leaf was used as the "left" node, but need to add split to leaf's parent
	new_entry := entry{split.computeBoundingBox(), split, nil}
	leaf.parent.entries = append(leaf.parent.entries, new_entry)

	// if split entry overflows parent, split parent and propagate
	if len(leaf.parent.entries) > tree.MaxBranch {
		return tree.adjustTree(leaf.parent.split(tree.MinBranch))
	}

	// otherwise continue to propagate changes upwards
	return tree.adjustTree(leaf.parent, nil)
}

// getEntry returns a pointer to the entry for the node n from n's parent
func (n *rTreeNode) getEntry() *entry {
	for i := range n.parent.entries {
		if n.parent.entries[i].child == n {
			return &n.parent.entries[i]
		}
	}
	panic("getEntry returns nil pointer!")
	return nil
}

// computeBoundingBox finds the bb of the children of n
func (n *rTreeNode) computeBoundingBox() *Rect {
	var bb Rect
	for i, e := range n.entries {
		if i == 0 {
			bb = *e.bb
		} else {
			bb.enlarge(e.bb)
		}
	}
	return &bb
}

// split splits a node into two groups while attempting to minimise the
// bounding box area of the split groups
func (n *rTreeNode) split(minBranch int) (left, right *rTreeNode) {

	// finds the initial split
	l, r := n.pickSeeds()
	leftSeed, rightSeed := n.entries[l], n.entries[r]

	// get remaining entries to be divided between left and right
	remaining := append(n.entries[:l], n.entries[l+1:r]...)
	remaining = append(remaining, n.entries[r+1:]...)

	// initialise new split nodes (reuse n as left node)
	left = n
	left.entries = []entry{leftSeed}

	right = &rTreeNode{
		parent:  n.parent,
		isLeaf:  n.isLeaf,
		level:   n.level,
		entries: []entry{rightSeed},
	}

	if rightSeed.child != nil {
		rightSeed.child.parent = right
	}

	if leftSeed.child != nil {
		leftSeed.child.parent = left
	}

	// distribute remaining entries into left and right split nodes
	assignGroup(remaining, left, right, minBranch)

	return
}

// pickSeeds chooses two child entries of n to start a split
// by choosing the entries that result in least overlap
// i.e. has greatest waste of space
func (n *rTreeNode) pickSeeds() (int, int) {
	left, right := 0, 1
	maxWastedSpace := -1.0
	var bb Rect
	for i, e1 := range n.entries {
		for j, e2 := range n.entries[i+1:] {
			initBoundingBox(&bb, e1.bb, e2.bb)
			diff := bb.size - e1.bb.size - e2.bb.size
			if diff > maxWastedSpace {
				maxWastedSpace = diff
				left, right = i, i+1+j
			}
		}
	}

	return left, right
}

// assign adds an entry to a split node
func assign(e entry, splitNode *rTreeNode) {
	if e.child != nil {
		e.child.parent = splitNode
	}
	splitNode.entries = append(splitNode.entries, e)
}

// assignGroup adds entries to either of the two split nodes
func assignGroup(remaining []entry, left, right *rTreeNode, minBranch int) {

	var nextIdx int
	var bestLeftDiff, bestRightDiff float64

	for len(remaining) > 0 {

		maxDiff := -1.0

		leftBB := left.computeBoundingBox()
		rightBB := right.computeBoundingBox()

		for i, e := range remaining {
			leftDiff := boundingBox(leftBB, e.bb).size - leftBB.size
			rightDiff := boundingBox(rightBB, e.bb).size - rightBB.size
			diff := math.Abs(leftDiff - rightDiff)
			if diff > maxDiff {
				maxDiff = diff

				bestLeftDiff = leftDiff
				bestRightDiff = rightDiff
				nextIdx = i
			}
		}

		nextE := remaining[nextIdx]

		switch {
		// prevent underflow of branches
		case len(left.entries)+len(remaining) <= minBranch:
			assign(nextE, left)
		case len(right.entries)+len(remaining) <= minBranch:
			assign(nextE, right)

		// choose the split node with least enlargement
		case bestLeftDiff < bestRightDiff:
			assign(nextE, left)
		case bestRightDiff < bestLeftDiff:
			assign(nextE, right)

		// in tie, choose the split node with the smaller area
		case leftBB.size < rightBB.size:
			assign(nextE, left)
		case rightBB.size < leftBB.size:
			assign(nextE, right)

		// in tie, choose node with fewer entries
		case len(left.entries) < len(right.entries):
			assign(nextE, left)
		default:
			assign(nextE, right)
		}

		// update remaining entries to be inserted
		remaining = append(remaining[:nextIdx], remaining[nextIdx+1:]...)
	}
}

/* Deletion */

// Delete removes an object from the tree.
// Returns true if object is not found, or false if it's a no-op
func (tree *Rtree) Delete(obj Spatial) bool {
	n := tree.findLeaf(tree.Root, obj)

	if n == nil {
		return false
	}

	idx := -1
	for i, e := range n.entries {
		if e.obj == obj {
			idx = i
		}
	}
	if idx == -1 {
		return false
	}

	n.entries = append(n.entries[:idx], n.entries[idx+1:]...)

	tree.condenseTree(n)
	tree.Size--

	// edge case: only Root is left and it's not a leaf node
	if !tree.Root.isLeaf && len(tree.Root.entries) == 1 {
		tree.Root = tree.Root.entries[0].child
	}

	return true
}

// findLeaf finds the leaf node containing obj
func (tree *Rtree) findLeaf(n *rTreeNode, obj Spatial) *rTreeNode {
	if n.isLeaf {
		return n
	}

	for _, e := range n.entries {
		if e.bb.containsRect(obj.ToRect()) {
			leaf := tree.findLeaf(e.child, obj)
			if leaf == nil {
				continue
			}
			for _, leafEntry := range leaf.entries {
				if leafEntry.obj == obj {
					return leaf
				}
			}
		}
	}

	return nil
}

// condenseTree deletes underflowing nodes and propagates changes upwards
func (tree *Rtree) condenseTree(n *rTreeNode) {
	deletedNodes := []*rTreeNode{}

	for n != tree.Root {
		if len(n.entries) < tree.MinBranch {

			// remove n from parent entries
			entries := []entry{}
			for _, e := range n.parent.entries {
				if e.child != n {
					entries = append(entries, e)
				}
			}

			if len(n.parent.entries) == len(entries) {
				panic("RTree Delete: Failed to remove entry from parent.")
			}
			n.parent.entries = entries

			// only add n to deleted nodes if it still has children to be reinserted
			if len(n.entries) > 0 {
				deletedNodes = append(deletedNodes, n)
			}
		} else {
			// child entry deletion, no underflow
			n.getEntry().bb = n.computeBoundingBox()
		}

		n = n.parent
	}

	for _, n := range deletedNodes {
		// reinsert entry at the same level
		entry := entry{n.computeBoundingBox(), n, nil}
		tree.insert(entry, n.level+1)
	}
}
