package mathutil

import (
	"math"
	"sort"
)

// 计算最大，最小，平均，分位值
//
// quantileList：需要求的分位值（如：0.99、0.95、0.90、0.50），计算结果会在 quantiles 中按顺序返回
func QueryQuantile(list []float64, quantileList ...float64) (max float64, min float64, avg float64, quantiles []float64) {
	sort.SliceStable(list, func(i, j int) bool {
		return list[i] < list[j]
	})

	if len(list) == 0 {
		return
	}

	sum := 0.0
	for _, v := range list {
		sum += v
	}

	max = list[len(list)-1]
	min = list[0]
	avg = sum / float64(len(list))

	for _, quantile := range quantileList {
		quantiles = append(quantiles, QuantileAlgorithm(list, quantile))
	}

	return
}

// 求指定分位值，例如：QuantileAlgorithm(list, 0.99)。
// 需要提前将 list 进行排序：
// 	sort.SliceStable(list, func(i, j int) bool {
//		return list[i] < list[j]
//	})
func QuantileAlgorithm(list []float64, quantile float64) float64 {
	index := quantile * float64(len(list))

	var result int
	if math.Ceil(index)-math.Floor(index) == 0 {
		result = int(index)
	} else {
		result = int(math.Ceil(index)) - 1
	}

	if result > len(list)-1 {
		result = len(list) - 1
	}

	return list[result]
}

// 四舍五入。
// dig：保留的小数位。
// 例如： Round(9.3456, 2) = 9.35。
func Round(x float64, dig int) float64 {
	pow := math.Pow10(dig)
	return math.Floor(x*pow+0.5) / pow
}
