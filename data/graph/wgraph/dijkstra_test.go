package wgraph

import (
	"testing"

	"github.com/artnoi43/gsl/gslutils"
)

type dijkstraTestUtils[T WeightDjikstra, S ~string] struct {
	inititalValue      T
	expectedFinalValue T
	expectedPathHops   int
	expectedPathway    []*NodeDijkstraImpl[T, S]

	edges []*struct {
		to     *NodeDijkstraImpl[T, S]
		weight T
	}
}

// TODO: Add tests for other types

func TestDijkstra(t *testing.T) {
	const (
		nameStart  = "start"
		nameFinish = "finish"
	)
	testDijkstra[uint](t, nameStart, nameFinish)
	testDijkstra[uint8](t, nameStart, nameFinish)
	testDijkstra[int32](t, nameStart, nameFinish)
	testDijkstra[float64](t, nameStart, nameFinish)
}

func constructDijkstraTestGraph[T WeightDjikstra, S ~string](nameStart, nameFinish S) map[NodeDijkstra[T, S]]*dijkstraTestUtils[T, S] {
	// TODO: infinity is way too low, because dijkstraWeight also has uint8
	infinity := T(100)
	nodeStart := &NodeDijkstraImpl[T, S]{
		NodeWeightedImpl: NodeWeightedImpl[T, S]{
			Name:        nameStart,
			ValueOrCost: T(0),
		},
	}
	nodeA := &NodeDijkstraImpl[T, S]{
		NodeWeightedImpl: NodeWeightedImpl[T, S]{
			Name:        "A",
			ValueOrCost: infinity,
		},
	}
	nodeB := &NodeDijkstraImpl[T, S]{
		NodeWeightedImpl: NodeWeightedImpl[T, S]{
			Name:        "B",
			ValueOrCost: infinity,
		},
	}
	nodeC := &NodeDijkstraImpl[T, S]{
		NodeWeightedImpl: NodeWeightedImpl[T, S]{
			Name:        "C",
			ValueOrCost: infinity,
		},
	}
	nodeD := &NodeDijkstraImpl[T, S]{
		NodeWeightedImpl: NodeWeightedImpl[T, S]{
			Name:        "D",
			ValueOrCost: infinity,
		},
	}
	nodeFinish := &NodeDijkstraImpl[T, S]{
		NodeWeightedImpl: NodeWeightedImpl[T, S]{
			Name:        nameFinish,
			ValueOrCost: infinity,
		},
	}
	m := map[NodeDijkstra[T, S]]*dijkstraTestUtils[T, S]{
		nodeStart: {
			inititalValue:      T(0),
			expectedFinalValue: T(0),
			expectedPathHops:   1,
			expectedPathway:    []*NodeDijkstraImpl[T, S]{},
			edges: []*struct {
				to     *NodeDijkstraImpl[T, S]
				weight T
			}{
				{
					to:     nodeA,
					weight: T(2),
				}, {
					to:     nodeB,
					weight: T(4),
				}, {
					to:     nodeD,
					weight: T(4),
				},
			},
		},
		nodeFinish: {
			inititalValue:      infinity,
			expectedFinalValue: T(7),
			expectedPathHops:   3,
			expectedPathway:    []*NodeDijkstraImpl[T, S]{nodeStart, nodeD, nodeFinish},
			edges: []*struct {
				to     *NodeDijkstraImpl[T, S]
				weight T
			}{},
		},
		nodeA: {
			inititalValue:      infinity,
			expectedFinalValue: T(2),
			expectedPathHops:   2,
			expectedPathway:    []*NodeDijkstraImpl[T, S]{nodeStart, nodeA},
			edges: []*struct {
				to     *NodeDijkstraImpl[T, S]
				weight T
			}{
				{
					to:     nodeB,
					weight: T(1),
				},
				{
					to:     nodeC,
					weight: T(2),
				},
			},
		},
		nodeB: {
			inititalValue:      infinity,
			expectedFinalValue: T(3),
			expectedPathHops:   3,
			expectedPathway:    []*NodeDijkstraImpl[T, S]{nodeStart, nodeA, nodeB},
			edges: []*struct {
				to     *NodeDijkstraImpl[T, S]
				weight T
			}{
				{
					to:     nodeFinish,
					weight: T(5),
				},
			},
		},
		nodeC: {
			inititalValue:      infinity,
			expectedFinalValue: T(4),
			expectedPathHops:   3,
			expectedPathway:    []*NodeDijkstraImpl[T, S]{nodeStart, nodeA, nodeC},
			edges: []*struct {
				to     *NodeDijkstraImpl[T, S]
				weight T
			}{
				{
					to:     nodeD,
					weight: T(2),
				},
			},
		},
		nodeD: {
			inititalValue:      infinity,
			expectedFinalValue: T(4),
			expectedPathHops:   2,
			expectedPathway:    []*NodeDijkstraImpl[T, S]{nodeStart, nodeD},
			edges: []*struct {
				to     *NodeDijkstraImpl[T, S]
				weight T
			}{
				{
					to:     nodeFinish,
					weight: T(3),
				},
			},
		},
	}
	return m
}

// The weighted graph used in this test can be viewed at assets/dijkstra_test_graph.png
func testDijkstra[T WeightDjikstra, S ~string](t *testing.T, nameStart, nameFinish S) {
	nodesMap := constructDijkstraTestGraph[T](nameStart, nameFinish)

	// Prepare graph
	directed := true
	djikGraph := NewDijkstraGraph[T, S](directed)
	for node, util := range nodesMap {
		// Add node
		djikGraph.AddNode(node)
		// Add edges
		nodeEdges := util.edges
		for _, edge := range nodeEdges {
			if err := djikGraph.AddEdge(node, edge.to, edge.weight); err != nil {
				t.Error(err.Error())
			}
		}
	}

	var startNode NodeDijkstra[T, S]
	for node := range nodesMap {
		if node.GetKey() == nameStart {
			startNode = node
		}
	}

	dijkShortestPaths := djikGraph.DijkstraShortestPathFrom(startNode)
	fatalMsgCost := "invalid answer for (%s->%s): %d, expecting %d"
	fatalMsgPathLen := "invalid returned path length (%s->%s): %d, expecting %d"
	fatalMsgPathVia := "invalid via path (%s->%s)[%d]: %s, expecting %s"

	// The check is hard-coded
	for node, util := range nodesMap {
		// Test costs
		if node.GetValue() != util.expectedFinalValue {
			t.Fatalf(fatalMsgCost, nameStart, node.GetKey(), node.GetValue(), util.expectedFinalValue)
		}

		if node == startNode {
			continue
		}
		nodeKey := node.GetKey()
		// t.Logf("dst node: %v\n", nodeKey)
		// Test paths
		pathToNode := dijkShortestPaths.ReconstructPathTo(node)
		gslutils.ReverseInPlace(pathToNode)
		if hops := len(pathToNode); hops != util.expectedPathHops {
			t.Fatalf(fatalMsgPathLen, nameStart, nodeKey, hops, util.expectedPathHops)
		}
		for i, actual := range pathToNode {
			expected := util.expectedPathway[i]
			// t.Logf("prev[%d]: %v, expected: %v\n", i, actual.GetKey(), expected.GetKey())
			if expected != actual {
				t.Fatalf(fatalMsgPathVia, nameStart, nodeKey, i, actual.GetKey(), expected.GetKey())
			}
		}
	}
}
