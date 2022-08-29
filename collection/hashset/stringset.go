package hashset

type StringSet struct {
	m map[string]struct{}
}

func NewStringSet(values ...string) *StringSet {
	s := &StringSet{
		m: make(map[string]struct{}, len(values)),
	}
	s.Add(values...)
	return s
}

func (s *StringSet) Add(values ...string) {
	for _, v := range values {
		s.m[v] = exists
	}
}

func (s *StringSet) Remove(values ...string) {
	for _, v := range values {
		delete(s.m, v)
	}
}

func (s *StringSet) Contains(value string) bool {
	_, ok := s.m[value]
	return ok
}

func (s *StringSet) Range(f func(value string) bool) {
	for k := range s.m {
		if !f(k) {
			break
		}
	}
}

func (s *StringSet) Merge(another *StringSet) {
	another.Range(func(str string) bool {
		s.Add(str)
		return true
	})
}

func (s *StringSet) GetSlice() []string {
	slice := make([]string, 0, len(s.m))
	s.Range(func(value string) bool {
		slice = append(slice, value)
		return true
	})
	return slice
}

func (s *StringSet) Len() int {
	return len(s.m)
}
