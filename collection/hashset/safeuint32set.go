package hashset

import "sync"

type SafeUint32Set struct {
	sync.RWMutex
	m map[uint32]struct{}
}

func NewSafeUint32Set(values ...uint32) *SafeUint32Set {
	s := &SafeUint32Set{
		m: make(map[uint32]struct{}, len(values)),
	}
	s.Add(values...)
	return s
}

func (s *SafeUint32Set) Add(values ...uint32) {
	s.Lock()
	for _, v := range values {
		s.m[v] = exists
	}
	s.Unlock()
}

func (s *SafeUint32Set) Remove(values ...uint32) {
	s.Lock()
	for _, v := range values {
		delete(s.m, v)
	}
	s.Unlock()
}

func (s *SafeUint32Set) Contains(value uint32) bool {
	s.RLock()
	_, ok := s.m[value]
	s.RUnlock()
	return ok
}

func (s *SafeUint32Set) Range(f func(value uint32) bool) {
	s.RLock()
	for k := range s.m {
		if !f(k) {
			break
		}
	}
	s.RUnlock()
}

func (s *SafeUint32Set) Merge(another *SafeUint32Set) {
	s.Lock()
	another.Range(func(str uint32) bool {
		s.Add(str)
		return true
	})
	s.Unlock()
}

func (s *SafeUint32Set) GetSlice() []uint32 {
	s.RLock()

	slice := make([]uint32, 0, len(s.m))
	s.Range(func(value uint32) bool {
		slice = append(slice, value)
		return true
	})

	s.RUnlock()
	return slice
}

func (s *SafeUint32Set) Len() int {
	s.RLock()
	l := len(s.m)
	s.RUnlock()
	return l
}
