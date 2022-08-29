package hashset

type Float32Set struct {
	m map[float32]struct{}
}

func NewFloat32Set(values ...float32) *Float32Set {
	s := &Float32Set{
		m: make(map[float32]struct{}, len(values)),
	}
	s.Add(values...)
	return s
}

func (s *Float32Set) Add(values ...float32) {
	for _, v := range values {
		s.m[v] = exists
	}
}

func (s *Float32Set) Remove(values ...float32) {
	for _, v := range values {
		delete(s.m, v)
	}
}

func (s *Float32Set) Contains(value float32) bool {
	_, ok := s.m[value]
	return ok
}

func (s *Float32Set) Range(f func(value float32) bool) {
	for k := range s.m {
		if !f(k) {
			break
		}
	}
}

func (s *Float32Set) Merge(another *Float32Set) {
	another.Range(func(str float32) bool {
		s.Add(str)
		return true
	})
}

func (s *Float32Set) GetSlice() []float32 {
	slice := make([]float32, 0, len(s.m))
	s.Range(func(value float32) bool {
		slice = append(slice, value)
		return true
	})
	return slice
}

func (s *Float32Set) Len() int {
	return len(s.m)
}
