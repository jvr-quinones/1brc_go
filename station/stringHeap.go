package station

import "strings"

type StringHeap []string

func (h StringHeap) Len() int {
	return len(h)
}

func (h StringHeap) Less(a, b int) bool {
	return strings.ToLower(h[a]) < strings.ToLower(h[b])
}

func (h *StringHeap) Pop() any {
	lastIdx := len(*h) - 1
	lastVal := (*h)[lastIdx]
	*h = (*h)[0:lastIdx]
	return lastVal
}

func (h *StringHeap) Push(str any) {
	*h = append(*h, str.(string))
}

func (h StringHeap) Swap(a, b int) {
	h[a], h[b] = h[b], h[a]
}
