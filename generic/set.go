// 提供非并发安全的 Set
package generic

type EmptyType struct{}

var empty EmptyType

type Set[T comparable] struct {
	m map[T]EmptyType
}

func NewSet[T comparable](items ...T) *Set[T] {
	s := &Set[T]{
		m: make(map[T]EmptyType),
	}
	s.Add(items...)
	return s
}

// Add 添加元素
func (s *Set[T]) Add(items ...T) *Set[T] {
	for _, v := range items {
		s.m[v] = empty
	}
	return s
}

// Remove 删除元素
func (s *Set[T]) Remove(items ...T) *Set[T] {
	for _, v := range items {
		delete(s.m, v)
	}
	return s
}

// Contains 检测是否包含元素，包含任意元素都会返回 true
func (s *Set[T]) Contains(items ...T) bool {
	for _, item := range items {
		if _, exist := s.m[item]; exist {
			return true
		}
	}
	return false
}

// Enumerate 枚举每个元素，如果 f 返回 false 则立即终止枚举
func (s *Set[T]) Enumerate(f func(item T) bool) {
	for key := range s.m {
		if !f(key) {
			break
		}
	}
}

// Merge 合并其他 Set 数据
func (s *Set[T]) Merge(another *Set[T]) {
	another.Enumerate(func(item T) bool {
		s.Add(item)
		return true
	})
}

// Items 返回 slice 数据
func (s *Set[T]) Items() []T {
	slice := make([]T, 0, len(s.m))
	s.Enumerate(func(item T) bool {
		slice = append(slice, item)
		return true
	})
	return slice
}

// Equals 判断两集合是否相等
func (s *Set[T]) Equals(other *Set[T]) bool {
	if other == nil || s.Len() != other.Len() {
		return false
	}

	for k := range s.m {
		if _, ok := other.m[k]; !ok {
			return false
		}
	}
	return true
}

// Len 获取长度
func (s *Set[T]) Len() int {
	return len(s.m)
}
