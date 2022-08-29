package hashset

type Uint16Set struct {
	m map[uint16]struct{}
}

func NewUint16Set(values ...uint16) *Uint16Set {
	s := &Uint16Set{
		m: make(map[uint16]struct{}, len(values)),
	}
	s.Add(values...)
	return s
}

func (s *Uint16Set) Add(values ...uint16) {
	for _, v := range values {
		s.m[v] = exists
	}
}

func (s *Uint16Set) Remove(values ...uint16) {
	for _, v := range values {
		delete(s.m, v)
	}
}

func (s *Uint16Set) Contains(value uint16) bool {
	_, ok := s.m[value]
	return ok
}

func (s *Uint16Set) Range(f func(value uint16) bool) {
	for k := range s.m {
		if !f(k) {
			break
		}
	}
}

func (s *Uint16Set) Merge(another *Uint16Set) {
	another.Range(func(str uint16) bool {
		s.Add(str)
		return true
	})
}

func (s *Uint16Set) GetSlice() []uint16 {
	slice := make([]uint16, 0, len(s.m))
	s.Range(func(value uint16) bool {
		slice = append(slice, value)
		return true
	})
	return slice
}

func (s *Uint16Set) Len() int {
	return len(s.m)
}
