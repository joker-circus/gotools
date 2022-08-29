package hashset

type Int16Set struct {
	m map[int16]struct{}
}

func NewInt16Set(values ...int16) *Int16Set {
	s := &Int16Set{
		m: make(map[int16]struct{}, len(values)),
	}
	s.Add(values...)
	return s
}

func (s *Int16Set) Add(values ...int16) {
	for _, v := range values {
		s.m[v] = exists
	}
}

func (s *Int16Set) Remove(values ...int16) {
	for _, v := range values {
		delete(s.m, v)
	}
}

func (s *Int16Set) Contains(value int16) bool {
	_, ok := s.m[value]
	return ok
}

func (s *Int16Set) Range(f func(value int16) bool) {
	for k := range s.m {
		if !f(k) {
			break
		}
	}
}

func (s *Int16Set) Merge(another *Int16Set) {
	another.Range(func(str int16) bool {
		s.Add(str)
		return true
	})
}

func (s *Int16Set) GetSlice() []int16 {
	slice := make([]int16, 0, len(s.m))
	s.Range(func(value int16) bool {
		slice = append(slice, value)
		return true
	})
	return slice
}

func (s *Int16Set) Len() int {
	return len(s.m)
}
