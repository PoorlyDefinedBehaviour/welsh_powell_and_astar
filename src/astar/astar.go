package astar

import (
	"container/heap"
	"fmt"
	"math"

	"github.com/pkg/errors"
)

type Item struct {
	vertex Vertex
	fCost  float32
}

type PriorityQueue []*Item

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].fCost < pq[j].fCost
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *PriorityQueue) Push(x interface{}) {
	item, ok := x.(*Item)
	if !ok {
		panic("pq.Push: argument is not a *Item")
	}

	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	*pq = old[0 : n-1]

	return item
}

type Vertex struct {
	Value     string  `json:"value"`
	Latitude  float32 `json:"latitude"`
	Longitude float32 `json:"longitude"`
}

var none = Vertex{
	Value:     "none",
	Latitude:  -1.0,
	Longitude: -1.0,
}

type Graph map[string][]Vertex

func (graph Graph) vertices() []string {
	vertices := make([]string, 0, len(graph))

	for vertex := range graph {
		vertices = append(vertices, vertex)
	}

	return vertices
}

func (graph Graph) neighbors(vertex string) []Vertex {
	return (graph)[vertex]
}

func manhattanDistance(x1 float32, y1 float32, x2 float32, y2 float32) float32 {
	return float32(math.Abs(float64(x1-x2)) + math.Abs(float64(y1-y2)))
}

func (graph Graph) distance(from Vertex, to Vertex) float32 {
	return manhattanDistance(from.Latitude, from.Longitude, to.Latitude, to.Longitude)
}

func reverse(xs []Vertex) []Vertex {
	ys := make([]Vertex, 0, len(xs))

	for i := len(xs) - 1; i > -1; i-- {
		ys = append(ys, xs[i])
	}

	return ys
}

func buildPath(pathMap map[string]Vertex, from Vertex, to Vertex) []Vertex {
	path := []Vertex{}

	current := to

	for current != from {
		path = append(path, current)
		current = pathMap[current.Value]
	}

	path = append(path, from)

	return reverse(path)
}

type Coordinates struct {
	Latitude  float32 `json:"latitude"`
	Longitude float32 `json:"longitude"`
}

type Path struct {
	From string `json:"from"`
	To   string `json:"to"`
}

type GraphInput struct {
	Cities map[string]Coordinates `json:"cities"`
	Path   Path                   `json:"path"`
	Graph  map[string][]string    `json:"graph"`
}

var errCityCoordinatesNotFound = errors.New("Coordinates not found")

func (input *GraphInput) coordinates(city string) (Coordinates, error) {
	coordinates, ok := input.Cities[city]
	if !ok {
		return coordinates, errors.Wrap(errCityCoordinatesNotFound, fmt.Sprintf("%s coordinates not found", city))
	}

	return coordinates, nil
}

func (input *GraphInput) cityToVertex(city string) (Vertex, error) {
	coordinates, err := input.coordinates(city)
	if err != nil {
		return Vertex{}, err
	}

	vertex := Vertex{
		Value:     city,
		Latitude:  coordinates.Latitude,
		Longitude: coordinates.Longitude,
	}

	return vertex, nil
}

func (input *GraphInput) toGraph() (Graph, error) {
	graph := Graph{}

	for city, neighbors := range input.Graph {
		adjacentVertices := make([]Vertex, 0, len(neighbors))

		for _, neighbor := range neighbors {
			vertex, err := input.cityToVertex(neighbor)
			if err != nil {
				return graph, err
			}

			adjacentVertices = append(adjacentVertices, vertex)
		}

		graph[city] = adjacentVertices
	}

	return graph, nil
}

func AStar(input *GraphInput) ([]Vertex, error) {
	graph, err := input.toGraph()
	if err != nil {
		return make([]Vertex, 0), err
	}

	from, err := input.cityToVertex(input.Path.From)
	if err != nil {
		return make([]Vertex, 0), err
	}

	to, err := input.cityToVertex(input.Path.To)
	if err != nil {
		return make([]Vertex, 0), err
	}

	priorityQueue := PriorityQueue{}
	visited := make(map[string]bool)
	previous := make(map[string]Vertex)

	for _, vertex := range graph.vertices() {
		previous[vertex] = none
	}

	heap.Push(&priorityQueue, &Item{
		vertex: from,
		fCost:  graph.distance(from, to),
	})

	for len(priorityQueue) != 0 {
		item, ok := heap.Pop(&priorityQueue).(*Item)
		if !ok {
			panic("couldn't cast priority queue element to *Item")
		}

		if item.vertex == to {
			break
		}

		visited[item.vertex.Value] = true

		for _, edge := range graph.neighbors(item.vertex.Value) {
			if visited[edge.Value] {
				continue
			}

			heap.Push(&priorityQueue, &Item{
				vertex: edge,
				fCost:  graph.distance(from, edge) + graph.distance(edge, to),
			})

			previous[edge.Value] = item.vertex
		}
	}

	path := buildPath(previous, from, to)

	return path, nil
}
