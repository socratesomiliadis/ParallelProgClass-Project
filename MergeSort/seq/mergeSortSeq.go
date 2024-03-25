package main

import (
	"fmt"
	"math/rand"
)

func mergeSort(src []int64) {
	if len(src) <= 1 {
		return
	}

	mid := len(src) / 2

	left := make([]int64, mid)
	right := make([]int64, len(src)-mid)
	copy(left, src[:mid])
	copy(right, src[mid:])

	mergeSort(left)
	mergeSort(right)

	merge(src, left, right)
}

func merge(result, left, right []int64) {
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
}

func main() {
	arr := genArr(100000000)
	failed := false

	mergeSort(arr)

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
