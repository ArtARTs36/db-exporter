package shared

type Set struct {
	set  map[string]bool
	list []string
}

func NewSet() *Set {
	return &Set{
		set: map[string]bool{},
	}
}

func (s *Set) Add(val string) {
	_, exists := s.set[val]
	if exists {
		return
	}

	s.list = append(s.list, val)
	s.set[val] = true
}

func (s *Set) List() []string {
	return s.list
}

func (s *Set) Valid() bool {
	return len(s.list) > 0
}
