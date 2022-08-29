package hashset

import "sync"

type SafeUint16Set struct {
	sync.RWMutex
	m map[uint16]struct{}
}

func NewSafeUint16Set(values ...uint16) *SafeUint16Set {
	s := &SafeUint16Set{
		m: make(map[uint16]struct{}, len(values)),
	}
	s.Add(values...)
	return s
}

func (s *SafeUint16Set) Add(values ...uint16) {
	s.Lock()
	for _, v := range values {
		s.m[v] = exists
	}
	s.Unlock()
}

func (s *SafeUint16Set) Remove(values ...uint16) {
	s.Lock()
	for _, v := range values {
		delete(s.m, v)
	}
	s.Unlock()
}

func (s *SafeUint16Set) Contains(value uint16) bool {
	s.RLock()
	_, ok := s.m[value]
	s.RUnlock()
	return ok
}

func (s *SafeUint16Set) Range(f func(value uint16) bool) {
	s.RLock()
	for k := range s.m {
		if !f(k) {
			break
		}
	}
	s.RUnlock()
}

func (s *SafeUint16Set) Merge(another *SafeUint16Set) {
	s.Lock()
	another.Range(func(str uint16) bool {
		s.Add(str)
		return true
	})
	s.Unlock()
}

func (s *SafeUint16Set) GetSlice() []uint16 {
	s.RLock()

	slice := make([]uint16, 0, len(s.m))
	s.Range(func(value uint16) bool {
		slice = append(slice, value)
		return true
	})

	s.RUnlock()
	return slice
}

func (s *SafeUint16Set) Len() int {
	s.RLock()
	l := len(s.m)
	s.RUnlock()
	return l
}
