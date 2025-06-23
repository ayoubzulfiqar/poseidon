package main

import (
	"fmt"
)

type Graph struct {
	adj map[int][]int
}

func NewGraph() *Graph {
	return &Graph{
		adj: make(map[int][]int),
	}
}

func (g *Graph) AddEdge(u, v int) {
	g.adj[u] = append(g.adj[u], v)
	g.adj[v] = append(g.adj[v], u)
}

func (g *Graph) dfs(node int, visited map[int]bool, component *[]int) {
	visited[node] = true
	*component = append(*component, node)

	for _, neighbor := range g.adj[node] {
		if !visited[neighbor] {
			g.dfs(neighbor, visited, component)
		}
	}
}

func (g *Graph) FindConnectedComponents() [][]int {
	visited := make(map[int]bool)
	var components [][]int

	nodesInGraph := make(map[int]struct{})
	for u, neighbors := range g.adj {
		nodesInGraph[u] = struct{}{}
		for _, v := range neighbors {
			nodesInGraph[v] = struct{}{}
		}
	}

	for node := range nodesInGraph {
		if !visited[node] {
			var currentComponent []int
			g.dfs(node, visited, &currentComponent)
			components = append(components, currentComponent)
		}
	}

	return components
}

func main() {
	graph := NewGraph()

	graph.AddEdge(0, 1)
	graph.AddEdge(0, 2)
	graph.AddEdge(1, 3)

	graph.AddEdge(4, 5)

	graph.AddEdge(6, 6)

	components := graph.FindConnectedComponents()

	fmt.Println("Connected Components:")
	for i, comp := range components {
		fmt.Printf("Component %d: %v\n", i+1, comp)
	}

	fmt.Println("\nTest Case 2: Empty graph")
	emptyGraph := NewGraph()
	emptyComponents := emptyGraph.FindConnectedComponents()
	fmt.Println("Connected Components:", emptyComponents)

	fmt.Println("\nTest Case 3: Graph with multiple components including isolated ones")
	graphWithIsolated := NewGraph()
	graphWithIsolated.AddEdge(10, 11)
	graphWithIsolated.AddEdge(12, 12)
	graphWithIsolated.AddEdge(13, 14)
	graphWithIsolated.AddEdge(14, 15)

	isolatedComponents := graphWithIsolated.FindConnectedComponents()
	fmt.Println("Connected Components:")
	for i, comp := range isolatedComponents {
		fmt.Printf("Component %d: %v\n", i+1, comp)
	}
}

// Additional implementation at 2025-06-23 02:55:53
package graph

import (
	"sort"
)

// Graph represents an undirected graph using an adjacency list.
type Graph struct {
	adj map[int][]int
	// vertices stores all unique vertices in the graph, useful for iterating through all nodes
	vertices map[int]struct{}
}

// NewGraph creates and returns a new empty Graph.
func NewGraph() *Graph {
	return &Graph{
		adj:      make(map[int][]int),
		vertices: make(map[int]struct{}),
	}
}

// AddEdge adds an undirected edge between vertices u and v.
// If vertices do not exist, they are added.
func (g *Graph) AddEdge(u, v int) {
	g.adj[u] = append(g.adj[u], v)
	g.adj[v] = append(g.adj[v], u)
	g.vertices[u] = struct{}{}
	g.vertices[v] = struct{}{}
}

// AddVertex ensures a vertex exists in the graph, even if it has no edges.
func (g *Graph) AddVertex(v int) {
	if _, exists := g.adj[v]; !exists {
		g.adj[v] = []int{} // Initialize an empty adjacency list for isolated vertex
	}
	g.vertices[v] = struct{}{}
}

// GetVertices returns a sorted slice of all vertices in the graph.
func (g *Graph) GetVertices() []int {
	var verts []int
	for v := range g.vertices {
		verts = append(verts, v)
	}
	sort.Ints(verts) // Sort for consistent output
	return verts
}

// ConnectedComponentsResult holds the results of a connected components analysis.
type ConnectedComponentsResult struct {
	// Components is a slice of slices, where each inner slice represents a connected component
	// and contains the vertices belonging to that component.
	Components [][]int
	// VertexToComponentID maps each vertex to its corresponding component ID (index in Components).
	VertexToComponentID map[int]int
	// NumVertices is the total number of unique vertices processed.
	NumVertices int
}

// NumComponents returns the total number of connected components found.
func (ccr *ConnectedComponentsResult) NumComponents() int {
	return len(ccr.Components)
}

// AreInSameComponent checks if two vertices belong to the same connected component.
// Returns true if they are in the same component, false otherwise.
// Returns false if either vertex does not exist in the graph that was analyzed.
func (ccr *ConnectedComponentsResult) AreInSameComponent(v1, v2 int) bool {
	id1, ok1 := ccr.VertexToComponentID[v1]
	id2, ok2 := ccr.VertexToComponentID[v2]

	if !ok1 || !ok2 {
		return false // One or both vertices not found in the graph
	}
	return id1 == id2
}

// ComponentOf returns the component ID (index in Components) for a given vertex.
// The second return value is true if the vertex was found, false otherwise.
func (ccr *ConnectedComponentsResult) ComponentOf(v int) (int, bool) {
	id, ok := ccr.VertexToComponentID[v]
	return id, ok
}

// GetComponentVertices returns a sorted slice of vertices for a given component ID.
// Returns nil if the component ID is out of bounds.
func (ccr *ConnectedComponentsResult) GetComponentVertices(componentID int) []int {
	if componentID < 0 || componentID >= len(ccr.Components) {
		return nil
	}
	// Return a copy to prevent external modification of internal state
	component := make([]int, len(ccr.Components[componentID]))
	copy(component, ccr.Components[componentID])
	sort.Ints(component) // Ensure sorted output
	return component
}

// FindConnectedComponents finds all connected components in the graph.
// It returns a ConnectedComponentsResult containing the components and a mapping
// from vertices to their component IDs.
func FindConnectedComponents(g *Graph) *ConnectedComponentsResult {
	visited := make(map[int]bool)
	vertexToComponentID := make(map[int]int)
	var components [][]int
	componentIDCounter := 0

	// Iterate through all vertices to ensure isolated vertices are also processed
	// and to start DFS from unvisited vertices.
	allVertices := g.GetVertices()

	for _, v := range allVertices {
		if !visited[v] {
			var currentComponent []int
			// Start a DFS from the unvisited vertex
			stack := []int{v}
			visited[v] = true
			vertexToComponentID[v] = componentIDCounter
			currentComponent = append(currentComponent, v)

			for len(stack) > 0 {
				curr := stack[len(stack)-1]
				stack = stack[:len(stack)-1] // Pop

				// Get neighbors, ensure they exist in the graph's adjacency map
				neighbors, ok := g.adj[curr]
				if !ok {
					// This case should ideally not happen if AddEdge/AddVertex are used correctly,
					// but good for robustness if graph is built manually.
					continue
				}

				for _, neighbor := range neighbors {
					if !visited[neighbor] {
						visited[neighbor] = true
						vertexToComponentID[neighbor] = componentIDCounter
						currentComponent = append(currentComponent, neighbor)
						stack = append(stack, neighbor) // Push
					}
				}
			}
			// Sort vertices within the component for consistent output
			sort.Ints(currentComponent)
			components = append(components, currentComponent)
			componentIDCounter++
		}
	}

	return &ConnectedComponentsResult{
		Components:          components,
		VertexToComponentID: vertexToComponentID,
		NumVertices:         len(allVertices),
	}
}

// Additional implementation at 2025-06-23 02:56:57
package main

import (
	"fmt"
)

// Graph represents an adjacency list for an unweighted graph.
// Keys are node IDs, values are slices of neighbor node IDs.
type Graph map[int][]int

// ConnectedComponentsFinder finds and manages connected components in a graph.
type ConnectedComponentsFinder struct {
	graph             Graph
	visited           map[int]bool
	components        [][]int
	nodeToComponentID map[int]int // Maps node ID to its component index in 'components'
	componentCount    int
}

// NewConnectedComponentsFinder creates a new finder for a given graph.
func NewConnectedComponentsFinder(g Graph) *ConnectedComponentsFinder {
	return &ConnectedComponentsFinder{
		graph:             g,
		visited:           make(map[int]bool),
		components:        make([][]int, 0),
		nodeToComponentID: make(map[int]int),
		componentCount:    0,
	}
}

// FindComponents executes the DFS algorithm to find all connected components.
// This method must be called before querying component information.
func (ccf *ConnectedComponentsFinder) FindComponents() {
	// Reset state for re-computation or initial computation
	ccf.visited = make(map[int]bool)
	ccf.components = make([][]int, 0)
	ccf.nodeToComponentID = make(map[int]int)
	ccf.componentCount = 0

	// Collect all unique nodes present in the graph (including isolated ones)
	nodes := make(map[int]struct{})
	for u, neighbors := range ccf.graph {
		nodes[u] = struct{}{}
		for _, v := range neighbors {
			nodes[v] = struct{}{}
		}
	}

	// Iterate through all nodes to find components
	for node := range nodes {
		if !ccf.visited[node] {
			currentComponent := []int{}
			ccf.dfs(node, &currentComponent, ccf.componentCount)
			ccf.components = append(ccf.components, currentComponent)
			ccf.componentCount++
		}
	}
}

// dfs performs a Depth First Search from a given node,
// populating the currentComponent slice and mapping nodes to their component index.
func (ccf *ConnectedComponentsFinder) dfs(node int, currentComponent *[]int, componentIdx int) {
	ccf.visited[node] = true
	*currentComponent = append(*currentComponent, node)
	ccf.nodeToComponentID[node] = componentIdx

	for _, neighbor := range ccf.graph[node] {
		if !ccf.visited[neighbor] {
			ccf.dfs(neighbor, currentComponent, componentIdx)
		}
	}
}

// GetComponents returns the list of all found connected components.
// Each inner slice represents a component.
// Call FindComponents() first to ensure results are up-to-date.
func (ccf *ConnectedComponentsFinder) GetComponents() [][]int {
	return ccf.components
}

// GetComponentCount returns the total number of connected components found.
// Call FindComponents() first.
func (ccf *ConnectedComponentsFinder) GetComponentCount() int {
	return ccf.componentCount
}

// AreConnected checks if two nodes belong to the same connected component.
// Returns true if both nodes exist in the graph and are in the same component, false otherwise.
// Call FindComponents() first.
func (ccf *ConnectedComponentsFinder) AreConnected(node1, node2 int) bool {
	id1, ok1 := ccf.nodeToComponentID[node1]
	id2, ok2 := ccf.nodeToComponentID[node2]

	// If either node is not in the graph or not processed, they cannot be connected
	if !ok1 || !ok2 {
		return false
	}
	return id1 == id2
}

// IsGraphConnected checks if the entire graph is a single connected component.
// Returns true if there is exactly one component, false otherwise.
// An empty graph will return false (0 components). A graph with a single node returns true.
// Call FindComponents() first.
func (ccf *ConnectedComponentsFinder) IsGraphConnected() bool {
	return ccf.componentCount == 1
}

// GetComponentSize returns the number of nodes in the component containing the given node.
// Returns -1 if the node is not found in any component (i.e., not in the graph).
// Call FindComponents() first.
func (ccf *ConnectedComponentsFinder) GetComponentSize(node int) int {
	componentIdx, ok := ccf.nodeToComponentID[node]
	if !ok {
		return -1 // Node not found or not part of any component
	}
	return len(ccf.components[componentIdx])
}

// GetAllComponentSizes returns a slice containing the size (number of nodes) of each component.
// The order of sizes corresponds to the order of components returned by GetComponents().
// Call FindComponents() first.
func (ccf *ConnectedComponentsFinder) GetAllComponentSizes() []int {
	sizes := make([]int, len(ccf.components))
	for i, comp := range ccf.components {
		sizes[i] = len(comp)
	}
	return sizes
}

func main() {
	// Example Graph 1: Disconnected graph with multiple components and an isolated node
	graph1 := Graph{
		1: {2},
		2: {1},
		3: {4},
		4: {3},
		5: {}, // Isolated node
	}

	finder1 := NewConnectedComponentsFinder(graph1)
	finder1.FindComponents()

	fmt.Println("--- Graph 1 ---")
	fmt.Println("Components:", finder1.GetComponents())
	fmt.Println("Component Count:", finder1.GetComponentCount())
	fmt.Println("Are 1 and 2 connected?", finder1.AreConnected(1, 2))
	fmt.Println("Are 1 and 3 connected?", finder1.AreConnected(1, 3))
	fmt.Println("Is Graph 1 connected?", finder1.IsGraphConnected())
	fmt.Println("Size of component containing 1:", finder1.GetComponentSize(1))
	fmt.Println("Size of component containing 5:", finder1.GetComponentSize(5))
	fmt.Println("All component sizes for Graph 1:", finder1.GetAllComponentSizes())
	fmt.Println("Are 1 and 99 connected (99 not in graph)?", finder1.AreConnected(1, 99))
	fmt.Println("Size of component containing 99 (99 not in graph)?", finder1.GetComponentSize(99))

	// Example Graph 2: Fully connected graph
	graph2 := Graph{
		1: {2, 3},
		2: {1, 4},
		3: {1, 5},
		4: {2, 6},
		5: {3, 6},
		6: {4, 5},
	}

	finder2 := NewConnectedComponentsFinder(graph2)
	finder2.FindComponents()

	fmt.Println("\n--- Graph 2 ---")
	fmt.Println("Components:", finder2.GetComponents())
	fmt.Println("Component Count:", finder2.GetComponentCount())
	fmt.Println("Are 1 and 6 connected?", finder2.AreConnected(1, 6))
	fmt.Println("Is Graph 2 connected?", finder2.IsGraphConnected())
	fmt.Println("Size of component containing 3:", finder2.GetComponentSize(3))
	fmt.Println("All component sizes for Graph 2:", finder2.GetAllComponentSizes())

	// Example Graph 3: Empty graph
	graph3 := Graph{}
	finder3 := NewConnectedComponentsFinder(graph3)
	finder3.FindComponents()

	fmt.Println("\n--- Graph 3 (Empty) ---")
	fmt.Println("Components:", finder3.GetComponents())
	fmt.Println("Component Count:", finder3.GetComponentCount())
	fmt.Println("Is Graph 3 connected?", finder3.IsGraphConnected())

	// Example Graph 4: Single node graph
	graph4 := Graph{
		10: {},
	}
	finder4 := NewConnectedComponentsFinder(graph4)
	finder4.FindComponents()

	fmt.Println("\n--- Graph 4 (Single Node) ---")
	fmt.Println("Components:", finder4.GetComponents())
	fmt.Println("Component Count:", finder4.GetComponentCount())
	fmt.Println("Is Graph 4 connected?", finder4.IsGraphConnected())
	fmt.Println("Size of component containing 10:", finder4.GetComponentSize(10))
}