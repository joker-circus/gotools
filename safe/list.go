package safe

import (
	"container/list"
	"sync"
)

type List struct {
	sync.RWMutex
	L *list.List
}

func NewSafeList() *List {
	return &List{L: list.New()}
}

func (l *List) PushFront(v interface{}) *list.Element {
	l.Lock()
	e := l.L.PushFront(v)
	l.Unlock()
	return e
}

func (l *List) PushFrontBatch(vs []interface{}) {
	l.Lock()
	for _, item := range vs {
		l.L.PushFront(item)
	}
	l.Unlock()
}

func (l *List) PopBack() interface{} {
	l.Lock()

	if elem := l.L.Back(); elem != nil {
		item := l.L.Remove(elem)
		l.Unlock()
		return item
	}

	l.Unlock()
	return nil
}

func (l *List) PopBackBy(max int) []interface{} {
	l.Lock()

	count := l.len()
	if count == 0 {
		l.Unlock()
		return []interface{}{}
	}

	if count > max {
		count = max
	}

	items := make([]interface{}, 0, count)
	for i := 0; i < count; i++ {
		item := l.L.Remove(l.L.Back())
		items = append(items, item)
	}

	l.Unlock()
	return items
}

func (l *List) PopBackAll() []interface{} {
	l.Lock()

	count := l.len()
	if count == 0 {
		l.Unlock()
		return []interface{}{}
	}

	items := make([]interface{}, 0, count)
	for i := 0; i < count; i++ {
		item := l.L.Remove(l.L.Back())
		items = append(items, item)
	}

	l.Unlock()
	return items
}

func (l *List) Remove(e *list.Element) interface{} {
	l.Lock()
	defer l.Unlock()
	return l.L.Remove(e)
}

func (l *List) RemoveAll() {
	l.Lock()
	l.L = list.New()
	l.Unlock()
}

func (l *List) FrontAll() []interface{} {
	l.RLock()
	defer l.RUnlock()

	count := l.len()
	if count == 0 {
		return []interface{}{}
	}

	items := make([]interface{}, 0, count)
	for e := l.L.Front(); e != nil; e = e.Next() {
		items = append(items, e.Value)
	}
	return items
}

func (l *List) BackAll() []interface{} {
	l.RLock()
	defer l.RUnlock()

	count := l.len()
	if count == 0 {
		return []interface{}{}
	}

	items := make([]interface{}, 0, count)
	for e := l.L.Back(); e != nil; e = e.Prev() {
		items = append(items, e.Value)
	}
	return items
}

func (l *List) Front() interface{} {
	l.RLock()

	if f := l.L.Front(); f != nil {
		l.RUnlock()
		return f.Value
	}

	l.RUnlock()
	return nil
}

func (l *List) Len() int {
	l.RLock()
	defer l.RUnlock()
	return l.len()
}

func (l *List) len() int {
	return l.L.Len()
}

// SafeList with Limited Size
type ListLimited struct {
	maxSize int
	SL      *List
}

func NewSafeListLimited(maxSize int) *ListLimited {
	return &ListLimited{SL: NewSafeList(), maxSize: maxSize}
}

func (l *ListLimited) PopBack() interface{} {
	return l.SL.PopBack()
}

func (l *ListLimited) PopBackBy(max int) []interface{} {
	return l.SL.PopBackBy(max)
}

func (l *ListLimited) PushFront(v interface{}) bool {
	if l.SL.Len() >= l.maxSize {
		return false
	}

	l.SL.PushFront(v)
	return true
}

// 批次插入数据的数量不能超出链表剩余容量
func (l *ListLimited) PushFrontBatch(vs []interface{}) bool {
	if l.SL.Len()+len(vs) >= l.maxSize {
		return false
	}

	l.SL.PushFrontBatch(vs)
	return true
}

func (l *ListLimited) PushFrontViolently(v interface{}) bool {
	l.SL.PushFront(v)
	if l.SL.Len() > l.maxSize {
		l.SL.PopBack()
	}

	return true
}

func (l *ListLimited) RemoveAll() {
	l.SL.RemoveAll()
}

func (l *ListLimited) Front() interface{} {
	return l.SL.Front()
}

func (l *ListLimited) FrontAll() []interface{} {
	return l.SL.FrontAll()
}

func (l *ListLimited) Len() int {
	return l.SL.Len()
}

// 队列使用过半则认为队列繁忙
func (l *ListLimited) Busy() bool {
	if l.SL.Len() == 0 {
		return false
	}
	return l.maxSize/l.SL.Len() <= 2
}
