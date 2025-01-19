package types

import "iter"

type Pair[T1 any, T2 any] struct {
	First  T1
	Second T2
}

func (t Pair[T1, T2]) U() (T1, T2) {
	return t.First, t.Second
}

type PairSlice[T1 any, T2 any] []Pair[T1, T2]

func (sl PairSlice[T1, T2]) I() iter.Seq2[T1, T2] {
	return func(yield func(T1, T2) bool) {
		for _, pair := range sl {
			if !yield(pair.U()) {
				return
			}
		}
	}
}
