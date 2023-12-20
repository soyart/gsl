package tree

import (
	"golang.org/x/exp/constraints"

	"github.com/soyart/gsl/data"
)

type Heap[T any] struct {
	Items   []data.Getter[T]
	CmpFunc data.LessFunc[data.Getter[T]]
}

func NewHeap[T constraints.Ordered](
	order data.SortOrder,
) *Heap[T] {
	return &Heap[T]{
		CmpFunc: data.FactoryLessFuncOrdered[T](order),
	}
}

func NewHeapCmp[T data.CmpOrdered[T]](
	order data.SortOrder,
) *Heap[T] {
	return &Heap[T]{
		CmpFunc: data.FactoryLessFuncCmp[T](order),
	}
}

func (h *Heap[T]) Push(item T) {
	getter := data.NewGetter[T](item)
	h.PushGetter(getter)
}

func (h *Heap[T]) PushGetter(getter data.Getter[T]) {
	h.Items = append(h.Items, getter)
	h.heapifyUp(h.Len() - 1)
}

func (h *Heap[T]) Pop() *T {
	rootValue := h.Items[0].GetValue()
	lastIdx := h.Len() - 1

	h.Items[0] = h.Items[lastIdx]
	h.Items = h.Items[:lastIdx]

	h.heapifyDown(0)

	return &rootValue
}

func (h *Heap[T]) PopGetter() data.Getter[T] {
	rootNode := h.Items[0]
	lastIdx := h.Len() - 1

	h.Items[0] = h.Items[lastIdx]
	h.Items = h.Items[:lastIdx]

	h.heapifyDown(0)

	return rootNode
}

func (h *Heap[T]) Len() int {
	return len(h.Items)
}

func (h *Heap[T]) IsEmpty() bool {
	return len(h.Items) == 0
}

func (h *Heap[T]) Peek() data.Getter[T] {
	return h.Items[0]
}

func (h *Heap[T]) PopValue() T {
	copied := *h.Pop()

	return copied
}

func (h *Heap[T]) PeekValue() T {
	return h.Peek().GetValue()
}

func (h *Heap[T]) heapifyUp(from int) {
	curr := from
	for curr != 0 {
		parent := ParentIdx(curr)

		if !h.CmpFunc(h.Items, curr, parent) {
			break
		}

		h.swap(curr, parent)
		curr = parent
	}
}

func (h *Heap[T]) heapifyDown(from int) {
	curr := from
	length := len(h.Items)

	for {
		childLeft := LeftChildIdx(curr)
		if childLeft >= length {
			break
		}

		childRight := RightChildIdx(curr)

		// Child to compare
		//nolint:ineffassign
		child := -1

		// Choose left child if:
		// 1) right child is null (out of range)
		// 2) left child has higher priority (lessFunc -> true)
		//
		// Otherwise use right child
		switch {
		case
			childRight >= length,
			h.CmpFunc(h.Items, childLeft, childRight):

			child = childLeft

		default:
			child = childRight
		}

		if h.CmpFunc(h.Items, curr, child) {
			break
		}

		h.swap(curr, child)

		curr = child
	}
}

func (h *Heap[T]) swap(i, j int) {
	h.Items[i], h.Items[j] = h.Items[j], h.Items[i]
}
