package hashset

import "sync"

type SafeIntSet struct {
	sync.RWMutex
	m map[int]struct{}
}

func NewSafeIntSet(values ...int) *SafeIntSet {
	s := &SafeIntSet{
		m: make(map[int]struct{}, len(values)),
	}
	s.Add(values...)
	return s
}

func (s *SafeIntSet) Add(values ...int) {
	s.Lock()
	for _, v := range values {
		s.m[v] = exists
	}
	s.Unlock()
}

func (s *SafeIntSet) Remove(values ...int) {
	s.Lock()
	for _, v := range values {
		delete(s.m, v)
	}
	s.Unlock()
}

func (s *SafeIntSet) Contains(value int) bool {
	s.RLock()
	_, ok := s.m[value]
	s.RUnlock()
	return ok
}

func (s *SafeIntSet) Range(f func(value int) bool) {
	s.RLock()
	for k := range s.m {
		if !f(k) {
			break
		}
	}
	s.RUnlock()
}

func (s *SafeIntSet) Merge(another *SafeIntSet) {
	s.Lock()
	another.Range(func(str int) bool {
		s.Add(str)
		return true
	})
	s.Unlock()
}

func (s *SafeIntSet) GetSlice() []int {
	s.RLock()

	slice := make([]int, 0, len(s.m))
	s.Range(func(value int) bool {
		slice = append(slice, value)
		return true
	})

	s.RUnlock()
	return slice
}

func (s *SafeIntSet) Len() int {
	s.RLock()
	l := len(s.m)
	s.RUnlock()
	return l
}
