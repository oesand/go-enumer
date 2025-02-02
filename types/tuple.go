package types

type Pair[T1 any, T2 any] struct {
	First  T1
	Second T2
}

func (p Pair[T1, T2]) U() (T1, T2) {
	return p.First, p.Second
}

type Triple[T1 any, T2 any, T3 any] struct {
	First  T1
	Second T2
	Third  T3
}

func (t Triple[T1, T2, T3]) U() (T1, T2, T3) {
	return t.First, t.Second, t.Third
}
