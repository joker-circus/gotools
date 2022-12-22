package types

import (
	"reflect"
	"unsafe"
)

// 将切片、数组转成 []interface。
// 如果类型不是 Array、Slice，程序会 panic。
func SliceInterface(slice interface{}) []interface{} {
	rv := reflect.ValueOf(slice)
	if rv.Kind() != reflect.Slice && rv.Kind() != reflect.Array {
		panic("kind is not Array or Slice")
	}

	data := make([]interface{}, rv.Len(), rv.Len())
	for i := 0; i < rv.Len(); i++ {
		data[i] = rv.Index(i).Interface()
	}
	return data
}

// Chunk 按照一定的间隔分割数组，如果类型不是 Slice 或 gap <= 0，程序会 panic。
// 例如：Chunk([1,2,3], 2)，返回 [[1,2],[3]]。
// 直接 copy 底层数组指针，减少了中间的性能损耗。
func Chunk(slice interface{}, gap int) interface{} {
	if gap <= 0 {
		panic("gap <= 0")
	}

	rv := reflect.Indirect(reflect.ValueOf(slice))
	if rv.Kind() != reflect.Slice {
		panic("kind is not Slice")
	}

	dataLen := rv.Len() / gap
	if rv.Len()%gap != 0 {
		dataLen += 1
	}

	base, size, sliceCap := sliceBase(rv)

	res := reflect.MakeSlice(reflect.SliceOf(rv.Type()), dataLen, dataLen)
	for i := 0; i < dataLen; i++ {
		start := i * gap
		end := (i + 1) * gap
		if end > rv.Len() {
			end = rv.Len()
		}

		// 直接赋值切片地址，对比 res.Index(i).Set(rv.Slice(start, end)) 减少中间 allocs 内存损耗，
		*(*reflect.SliceHeader)(unsafe.Pointer(res.Index(i).UnsafeAddr())) = *(*reflect.SliceHeader)(sliceUnsafePointer(base, size, sliceCap, start, end))
	}

	return res.Interface()
}

// 返回 slice 基础地址，元素的 size 大小，及当前数组的 cap 值。
// copy from reflect.Slice。
func sliceBase(slice reflect.Value) (base unsafe.Pointer, size uintptr, cap int) {
	return unsafe.Pointer(slice.Pointer()), slice.Type().Elem().Size(), slice.Cap()
}

// 返回 slice[i:j] 的地址。
// copy from reflect.Slice。
func sliceUnsafePointer(base unsafe.Pointer, size uintptr, cap, i, j int) unsafe.Pointer {
	// Declare slice so that gc can see the base pointer in it.
	var x []unsafe.Pointer

	s := (*reflect.SliceHeader)(unsafe.Pointer(&x))
	s.Len = j - i
	s.Cap = cap - i
	if cap-i > 0 {
		s.Data = uintptr(base) + uintptr(i)*size
	} else {
		// do not advance pointer, to avoid pointing beyond end of slice
		s.Data = uintptr(base)
	}

	return unsafe.Pointer(&x)
}
