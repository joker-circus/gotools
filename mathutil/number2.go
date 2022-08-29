package mathutil

import (
	"math/bits"
)

// 获取数值是 2 的几次方，
// 等同于 math.Log2() 取整
func Log2(x int) int {
	return bits.Len(uint(x)) - 1
}

// 是否是 2 的次方
func IsPowerOfTwo(x int) bool {
	return (x & (-x)) == x
}
