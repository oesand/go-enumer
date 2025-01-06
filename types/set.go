package types

type Set[K comparable] map[K]struct{}

func (m *Set[K]) create() {
	if *m == nil {
		*m = make(Set[K])
	}
}

func (m *Set[K]) Add(key K) {
	m.create()
	(*m)[key] = struct{}{}
}

func (m *Set[K]) Contains(key K) bool {
	if *m == nil {
		return false
	}
	_, has := (*m)[key]
	return has
}

func (m *Set[K]) CopyFrom(other Set[K]) {
	if other == nil {
		return
	}
	m.create()
	for k := range other {
		m.Add(k)
	}
}

func (m *Set[K]) Values() []K {
	if *m == nil {
		return nil
	}
	keys := make([]K, 0, len(*m))
	for k := range *m {
		keys = append(keys, k)
	}
	return keys
}
