package hashset

import "sync"

type SafeStringSet struct {
	sync.RWMutex
	m map[string]struct{}
}

func NewSafeStringSet(values ...string) *SafeStringSet {
	s := &SafeStringSet{
		m: make(map[string]struct{}, len(values)),
	}
	s.Add(values...)
	return s
}

func (s *SafeStringSet) Add(values ...string) {
	s.Lock()
	for _, v := range values {
		s.m[v] = exists
	}
	s.Unlock()
}

func (s *SafeStringSet) Remove(values ...string) {
	s.Lock()
	for _, v := range values {
		delete(s.m, v)
	}
	s.Unlock()
}

func (s *SafeStringSet) Contains(value string) bool {
	s.RLock()
	_, ok := s.m[value]
	s.RUnlock()
	return ok
}

func (s *SafeStringSet) Range(f func(value string) bool) {
	s.RLock()
	for k := range s.m {
		if !f(k) {
			break
		}
	}
	s.RUnlock()
}

func (s *SafeStringSet) Merge(another *SafeStringSet) {
	s.Lock()
	another.Range(func(str string) bool {
		s.Add(str)
		return true
	})
	s.Unlock()
}

func (s *SafeStringSet) GetSlice() []string {
	s.RLock()

	slice := make([]string, 0, len(s.m))
	s.Range(func(value string) bool {
		slice = append(slice, value)
		return true
	})

	s.RUnlock()
	return slice
}

func (s *SafeStringSet) Len() int {
	s.RLock()
	l := len(s.m)
	s.RUnlock()
	return l
}
