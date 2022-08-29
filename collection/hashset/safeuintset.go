package hashset

import "sync"

type SafeUintSet struct {
	sync.RWMutex
	m map[uint]struct{}
}

func NewSafeUintSet(values ...uint) *SafeUintSet {
	s := &SafeUintSet{
		m: make(map[uint]struct{}, len(values)),
	}
	s.Add(values...)
	return s
}

func (s *SafeUintSet) Add(values ...uint) {
	s.Lock()
	for _, v := range values {
		s.m[v] = exists
	}
	s.Unlock()
}

func (s *SafeUintSet) Remove(values ...uint) {
	s.Lock()
	for _, v := range values {
		delete(s.m, v)
	}
	s.Unlock()
}

func (s *SafeUintSet) Contains(value uint) bool {
	s.RLock()
	_, ok := s.m[value]
	s.RUnlock()
	return ok
}

func (s *SafeUintSet) Range(f func(value uint) bool) {
	s.RLock()
	for k := range s.m {
		if !f(k) {
			break
		}
	}
	s.RUnlock()
}

func (s *SafeUintSet) Merge(another *SafeUintSet) {
	s.Lock()
	another.Range(func(str uint) bool {
		s.Add(str)
		return true
	})
	s.Unlock()
}

func (s *SafeUintSet) GetSlice() []uint {
	s.RLock()

	slice := make([]uint, 0, len(s.m))
	s.Range(func(value uint) bool {
		slice = append(slice, value)
		return true
	})

	s.RUnlock()
	return slice
}

func (s *SafeUintSet) Len() int {
	s.RLock()
	l := len(s.m)
	s.RUnlock()
	return l
}
