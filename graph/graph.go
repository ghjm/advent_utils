package graph

import (
	"fmt"
	"math"
)

type Edge[T comparable] struct {
	Dest T
	Cost uint64
}

type Graph[T comparable] struct {
	Nodes map[T][]Edge[T]
}

// checkInit checks that the graph data is initialized
func (g *Graph[T]) checkInit() {
	if g.Nodes == nil {
		g.Nodes = make(map[T][]Edge[T])
	}
}

// AddNode adds a node to the graph
func (g *Graph[T]) AddNode(node T) {
	g.checkInit()
	_, ok := g.Nodes[node]
	if !ok {
		g.Nodes[node] = nil
	}
}

// AddEdge adds an edge to the graph
func (g *Graph[T]) AddEdge(from, to T, cost uint64) {
	g.checkInit()
	g.AddNode(from)
	g.AddNode(to)
	g.Nodes[from] = append(g.Nodes[from], Edge[T]{Dest: to, Cost: cost})
}

// BuildStateGraph builds a graph, given a transition function between states
func (g *Graph[T]) BuildStateGraph(initialState T, transitionFunc func(T) []Edge[T]) {
	g.Nodes = nil
	g.checkInit()
	open := []T{initialState}
	visited := make(map[T]struct{})
	for len(open) > 0 {
		s := open[0]
		open = open[1:]
		if _, ok := visited[s]; ok {
			continue
		}
		visited[s] = struct{}{}
		g.AddNode(s)
		for _, tr := range transitionFunc(s) {
			g.AddEdge(s, tr.Dest, tr.Cost)
			open = append(open, tr.Dest)
		}
	}
	g.AddNode(initialState)
}

// Dijkstra finds the cost to all reachable states from a given source, along with "prev" data that can be used
// to reconstruct the path taken to each reachable destination.
func (g *Graph[T]) Dijkstra(source T) (map[T]uint64, map[T][]T) {
	Q := PriorityQueue[T, uint64]{}
	dist := make(map[T]uint64)
	prev := make(map[T][]T)
	dist[source] = 0
	Q.Insert(source, 0)
	for v := range g.Nodes {
		if v != source {
			prev[v] = nil
			dist[v] = math.MaxUint64
			Q.Insert(v, math.MaxUint64)
		}
	}
	for Q.Len() > 0 {
		u, err := Q.Pop()
		if err != nil {
			panic(fmt.Errorf("error popping value: %w", err))
		}
		for _, e := range g.Nodes[u] {
			v := e.Dest
			alt := dist[u] + e.Cost
			if alt < dist[v] {
				prev[v] = []T{u}
				dist[v] = alt
				Q.UpdatePriority(v, alt)
			} else if alt == dist[v] {
				prev[v] = append(prev[v], u)
			}
		}
	}
	return dist, prev
}
