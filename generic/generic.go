package generic

import (
	"encoding/json"
	"math"
)

func MapWitchError[T any, R any](collection []T, iteratee func(item T, index int) (R, error)) ([]R, error) {
	result := make([]R, len(collection))

	var err error
	for i, item := range collection {
		result[i], err = iteratee(item, i)
		if err != nil {
			return result, err
		}
	}

	return result, nil
}

// 保证 x 是正数，如果x<0则返回0。
func PositiveNumber[T Integer | Float](x T) T {
	if x < 0 {
		return 0
	}
	return x
}

// 去重
func Uniq[T comparable](collection []T) []T {
	result := make([]T, 0, len(collection))
	seen := make(map[T]struct{}, len(collection))

	for _, item := range collection {
		if _, ok := seen[item]; ok {
			continue
		}

		seen[item] = struct{}{}
		result = append(result, item)
	}

	return result
}

// Union returns all distinct elements from given collections.
// result returns will not change the order of elements relatively.
func Union[T comparable](lists ...[]T) []T {
	result := make([]T, 0)
	seen := map[T]struct{}{}

	for _, list := range lists {
		for _, e := range list {
			if _, ok := seen[e]; !ok {
				seen[e] = struct{}{}
				result = append(result, e)
			}
		}
	}

	return result
}

// Filter iterates over elements of collection, returning an array of all elements predicate returns truthy for.
// Play: https://go.dev/play/p/Apjg3WeSi7K
func Filter[V any](collection []V, predicate func(item V, index int) bool) []V {
	result := make([]V, 0, len(collection))

	for i, item := range collection {
		if predicate(item, i) {
			result = append(result, item)
		}
	}

	return result
}

// Map manipulates a slice and transforms it to a slice of another type.
// Play: https://go.dev/play/p/OkPcYAhBo0D
func Map[T any, R any](collection []T, iteratee func(item T, index int) R) []R {
	result := make([]R, len(collection))

	for i, item := range collection {
		result[i] = iteratee(item, i)
	}

	return result
}

// FlatMap manipulates a slice and transforms and flattens it to a slice of another type.
// The transform function can either return a slice or a `nil`, and in the `nil` case
// no value is added to the final slice.
// Play: https://go.dev/play/p/YSoYmQTA8-U
func FlatMap[T any, R any](collection []T, iteratee func(item T, index int) []R) []R {
	result := make([]R, 0, len(collection))

	for i, item := range collection {
		result = append(result, iteratee(item, i)...)
	}

	return result
}

// Difference (news, old []T) (add, del []T) 对比两个数组，返回新增项、删除项、是否相同。
// 非 comparable 类型数组对比使用 DifferenceUncomparable 方法。
// Difference returns the difference between two collections.
// The first value is the collection of element absent of list2.
// The second value is the collection of element absent of list1.
func Difference[T comparable](list1 []T, list2 []T) (left []T, right []T) {
	seenLeft := map[T]struct{}{}
	seenRight := map[T]struct{}{}

	for _, elem := range list1 {
		seenLeft[elem] = struct{}{}
	}

	for _, elem := range list2 {
		seenRight[elem] = struct{}{}
	}

	for _, elem := range list1 {
		if _, ok := seenRight[elem]; !ok {
			left = append(left, elem)
		}
	}

	for _, elem := range list2 {
		if _, ok := seenLeft[elem]; !ok {
			right = append(right, elem)
		}
	}

	return left, right
}

// DifferenceUncomparable (news, old []T, key func(elem T) K, value func(elem T) V) (add, del []T)。
// DifferenceUncomparable 对比两个非 comparable 类型的数组，返回新增项、删除项、是否相同。
// key 是获取每个元素唯一标识的方法，value 是获取每个元素对比内容的方法，通常用 json.Marshal 即可。
// DifferenceUncomparable returns the difference between two collections.
// The first value is the collection of element absent of list2.
// The second value is the collection of element absent of list1.
func DifferenceUncomparable[T any, K, V comparable](list1 []T, list2 []T, key func(elem T) K, value func(elem T) V) (left []T, right []T) {
	comparableList1, seenLeft := ComparableValues(list1, key, value)
	comparableList2, seenRight := ComparableValues(list2, key, value)
	comparableLeft, comparableRight := Difference(comparableList1, comparableList2)

	left = make([]T, 0, len(comparableLeft))
	for _, elem := range comparableLeft {
		left = append(left, seenLeft[elem.Key])
	}
	right = make([]T, 0, len(comparableRight))
	for _, elem := range comparableRight {
		right = append(right, seenRight[elem.Key])
	}
	return
}

// Compact returns a slice of all non-zero elements.
// Play: https://go.dev/play/p/tXiy-iK6PAc
func Compact[T comparable](collection []T) []T {
	var zero T

	result := make([]T, 0, len(collection))

	for _, item := range collection {
		if item != zero {
			result = append(result, item)
		}
	}

	return result
}

// Chunk returns an array of elements split into groups the length of size. If array can't be split evenly,
// the final chunk will be the remaining elements.
// Play: https://go.dev/play/p/EeKl0AuTehH
func Chunk[T any](collection []T, size int) [][]T {
	if size <= 0 {
		panic("Second parameter must be greater than 0")
	}

	chunksNum := len(collection) / size
	if len(collection)%size != 0 {
		chunksNum += 1
	}

	result := make([][]T, 0, chunksNum)

	for i := 0; i < chunksNum; i++ {
		last := (i + 1) * size
		if last > len(collection) {
			last = len(collection)
		}
		result = append(result, collection[i*size:last])
	}

	return result
}

// Contains returns true if an element is present in a collection.
func Contains[T comparable](collection []T, element T) bool {
	for _, item := range collection {
		if item == element {
			return true
		}
	}

	return false
}

// Entry defines a key/value pairs.
type Entry[K comparable, V any] struct {
	Key   K
	Value V
}

// 将不可以比较的结构体转换成可比较的结构体，返回对应的可比较元素、对应的 Map。
// key 是获取每个元素唯一标识的方法，value 是获取每个元素对比内容的方法，通常用 json.Marshal 即可。
func ComparableValues[T any, K, V comparable](data []T, key func(elem T) K, value func(elem T) V) (entries []Entry[K, V], in map[K]T) {
	entries = make([]Entry[K, V], 0, len(data))
	in = make(map[K]T)
	for _, elem := range data {
		k := key(elem)
		entries = append(entries, Entry[K, V]{
			Key:   k,
			Value: value(elem),
		})
		in[k] = elem
	}
	return entries, in
}

// GroupBy returns an object composed of keys generated from the results of running each element of collection through iteratee.
// Play: https://go.dev/play/p/XnQBd_v6brd
func GroupBy[T any, U comparable](collection []T, iteratee func(item T) U) map[U][]T {
	result := map[U][]T{}

	for _, item := range collection {
		key := iteratee(item)

		result[key] = append(result[key], item)
	}

	return result
}

// GroupByWithPoint 返回指针结构体数组，防止申请多余的内存空间。
func GroupByWithPoint[T any, U comparable](collection []T, iteratee func(item T) U) map[U][]*T {
	result := map[U][]*T{}

	for i := range collection {
		item := collection[i]
		key := iteratee(item)

		result[key] = append(result[key], &item)
	}

	return result
}

// 对数据进行分类后，再进行聚合处理，返回聚合后的结果值
func GroupByThenReduce[T any, K comparable, V any](data []T, groupBy func(item T) K, reduce func(key K, values []T) V) []V {
	result := make([]V, 0, len(data))
	for k, values := range GroupBy(data, groupBy) {
		result = append(result, reduce(k, values))
	}
	return result
}

// 对数据进行分类后，再进行聚合处理，返回聚合后的结果值
func GroupByThenMapReduce[T any, K comparable, V any](data []T, groupBy func(item T) K, reduce func(key K, values []T) V) map[K]V {
	result := make(map[K]V)
	for k, values := range GroupBy(data, groupBy) {
		result[k] = reduce(k, values)
	}
	return result
}

// 对数据进行分类后，再进行聚合处理，返回聚合后的结果值。
// 针对特殊情况 key 值不可比较，对 key 值使用 JSON 序列化后聚合。
func GroupByJsonThenReduce[T, K, V any](data []T, groupBy func(item T) K, reduce func(key K, values []T) V) []V {
	// 内部使用 *T 类型，是为了减少产生过多的内存
	keyMap := map[string]*K{}
	prtData := ToSlicePtr(data)
	groupData := GroupBy(prtData, func(item *T) string {
		key := groupBy(*item)
		keyStr := jsonString(key)
		keyMap[keyStr] = &key
		return keyStr
	})

	result := make([]V, 0, len(data))
	for k, prtValues := range groupData {
		values := FromSlicePtr(prtValues)
		result = append(result, reduce(*keyMap[k], values))
	}
	return result
}

// Keys creates an array of the map keys.
// Play: https://go.dev/play/p/Uu11fHASqrU
func Keys[K comparable, V any](in map[K]V) []K {
	result := make([]K, 0, len(in))

	for k := range in {
		result = append(result, k)
	}

	return result
}

// KeyBy transforms a slice or an array of structs to a map based on a pivot callback.
// Play: https://go.dev/play/p/mdaClUAT-zZ
func KeyBy[K comparable, V any](collection []V, iteratee func(item V) K) map[K]V {
	result := make(map[K]V, len(collection))

	for _, v := range collection {
		k := iteratee(v)
		result[k] = v
	}

	return result
}

// Values creates an array of the map values.
// Play: https://go.dev/play/p/nnRTQkzQfF6
func Values[K comparable, V any](in map[K]V) []V {
	result := make([]V, 0, len(in))

	for _, v := range in {
		result = append(result, v)
	}

	return result
}

// Associate returns a map containing key-value pairs provided by transform function applied to elements of the given slice.
// If any of two pairs would have the same key the last one gets added to the map.
// The order of keys in returned map is not specified and is not guaranteed to be the same from the original array.
// Play: https://go.dev/play/p/WHa2CfMO3Lr
func Associate[T any, K comparable, V any](collection []T, transform func(item T) (K, V)) map[K]V {
	result := make(map[K]V, len(collection))

	for _, t := range collection {
		k, v := transform(t)
		result[k] = v
	}

	return result
}

// SliceToMap returns a map containing key-value pairs provided by transform function applied to elements of the given slice.
// If any of two pairs would have the same key the last one gets added to the map.
// The order of keys in returned map is not specified and is not guaranteed to be the same from the original array.
// Alias of Associate().
// Play: https://go.dev/play/p/WHa2CfMO3Lr
func SliceToMap[T any, K comparable, V any](collection []T, transform func(item T) (K, V)) map[K]V {
	return Associate(collection, transform)
}

// Reverse reverses array so that the first element becomes the last, the second element becomes the second to last, and so on.
// Play: https://go.dev/play/p/fhUMLvZ7vS6
func Reverse[T any](collection []T) []T {
	length := len(collection)
	half := length / 2

	for i := 0; i < half; i = i + 1 {
		j := length - 1 - i
		collection[i], collection[j] = collection[j], collection[i]
	}

	return collection
}

// MapToSlice transforms a map into a slice based on specific iteratee
// Play: https://go.dev/play/p/ZuiCZpDt6LD
func MapToSlice[K comparable, V any, R any](in map[K]V, iteratee func(key K, value V) R) []R {
	result := make([]R, 0, len(in))

	for k, v := range in {
		result = append(result, iteratee(k, v))
	}

	return result
}

// 对数组每个 map 固定情况下，提取表的标题 columns 行和每一行值。
func MapToTable[K comparable, V any](data []map[K]V) (columns []K, rows [][]V) {
	if len(data) == 0 {
		return
	}

	columns = make([]K, 0, len(data[0]))
	columnIndex := make(map[K]int)
	for k := range data[0] {
		columnIndex[k] = len(columns)
		columns = append(columns, k)
	}

	rows = make([][]V, 0, len(data))
	for _, iterm := range data {
		row := make([]V, len(columns))
		for k, v := range iterm {
			idx, ok := columnIndex[k]
			if !ok {
				continue
			}
			row[idx] = v
		}
		rows = append(rows, row)
	}
	return columns, rows
}

// ForEach iterates over elements of collection and invokes iteratee for each element.
// Play: https://go.dev/play/p/oofyiUPRf8t
func ForEach[T any](collection []T, iteratee func(item T, index int)) {
	for i, item := range collection {
		iteratee(item, i)
	}
}

// 求众数。摩尔投票法，此方法只适用于一定有结果的情况。
func MajorityElement[T comparable](list []T) (majority T) {
	vote := 0
	for _, num := range list {
		if vote == 0 { //如果票数等于0，则换投票的数字
			majority = num
		}
		if num == majority { //如果遍历的数字与被投票的数字相等，则票数+1
			vote++
		}
		if num != majority { //如果票数不相等，则票数-1
			vote--
		}
	}
	return majority
}



// ToSlicePtr returns a slice of pointer copy of value.
func FromSlicePtr[T any](collection []*T) []T {
	var zero T
	return Map(collection, func(x *T, _ int) T {
		if x == nil {
			return zero
		}
		return *x
	})
}

// ToSlicePtr returns a slice of pointer copy of value.
func ToSlicePtr[T any](collection []T) []*T {
	return Map(collection, func(x T, _ int) *T {
		return &x
	})
}

// Abs returns the absolute value of x.
//
// Special cases are:
//	Abs(±Inf) = +Inf
//	Abs(NaN) = NaN
func Abs[T Integer | Float](x T) T {
	return T(math.Abs(float64(x)))
}

func jsonString(data interface{}) string {
	if v, ok := data.(string); ok {
		return v
	}

	if v, ok := data.([]byte); ok {
		return string(v)
	}

	b, _ := json.Marshal(data)
	return string(b)
}