/* 泛型约束集 golang.org/x/exp/constraints */
package generic

// Ordered 代表所有可比大小排序的类型。
// comparable 比较容易引起误解的一点是与可排序搞混淆。
// 可比较指的是 可以执行 != == 操作的类型，并没确保这个类型可以执行大小比较（ >,<,<=,>= ）
type Ordered interface {
	Integer | Float | ~string
}

// 整数
type Integer interface {
	Signed | Unsigned
}

// 有符号位数字
type Signed interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

// 无符号位数字
type Unsigned interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

// 浮点数
type Float interface {
	~float32 | ~float64
}
