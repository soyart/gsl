package data

import (
	"golang.org/x/exp/constraints"
)

type SetValuer[T any] interface {
	SetValue(T)
}

type GetValuer[T any] interface {
	GetValue() T
}

type Valuer[T any] interface {
	SetValuer[T]
	GetValuer[T]
}

type Set[T comparable] interface {
	HasDuplicate(x T) bool
}

type wrapper[T any] struct {
	value T
}

func (w *wrapper[T]) SetValue(value T) {
	w.value = value
}

func (w *wrapper[T]) GetValue() T {
	return w.value
}

func NewValuer[T any](t T) SetValuer[T] {
	return &wrapper[T]{value: t}
}

func NewGetValuer[T any](t T) GetValuer[T] {
	return &wrapper[T]{value: t}
}

func MaxValuer[T constraints.Ordered](values []GetValuer[T]) T {
	valuesT := make([]T, len(values))
	for i := range values {
		valuesT[i] = values[i].GetValue()
	}

	max := valuesT[0]
	for i := range valuesT {
		v := valuesT[i]
		if v > max {
			max = v
		}
	}

	return max
}
