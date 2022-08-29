package hashset

type Float64Set struct {
	m map[float64]struct{}
}

func NewFloat64Set(values ...float64) *Float64Set {
	s := &Float64Set{
		m: make(map[float64]struct{}, len(values)),
	}
	s.Add(values...)
	return s
}

func (s *Float64Set) Add(values ...float64) {
	for _, v := range values {
		s.m[v] = exists
	}
}

func (s *Float64Set) Remove(values ...float64) {
	for _, v := range values {
		delete(s.m, v)
	}
}

func (s *Float64Set) Contains(value float64) bool {
	_, ok := s.m[value]
	return ok
}

func (s *Float64Set) Range(f func(value float64) bool) {
	for k := range s.m {
		if !f(k) {
			break
		}
	}
}

func (s *Float64Set) Merge(another *Float64Set) {
	another.Range(func(str float64) bool {
		s.Add(str)
		return true
	})
}

func (s *Float64Set) GetSlice() []float64 {
	slice := make([]float64, 0, len(s.m))
	s.Range(func(value float64) bool {
		slice = append(slice, value)
		return true
	})
	return slice
}

func (s *Float64Set) Len() int {
	return len(s.m)
}
