package hashset

type UintSet struct {
	m map[uint]struct{}
}

func NewUintSet(values ...uint) *UintSet {
	s := &UintSet{
		m: make(map[uint]struct{}, len(values)),
	}
	s.Add(values...)
	return s
}

func (s *UintSet) Add(values ...uint) {
	for _, v := range values {
		s.m[v] = exists
	}
}

func (s *UintSet) Remove(values ...uint) {
	for _, v := range values {
		delete(s.m, v)
	}
}

func (s *UintSet) Contains(value uint) bool {
	_, ok := s.m[value]
	return ok
}

func (s *UintSet) Range(f func(value uint) bool) {
	for k := range s.m {
		if !f(k) {
			break
		}
	}
}

func (s *UintSet) Merge(another *UintSet) {
	another.Range(func(str uint) bool {
		s.Add(str)
		return true
	})
}

func (s *UintSet) GetSlice() []uint {
	slice := make([]uint, 0, len(s.m))
	s.Range(func(value uint) bool {
		slice = append(slice, value)
		return true
	})
	return slice
}

func (s *UintSet) Len() int {
	return len(s.m)
}
