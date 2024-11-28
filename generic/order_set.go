// 提供非并发安全的可排序 OrderSet
package generic

import "sort"

type OrderSet[T Ordered] struct {
	*Set[T]
}

func NewOrderSet[T Ordered](items ...T) *OrderSet[T] {
	s := &OrderSet[T]{
		Set: NewSet[T](items...),
	}
	s.Add(items...)
	return s
}

// SortItems 返回已排序的 slice 数据
func (s *OrderSet[T]) SortItems() []T {
	slice := s.Items()
	sort.Slice(slice, func(i, j int) bool {
		return slice[i] < slice[j]
	})
	return slice
}
