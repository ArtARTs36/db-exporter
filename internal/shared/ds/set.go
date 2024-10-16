package ds

type Set[T comparable] struct {
	set  map[T]bool
	list []T
}

func NewSet[T comparable](values ...T) *Set[T] {
	set := &Set[T]{
		set:  map[T]bool{},
		list: make([]T, 0, len(values)),
	}

	for _, value := range values {
		set.Add(value)
	}

	return set
}

func (s *Set[T]) Add(val T) {
	_, exists := s.set[val]
	if exists {
		return
	}

	s.list = append(s.list, val)
	s.set[val] = true
}

func (s *Set[T]) List() []T {
	return s.list
}

func (s *Set[T]) Valid() bool {
	return len(s.list) > 0
}

func (s *Set[T]) Has(value T) bool {
	return s.set[value]
}
