package rtree

import (
	"math"
)

type Rect struct {
	bottomLeft, topRight RTreePoint
	size                 float64
}

type RTreePoint struct {
	X float64
	Y float64
}

// ToRect constructs a bounding box containing the RTreePoint
func (n *RTreePoint) ToRect() *Rect {
	bottomLeft := RTreePoint{
		X: n.X - 0.00002,
		Y: n.Y - 0.00002,
	}

	topRight := RTreePoint{
		X: n.X + 0.00002,
		Y: n.Y + 0.00002,
	}

	return &Rect{
		bottomLeft: bottomLeft,
		topRight:   topRight,
		size:       0.00002 * 0.00002,
	}
}

// SquaredDist returns the square of the distance from point to rectangle
func (n *RTreePoint) SquaredDist(r *Rect) float64 {

	dist_x := math.Max(math.Max(r.bottomLeft.X-n.X, n.X-r.topRight.X), 0.0)
	dist_y := math.Max(math.Max(r.bottomLeft.Y-n.Y, n.Y-r.topRight.Y), 0.0)

	return dist_x*dist_x + dist_y*dist_y
}

// Dist returns the Euclidean distance of a point to a rectangle
func (n *RTreePoint) Dist(r *Rect) float64 {
	return math.sqrt(n.SquaredDist(r))
}

// NewRect initialises a new rectangle from two points
func NewRect(u, v *RTreePoint) *Rect {
	bottomLeft := RTreePoint{
		X: math.Min(u.X, v.X),
		Y: math.Min(u.Y, v.Y),
	}

	topRight := RTreePoint{
		X: math.Max(u.X, v.X),
		Y: math.Max(u.Y, v.Y),
	}

	return &Rect{
		bottomLeft: bottomLeft,
		topRight:   topRight,
		size:       (topRight.Y - bottomLeft.Y) * (topRight.X - bottomLeft.X),
	}
}

// containsRect tests whether r2 is located inside r1
func (r1 *Rect) containsRect(r2 *Rect) bool {
	if r1.bottomLeft.Y > r2.bottomLeft.Y || r1.bottomLeft.X > r2.bottomLeft.X {
		return false
	} else if r1.topRight.Y < r2.topRight.Y || r1.topRight.X < r1.topRight.X {
		return false
	}

	return true
}

// enlarge increases a rectangle bound to include
func (r1 *Rect) enlarge(r2 *Rect) {

	if r1.bottomLeft.X > r2.bottomLeft.X {
		r1.bottomLeft.X = r2.bottomLeft.X
	}

	if r1.topRight.X < r2.topRight.X {
		r1.topRight.X = r2.topRight.X
	}

	if r1.bottomLeft.Y > r2.bottomLeft.Y {
		r1.bottomLeft.Y = r2.bottomLeft.Y
	}

	if r1.topRight.Y < r2.topRight.Y {
		r1.topRight.Y = r2.topRight.Y
	}

	r1.size = (r1.topRight.Y - r1.bottomLeft.Y) * (r1.topRight.X - r1.bottomLeft.X)
}

// intersect tests whether there is an intersection between two rectangles
func intersect(r1, r2 *Rect) bool {

	if r2.topRight.X < r1.bottomLeft.X || r1.topRight.X < r2.bottomLeft.X {
		return false
	}

	if r2.topRight.Y < r1.bottomLeft.Y || r1.topRight.Y < r2.bottomLeft.Y {
		return false
	}

	return true

}

// boundingBox returns a rectangle that bounds both rectangles
func boundingBox(r1, r2 *Rect) *Rect {
	var r Rect
	initBoundingBox(&r, r1, r2)
	return &r
}

// initBoundingBox returns a rectangle that bounds both rectangles
func initBoundingBox(r, r1, r2 *Rect) {
	*r = *r1

	r.enlarge(r2)
}
