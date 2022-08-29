package hashset

import "sync"

type SafeInt16Set struct {
	sync.RWMutex
	m map[int16]struct{}
}

func NewSafeInt16Set(values ...int16) *SafeInt16Set {
	s := &SafeInt16Set{
		m: make(map[int16]struct{}, len(values)),
	}
	s.Add(values...)
	return s
}

func (s *SafeInt16Set) Add(values ...int16) {
	s.Lock()
	for _, v := range values {
		s.m[v] = exists
	}
	s.Unlock()
}

func (s *SafeInt16Set) Remove(values ...int16) {
	s.Lock()
	for _, v := range values {
		delete(s.m, v)
	}
	s.Unlock()
}

func (s *SafeInt16Set) Contains(value int16) bool {
	s.RLock()
	_, ok := s.m[value]
	s.RUnlock()
	return ok
}

func (s *SafeInt16Set) Range(f func(value int16) bool) {
	s.RLock()
	for k := range s.m {
		if !f(k) {
			break
		}
	}
	s.RUnlock()
}

func (s *SafeInt16Set) Merge(another *SafeInt16Set) {
	s.Lock()
	another.Range(func(str int16) bool {
		s.Add(str)
		return true
	})
	s.Unlock()
}

func (s *SafeInt16Set) GetSlice() []int16 {
	s.RLock()

	slice := make([]int16, 0, len(s.m))
	s.Range(func(value int16) bool {
		slice = append(slice, value)
		return true
	})

	s.RUnlock()
	return slice
}

func (s *SafeInt16Set) Len() int {
	s.RLock()
	l := len(s.m)
	s.RUnlock()
	return l
}
