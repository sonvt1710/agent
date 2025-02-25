package dag

// WalkFunc is a function that gets invoked when walking a Graph. Walking will
// stop if WalkFunc returns a non-nil error.
type WalkFunc func(n Node) error

// WalkFuncWithParent is like WalkFunc, but also provides the parent node of each node visited (which can be nil).
type WalkFuncWithParent func(n Node, parent Node) error

// Walk performs a depth-first walk of outgoing edges for all nodes in start,
// invoking the provided fn for each node. Walk returns the error returned by
// fn.
//
// Nodes unreachable from start will not be passed to fn.
func Walk(g *Graph, start []Node, fn WalkFunc) error {
	var (
		visited   = make(nodeSet)
		unchecked = make([]Node, 0, len(start))
	)

	// Prefill the set of unchecked nodes with our start set.
	unchecked = append(unchecked, start...)

	// Iterate through unchecked nodes, visiting each in turn and adding outgoing
	// edges to the unchecked list until all reachable nodes have been processed.
	for len(unchecked) > 0 {
		check := unchecked[len(unchecked)-1]
		unchecked = unchecked[:len(unchecked)-1]

		if visited.Has(check) {
			continue
		}
		visited.Add(check)

		if err := fn(check); err != nil {
			return err
		}

		for n := range g.outEdges[check] {
			unchecked = append(unchecked, n)
		}
	}

	return nil
}

// WalkReverse performs a depth-first walk of incoming edges for start node, invoking the provided fn for each node.
// Walk returns the error returned by fn.
//
// Nodes unreachable from start will not be passed to fn.
func WalkReverse(g *Graph, start Node, fn WalkFuncWithParent) error {
	type nodeWithParent struct {
		node   Node
		parent Node
	}

	var (
		visited   = make(nodeSet)
		unchecked = make([]nodeWithParent, 0, 1)
	)

	// Prefill the set of unchecked nodes with our start set.
	unchecked = append(unchecked, nodeWithParent{
		node:   start,
		parent: nil,
	})

	// Iterate through unchecked nodes, visiting each in turn and adding incoming
	// edges to the unchecked list until all reachable nodes have been processed.
	for len(unchecked) > 0 {
		check := unchecked[len(unchecked)-1]
		unchecked = unchecked[:len(unchecked)-1]

		if visited.Has(check.node) {
			continue
		}
		visited.Add(check.node)

		if err := fn(check.node, check.parent); err != nil {
			return err
		}

		for n := range g.inEdges[check.node] {
			unchecked = append(unchecked, nodeWithParent{
				node:   n,
				parent: check.node,
			})
		}
	}

	return nil
}

// WalkTopological performs a topological walk of all nodes in start in
// dependency order: a node will not be visited until its outgoing edges are
// visited first.
//
// Nodes will not be passed to fn if they are not reachable from start or if
// not all of their outgoing edges are reachable from start.
func WalkTopological(g *Graph, start []Node, fn WalkFunc) error {
	// NOTE(rfratto): WalkTopological is an implementation of Kahn's algorithm
	// which leaves g unmodified.

	var (
		visited   = make(nodeSet)
		unchecked = make([]Node, 0, len(start))

		remainingDeps = make(map[Node]int)
	)

	// Pre-fill the set of nodes to check from the start list.
	unchecked = append(unchecked, start...)

	for len(unchecked) > 0 {
		check := unchecked[len(unchecked)-1]
		unchecked = unchecked[:len(unchecked)-1]

		if visited.Has(check) {
			continue
		}
		visited.Add(check)

		if err := fn(check); err != nil {
			return err
		}

		// Iterate through the incoming edges to check and queue nodes if we're the
		// last edge to be walked.
		for n := range g.inEdges[check] {
			// remainingDeps starts with the number of edges, and we subtract one for
			// each outgoing edge that's visited.
			if _, ok := remainingDeps[n]; !ok {
				remainingDeps[n] = len(g.outEdges[n])
			}
			remainingDeps[n]--

			// Only enqueue the incoming edge once all of its outgoing edges have
			// been consumed. This prevents it from being visited before its
			// dependencies.
			if remainingDeps[n] == 0 {
				unchecked = append(unchecked, n)
			}
		}
	}

	return nil
}
