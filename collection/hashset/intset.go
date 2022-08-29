package hashset

type IntSet struct {
	m map[int]struct{}
}

func NewIntSet(values ...int) *IntSet {
	s := &IntSet{
		m: make(map[int]struct{}, len(values)),
	}
	s.Add(values...)
	return s
}

func (s *IntSet) Add(values ...int) {
	for _, v := range values {
		s.m[v] = exists
	}
}

func (s *IntSet) Remove(values ...int) {
	for _, v := range values {
		delete(s.m, v)
	}
}

func (s *IntSet) Contains(value int) bool {
	_, ok := s.m[value]
	return ok
}

func (s *IntSet) Range(f func(value int) bool) {
	for k := range s.m {
		if !f(k) {
			break
		}
	}
}

func (s *IntSet) Merge(another *IntSet) {
	another.Range(func(str int) bool {
		s.Add(str)
		return true
	})
}

func (s *IntSet) GetSlice() []int {
	slice := make([]int, 0, len(s.m))
	s.Range(func(value int) bool {
		slice = append(slice, value)
		return true
	})
	return slice
}

func (s *IntSet) Len() int {
	return len(s.m)
}
