package hashset

type Uint32Set struct {
	m map[uint32]struct{}
}

func NewUint32Set(values ...uint32) *Uint32Set {
	s := &Uint32Set{
		m: make(map[uint32]struct{}, len(values)),
	}
	s.Add(values...)
	return s
}

func (s *Uint32Set) Add(values ...uint32) {
	for _, v := range values {
		s.m[v] = exists
	}
}

func (s *Uint32Set) Remove(values ...uint32) {
	for _, v := range values {
		delete(s.m, v)
	}
}

func (s *Uint32Set) Contains(value uint32) bool {
	_, ok := s.m[value]
	return ok
}

func (s *Uint32Set) Range(f func(value uint32) bool) {
	for k := range s.m {
		if !f(k) {
			break
		}
	}
}

func (s *Uint32Set) Merge(another *Uint32Set) {
	another.Range(func(str uint32) bool {
		s.Add(str)
		return true
	})
}

func (s *Uint32Set) GetSlice() []uint32 {
	slice := make([]uint32, 0, len(s.m))
	s.Range(func(value uint32) bool {
		slice = append(slice, value)
		return true
	})
	return slice
}

func (s *Uint32Set) Len() int {
	return len(s.m)
}
