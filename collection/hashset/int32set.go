package hashset

type Int32Set struct {
	m map[int32]struct{}
}

func NewInt32Set(values ...int32) *Int32Set {
	s := &Int32Set{
		m: make(map[int32]struct{}, len(values)),
	}
	s.Add(values...)
	return s
}

func (s *Int32Set) Add(values ...int32) {
	for _, v := range values {
		s.m[v] = exists
	}
}

func (s *Int32Set) Remove(values ...int32) {
	for _, v := range values {
		delete(s.m, v)
	}
}

func (s *Int32Set) Contains(value int32) bool {
	_, ok := s.m[value]
	return ok
}

func (s *Int32Set) Range(f func(value int32) bool) {
	for k := range s.m {
		if !f(k) {
			break
		}
	}
}

func (s *Int32Set) Merge(another *Int32Set) {
	another.Range(func(str int32) bool {
		s.Add(str)
		return true
	})
}

func (s *Int32Set) GetSlice() []int32 {
	slice := make([]int32, 0, len(s.m))
	s.Range(func(value int32) bool {
		slice = append(slice, value)
		return true
	})
	return slice
}

func (s *Int32Set) Len() int {
	return len(s.m)
}
