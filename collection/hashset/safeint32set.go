package hashset

import "sync"

type SafeInt32Set struct {
	sync.RWMutex
	m map[int32]struct{}
}

func NewSafeInt32Set(values ...int32) *SafeInt32Set {
	s := &SafeInt32Set{
		m: make(map[int32]struct{}, len(values)),
	}
	s.Add(values...)
	return s
}

func (s *SafeInt32Set) Add(values ...int32) {
	s.Lock()
	for _, v := range values {
		s.m[v] = exists
	}
	s.Unlock()
}

func (s *SafeInt32Set) Remove(values ...int32) {
	s.Lock()
	for _, v := range values {
		delete(s.m, v)
	}
	s.Unlock()
}

func (s *SafeInt32Set) Contains(value int32) bool {
	s.RLock()
	_, ok := s.m[value]
	s.RUnlock()
	return ok
}

func (s *SafeInt32Set) Range(f func(value int32) bool) {
	s.RLock()
	for k := range s.m {
		if !f(k) {
			break
		}
	}
	s.RUnlock()
}

func (s *SafeInt32Set) Merge(another *SafeInt32Set) {
	s.Lock()
	another.Range(func(str int32) bool {
		s.Add(str)
		return true
	})
	s.Unlock()
}

func (s *SafeInt32Set) GetSlice() []int32 {
	s.RLock()

	slice := make([]int32, 0, len(s.m))
	s.Range(func(value int32) bool {
		slice = append(slice, value)
		return true
	})

	s.RUnlock()
	return slice
}

func (s *SafeInt32Set) Len() int {
	s.RLock()
	l := len(s.m)
	s.RUnlock()
	return l
}
