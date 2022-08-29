package timeutil

import (
	"sort"
	"time"
)

type TimeArray []string

func (t TimeArray) Len() int {
	return len(t)
}

func (t TimeArray) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

func (t TimeArray) Less(i, j int) bool {
	t1 := t[i]
	t2 := t[j]
	time1, _ := time.Parse(DateFormat, t1)
	time2, _ := time.Parse(DateFormat, t2)

	return time1.Before(time2)
}

func SortTime(input TimeArray) TimeArray {
	sort.Sort(TimeArray(input))
	return input
}
