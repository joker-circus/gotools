package hashset

import "sync"

type SafeUint64Set struct {
	sync.RWMutex
	m map[uint64]struct{}
}

func NewSafeUint64Set(values ...uint64) *SafeUint64Set {
	s := &SafeUint64Set{
		m: make(map[uint64]struct{}, len(values)),
	}
	s.Add(values...)
	return s
}

func (s *SafeUint64Set) Add(values ...uint64) {
	s.Lock()
	for _, v := range values {
		s.m[v] = exists
	}
	s.Unlock()
}

func (s *SafeUint64Set) Remove(values ...uint64) {
	s.Lock()
	for _, v := range values {
		delete(s.m, v)
	}
	s.Unlock()
}

func (s *SafeUint64Set) Contains(value uint64) bool {
	s.RLock()
	_, ok := s.m[value]
	s.RUnlock()
	return ok
}

func (s *SafeUint64Set) Range(f func(value uint64) bool) {
	s.RLock()
	for k := range s.m {
		if !f(k) {
			break
		}
	}
	s.RUnlock()
}

func (s *SafeUint64Set) Merge(another *SafeUint64Set) {
	s.Lock()
	another.Range(func(str uint64) bool {
		s.Add(str)
		return true
	})
	s.Unlock()
}

func (s *SafeUint64Set) GetSlice() []uint64 {
	s.RLock()

	slice := make([]uint64, 0, len(s.m))
	s.Range(func(value uint64) bool {
		slice = append(slice, value)
		return true
	})

	s.RUnlock()
	return slice
}

func (s *SafeUint64Set) Len() int {
	s.RLock()
	l := len(s.m)
	s.RUnlock()
	return l
}
