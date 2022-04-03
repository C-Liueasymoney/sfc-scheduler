package graph

import (
	"SFC-Scheduler/pkg/heap"
)

type Node struct {
	name string
}

type Edge struct {
	from   string
	to     string
	weight int
}

type Graph struct {
	nodes map[string][]Edge
}

func NewGraph() *Graph {
	return &Graph{nodes: make(map[string][]Edge)}
}

func (g *Graph) AddEdge(from, to string, weight int) {
	g.nodes[from] = append(g.nodes[from], Edge{from: from, to: to, weight: weight})
	g.nodes[to] = append(g.nodes[to], Edge{from: to, to: from, weight: weight})
}

func (g *Graph) GetEdge(nodeName string) []Edge {
	return g.nodes[nodeName]
}

func (g *Graph) GetPath(origin, target string) (int, []string) {
	h := heap.NewHeap()
	// 先在堆中压入源节点，权重是0
	h.Push(heap.Path{Value: 0, Nodes: []string{origin}})
	visited := make(map[string]bool)

	// 实际是使用BFS寻找源节点到目的节点权重最小路径,利用小顶堆的特性使堆顶取出的路径一直是权重和最小的，直到到达目的节点
	for len(*h.Values) > 0 {
		p := h.Pop()
		// 拿到路径中最后一个节点
		node := p.Nodes[len(p.Nodes)-1]

		if visited[node] {
			continue
		}

		if node == target {
			return p.Value, p.Nodes
		}

		// 遍历节点发散出去的所有边
		for _, e := range g.GetEdge(node) {
			if !visited[e.to] {
				h.Push(heap.Path{Value: p.Value + e.weight, Nodes: append([]string{}, append(p.Nodes, e.to)...)})
			}
		}

		visited[node] = true
	}

	return 0, nil
}

//func main() {
//	g := NewGraph()
//	g.AddEdge("Bg", "G", 15)
//	g.AddEdge("G", "Br", 12)
//	g.AddEdge("G", "A", 6)
//	g.AddEdge("Bg", "A", 30)
//	g.AddEdge("Bg", "Br", 10)
//
//	fmt.Println(g.nodes)
//
//	v, p := g.GetPath("Bg", "A")
//	fmt.Println(v)
//	fmt.Println(p)
//	//e := g.GetEdge("node2")
//	//fmt.Println(e)
//
//}
