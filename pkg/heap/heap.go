package heap

import (
	hp "container/heap"
)

// 此文件是小顶堆，按照一条路径的权重和排序

type Path struct {
	Value int
	Nodes []string
}

type minPath []Path

func (m minPath) Len() int           { return len(m) }
func (m minPath) Less(i, j int) bool { return m[i].Value < m[j].Value }
func (m minPath) Swap(i, j int)      { m[i], m[j] = m[j], m[i] }

func (m *minPath) Push(x interface{}) {
	*m = append(*m, x.(Path))
}

func (m *minPath) Pop() interface{} {
	old := *m
	n := len(old)
	x := old[n-1]
	*m = old[0 : n-1]
	return x
}

type heap struct {
	Values *minPath
}

func NewHeap() *heap {
	return &heap{Values: &minPath{}}
}

func (h *heap) Push(p Path) {
	hp.Push(h.Values, p)
}

func (h *heap) Pop() Path {
	i := hp.Pop(h.Values)
	return i.(Path)
}
