package hashset

import "sync"

type SafeFloat64Set struct {
	sync.RWMutex
	m map[float64]struct{}
}

func NewSafeFloat64Set(values ...float64) *SafeFloat64Set {
	s := &SafeFloat64Set{
		m: make(map[float64]struct{}, len(values)),
	}
	s.Add(values...)
	return s
}

func (s *SafeFloat64Set) Add(values ...float64) {
	s.Lock()
	for _, v := range values {
		s.m[v] = exists
	}
	s.Unlock()
}

func (s *SafeFloat64Set) Remove(values ...float64) {
	s.Lock()
	for _, v := range values {
		delete(s.m, v)
	}
	s.Unlock()
}

func (s *SafeFloat64Set) Contains(value float64) bool {
	s.RLock()
	_, ok := s.m[value]
	s.RUnlock()
	return ok
}

func (s *SafeFloat64Set) Range(f func(value float64) bool) {
	s.RLock()
	for k := range s.m {
		if !f(k) {
			break
		}
	}
	s.RUnlock()
}

func (s *SafeFloat64Set) Merge(another *SafeFloat64Set) {
	s.Lock()
	another.Range(func(str float64) bool {
		s.Add(str)
		return true
	})
	s.Unlock()
}

func (s *SafeFloat64Set) GetSlice() []float64 {
	s.RLock()

	slice := make([]float64, 0, len(s.m))
	s.Range(func(value float64) bool {
		slice = append(slice, value)
		return true
	})

	s.RUnlock()
	return slice
}

func (s *SafeFloat64Set) Len() int {
	s.RLock()
	l := len(s.m)
	s.RUnlock()
	return l
}
