package graph

// BreadthFirstSearch traverses the graph via breadth first search.
func (g *DirectedGraph) BreadthFirstSearch(from *Node, visit func(u, v *Node)) {

	visited := make([]int, len(g.Nodes))
	visited[from.ID] = true
	queue := []*Node{from}

	for len(queue) > 0 {
		u := queue[0]
		queue = queue[1:]

		for _, v := range g.FindNeighboursFrom(u) {
			if visited[v.ID] {
				continue
			}

			//process vertex u, v
			if visit != nil {
				visit(u, v)
			}

			visited[v.ID] = true
			queue = append(queue, v)
		}
	}

}

// DepthFirstSearch traverses the graph via depth first search.
func (g *DirectedGraph) DepthFirstSearch(from *Node, visit func(u, v *Node)) {

	visited := make([]int, len(g.Nodes))
	visited[from.ID] = true
	stack := []*Node{from}

	for len(stack) > 0 {
		u := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		for _, v := range g.FindNeighboursFrom(u) {
			if visited[v.ID] {
				continue
			}

			//process vertex u, v
			if visit != nil {
				visit(u, v)
			}

			visited[v.ID] = true
			stack = append(stack, v)
		}
	}
}
