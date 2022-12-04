package main

import (
	"fmt"
)

type uint64Slice []uint64

func main()  {
	numbers := []uint64{3,44,38,5,47,15,36,26,27,2,46,4,19,50,48}
	res := mergeSort(numbers)
	fmt.Println(res)
}

func mergeSort(numbers uint64Slice) uint64Slice {
	length := len(numbers)
	if length < 2{
		return numbers
	}
	middle := length/2
	left := numbers[0:middle]
	right := numbers[middle:]
	return merge(mergeSort(left),mergeSort(right))
}

func merge(left uint64Slice,right uint64Slice) uint64Slice {
	result := make(uint64Slice,0)
	for len(left) != 0 && len(right) != 0 {
		if left[0] <= right[0] {
			result = append(result, left[0])
			left = left[1:]
		} else {
			result = append(result, right[0])
			right = right[1:]
		}
	}

	for len(left) != 0 {
		result = append(result, left[0])
		left = left[1:]
	}

	for len(right) != 0 {
		result = append(result, right[0])
		right = right[1:]
	}

	return result
}

// 交换方法
func (numbers uint64Slice)swap(i,j int)  {
	numbers[i],numbers[j] = numbers[j],numbers[i]
}