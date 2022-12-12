package main

import (
	"fmt"
	"math"
	"reflect"

	"github.com/artnoi43/gsl/data/graph"
	"github.com/artnoi43/gsl/data/graph/wgraph"
	"github.com/artnoi43/gsl/gslutils"
)

// This code purpose is to show an example of how to use the graph and wgraph packages
// With your own data structures, that is, without using the *Impl types.

// The code provides an example of solving the best flight path between major cities.
// In this case, type *city is used as nodes for both BFS and Dijkstra sort.

// Toward the end of this file, example functions are provided to get a better picture
// of how write function params for these graphs.

type cityName string

// *city implements wgraph.UndirectedNode
type city struct {
	name        cityName
	flightHours float64
	through     *city
}

func main() {
	infinity := math.MaxFloat64

	tokyo := &city{name: cityName("Tokyo"), flightHours: 0, through: nil}
	bangkok := &city{name: cityName("Bangkok"), flightHours: infinity, through: nil}
	hongkong := &city{name: cityName("Hongkong"), flightHours: infinity, through: nil}
	dubai := &city{name: cityName("Dubai"), flightHours: infinity, through: nil}
	helsinki := &city{name: cityName("Helsinki"), flightHours: infinity, through: nil}
	budapest := &city{name: cityName("Budapest"), flightHours: infinity, through: nil}

	// See file flight_graph.png
	graphEdges := map[wgraph.NodeWeighted[float64, cityName]][]struct {
		to       wgraph.NodeWeighted[float64, cityName]
		flighDur float64
	}{
		tokyo: {
			{
				to:       bangkok,
				flighDur: 4,
			},
			{
				to:       hongkong,
				flighDur: 1.5,
			},
			{
				to:       dubai,
				flighDur: 7,
			},
		},
		hongkong: {
			{
				to:       helsinki,
				flighDur: 11.5,
			},
		},
		bangkok: {
			{
				to:       dubai,
				flighDur: 6,
			},
			{
				to:       helsinki,
				flighDur: 9,
			},
		},
		dubai: {
			{
				to:       helsinki,
				flighDur: 3,
			},
			{
				to:       budapest,
				flighDur: 5,
			},
		},
		helsinki: {
			{
				to:       budapest,
				flighDur: 1.5,
			},
		},
		budapest: {},
	}

	directed := false
	dijkGraph := wgraph.NewDijkstraGraph[float64, cityName](directed)
	unweightedGraph := graph.NewGraph[float64](directed)

	// Add edges and nodes to graphs
	for node, nodeEdges := range graphEdges {
		dijkGraph.AddNode(node)
		unweightedGraph.AddNode(node)
		for _, nodeEdge := range nodeEdges {
			// fmt.Println(node.GetKey()+":", "adding edge", nodeEdge.to.GetKey(), "weight", nodeEdge.flighDur)
			if err := dijkGraph.AddEdge(node, nodeEdge.to, nodeEdge.flighDur); err != nil {
				panic("failed to add dijkstra-compatible graph edge: " + err.Error())
			}
			unweightedGraph.AddEdge(node, nodeEdge.to, nil)
		}
	}

	fromNode := tokyo
	shortestPathsFromTokyo := dijkGraph.DijkstraShortestPathFrom(fromNode)

	fmt.Println("Dijkstra result")
	for _, dst := range dijkGraph.GetNodes() {
		if dst == fromNode {
			continue
		}

		pathToNode := wgraph.DijkstraShortestPathReconstruct(shortestPathsFromTokyo.Paths, shortestPathsFromTokyo.From, dst)
		gslutils.ReverseInPlace(pathToNode)

		fmt.Println("> from", fromNode.GetKey(), "to", dst.GetKey(), "min flightHours", dst.GetValue())
		for _, via := range pathToNode {
			fmt.Printf("%s (%v) ", via.GetKey(), via.GetValue())
		}
		fmt.Println()
	}
	fmt.Println()

	fmt.Println("BFS result")
	for _, dst := range unweightedGraph.GetNodes() {
		shortestHopsFromTokyo, hops, found := graph.BFS[float64](unweightedGraph, fromNode, dst)
		fmt.Println("path to", dst.(*city).GetKey(), "found", found, "shortestHops", hops)
		gslutils.ReverseInPlace(shortestHopsFromTokyo)
		for i, hop := range shortestHopsFromTokyo {
			fmt.Println("hop", i, hop.(*city).GetKey())
		}
	}

	takeUnweightedGraph(unweightedGraph, fromNode)
	takeGraphWeighted(dijkGraph, fromNode)

	gg := unweightedGraph.(graph.GenericGraph[graph.Node[float64], graph.Node[float64], any])
	takeGenericGraph(gg)
}

func (c *city) GetValue() float64 {
	return c.flightHours
}

func (c *city) SetValueOrCost(newCost float64) {
	c.flightHours = newCost
}

func (c *city) GetKey() cityName {
	return c.name
}

func (c *city) GetPrevious() wgraph.NodeWeighted[float64, cityName] {
	if c.through == nil {
		return nil
	}
	return c.through
}

func (c *city) SetPrevious(node wgraph.NodeWeighted[float64, cityName]) {
	if node == nil {
		c.through = nil
		return
	}

	via, ok := node.(*city)
	if !ok {
		typeOfNode := reflect.TypeOf(node)
		panic(fmt.Sprintf("node not *city but %s", typeOfNode))
	}
	c.through = via
}

// graph.GenericGraph is impractical, as you can see from the type parameters..
func takeGenericGraph(
	gg graph.GenericGraph[graph.Node[float64], graph.Node[float64], any],
) {

	gg.AddEdge(nil, nil, 0)
}

// Instead of GenericGraph, use Graph or GraphWeighted
func takeUnweightedGraph(g graph.HashMapGraph[float64], from graph.Node[float64]) {

}

func takeGraphWeighted(g wgraph.HashMapGraphWeighted[float64, cityName], from wgraph.NodeWeighted[float64, cityName]) {

}