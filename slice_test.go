package gotools

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var benchmarkData = make([]string, 100000)

func init() {
	for i := 0; i < 100000; i++ {
		benchmarkData[i] = fmt.Sprintf("%d", i+1)
	}
}

func TestSliceInterface(t *testing.T) {
	data := make([]string, 0)
	data = append(data, "1", "2", "3", "4", "5")
	for i, v := range SliceInterface(data) {
		assert.Equal(t, v, data[i], "the should be equal")
	}
}

// 按照一定的间隔分割数组
func SplitString(a []string, gap int) [][]string {
	dataLen := len(a) / gap
	if len(a)%gap != 0 {
		dataLen += 1
	}

	res := make([][]string, 0, dataLen)
	for i := 0; i < dataLen; i++ {
		start := i * gap
		end := (i + 1) * gap
		if end > len(a) {
			end = len(a)
		}
		res = append(res, a[start:end])
	}
	return res
}

func TestSplitSlice(t *testing.T) {
	data := make([]string, 16)
	for i := 0; i < 16; i++ {
		data[i] = fmt.Sprintf("%d", i+1)
	}
	v := SplitString(data, 4)
	v2 := Chunk(data, 4).([][]string)
	data[0], data[8] = "100", "100"

	assert.Equal(t, v, v2, "the should be equal")
	fmt.Println(v, len(v), cap(v))
	fmt.Println(v2, len(v2), cap(v2))
}

func BenchmarkSplitString(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = SplitString(benchmarkData, 4)
	}
}

func BenchmarkChunk(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = Chunk(benchmarkData, 4).([][]string)
	}
}
