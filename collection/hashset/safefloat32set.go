package hashset

import "sync"

type SafeFloat32Set struct {
	sync.RWMutex
	m map[float32]struct{}
}

func NewSafeFloat32Set(values ...float32) *SafeFloat32Set {
	s := &SafeFloat32Set{
		m: make(map[float32]struct{}, len(values)),
	}
	s.Add(values...)
	return s
}

func (s *SafeFloat32Set) Add(values ...float32) {
	s.Lock()
	for _, v := range values {
		s.m[v] = exists
	}
	s.Unlock()
}

func (s *SafeFloat32Set) Remove(values ...float32) {
	s.Lock()
	for _, v := range values {
		delete(s.m, v)
	}
	s.Unlock()
}

func (s *SafeFloat32Set) Contains(value float32) bool {
	s.RLock()
	_, ok := s.m[value]
	s.RUnlock()
	return ok
}

func (s *SafeFloat32Set) Range(f func(value float32) bool) {
	s.RLock()
	for k := range s.m {
		if !f(k) {
			break
		}
	}
	s.RUnlock()
}

func (s *SafeFloat32Set) Merge(another *SafeFloat32Set) {
	s.Lock()
	another.Range(func(str float32) bool {
		s.Add(str)
		return true
	})
	s.Unlock()
}

func (s *SafeFloat32Set) GetSlice() []float32 {
	s.RLock()

	slice := make([]float32, 0, len(s.m))
	s.Range(func(value float32) bool {
		slice = append(slice, value)
		return true
	})

	s.RUnlock()
	return slice
}

func (s *SafeFloat32Set) Len() int {
	s.RLock()
	l := len(s.m)
	s.RUnlock()
	return l
}
