package utils

import "container/heap"

type sortedHeapItem struct {
	Data  interface{}
	Score float64
}

type SortedHeap interface {
	Len() int
	Add(item interface{}, score float64)
	PopMax() (interface{}, float64)
}

type sortedHeap struct {
	items []sortedHeapItem
}

func (h *sortedHeap) Len() int {
	return len(h.items)
}

func (h *sortedHeap) Less(i, j int) bool {
	return h.items[i].Score > h.items[j].Score
}

func (h *sortedHeap) Swap(i, j int) {
	h.items[i], h.items[j] = h.items[j], h.items[i]
}

func (h *sortedHeap) Push(x interface{}) {
	h.items = append(h.items, x.(sortedHeapItem))
}

func (h *sortedHeap) Pop() interface{} {
	n := len(h.items)
	x := h.items[n-1]
	h.items = h.items[:n-1]
	return x
}

func (h *sortedHeap) Add(item interface{}, score float64) {
	heap.Push(h, sortedHeapItem{item, score})
}

func (h *sortedHeap) PopMax() (interface{}, float64) {
	if h.Len() == 0 {
		return nil, 0.0
	}

	item := heap.Pop(h).(sortedHeapItem)
	return item.Data, item.Score
}
func NewSortedHeap(size int) SortedHeap {
	h := &sortedHeap{
		items: make([]sortedHeapItem, 0, size),
	}
	heap.Init(h)

	return h
}
