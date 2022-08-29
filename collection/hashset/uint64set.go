package hashset

type Uint64Set struct {
	m map[uint64]struct{}
}

func NewUint64Set(values ...uint64) *Uint64Set {
	s := &Uint64Set{
		m: make(map[uint64]struct{}, len(values)),
	}
	s.Add(values...)
	return s
}

func (s *Uint64Set) Add(values ...uint64) {
	for _, v := range values {
		s.m[v] = exists
	}
}

func (s *Uint64Set) Remove(values ...uint64) {
	for _, v := range values {
		delete(s.m, v)
	}
}

func (s *Uint64Set) Contains(value uint64) bool {
	_, ok := s.m[value]
	return ok
}

func (s *Uint64Set) Range(f func(value uint64) bool) {
	for k := range s.m {
		if !f(k) {
			break
		}
	}
}

func (s *Uint64Set) Merge(another *Uint64Set) {
	another.Range(func(str uint64) bool {
		s.Add(str)
		return true
	})
}

func (s *Uint64Set) GetSlice() []uint64 {
	slice := make([]uint64, 0, len(s.m))
	s.Range(func(value uint64) bool {
		slice = append(slice, value)
		return true
	})
	return slice
}

func (s *Uint64Set) Len() int {
	return len(s.m)
}
