package main

import (
	"fmt"
	"math/rand"
	"sort"
)

func mergeSort(src, temp []int64) {
	if len(src) <= 10000 {
		sort.Slice(src, func(i int, j int) bool { return src[i] <= src[j] })
		return
	}

	mid := len(src) / 2

	left, lTemp := src[:mid], temp[:mid]
	right, rTemp := src[mid:], temp[mid:]

	var done chan bool

	done = make(chan bool)
	go func() {
		mergeSort(left, lTemp)
		done <- true
	}()

	mergeSort(right, rTemp)

	<-done

	merge(src, temp, left, right)
}

func merge(src, result, left, right []int64) {
	var l, r, i int

	for l < len(left) || r < len(right) {
		if l < len(left) && r < len(right) {
			if left[l] <= right[r] {
				result[i] = left[l]
				l++
			} else {
				result[i] = right[r]
				r++
			}
		} else if l < len(left) {
			result[i] = left[l]
			l++
		} else if r < len(right) {
			result[i] = right[r]
			r++
		}
		i++
	}
	copy(src, result)
}

func main() {
	arr := genArr(100000000)
	temp := make([]int64, len(arr))
	failed := false

	mergeSort(arr, temp)

	for i := 0; i < len(arr)-1; i++ {
		if arr[i] >= arr[i+1] {
			failed = true
		}
	}

	if failed {
		fmt.Println("Sort Failed")
	} else {
		fmt.Println("Sort Passed")
	}
}

func genArr(n int) []int64 {
	arr := make([]int64, n)
	for i := range arr {
		arr[i] = rand.Int63()
	}
	return arr
}
