package main

import (
	"fmt"
)

func main()  {
	numbers := []uint64{3221, 1, 10, 9680, 577, 9420, 7, 5622, 4793, 2030, 3138, 82, 2599, 743, 4127}
	radixSort(numbers)
	fmt.Println(numbers)
}

func radixSort(numbers []uint64)  {
	key := maxDigits(numbers)
	tmp := make([]uint64,len(numbers),len(numbers))
	count := new([10]uint64)
	length := uint64(len(numbers))
	var radix uint64 =  1
	var i, j, k uint64
	for i = 0; i < key; i++ { //进行key次排序
		for j = 0; j < 10; j++ {
			count[j] = 0
		}
		for j = 0; j < length; j++ {
			k = (numbers[j] / radix) % 10
			count[k]++
		}
		for j = 1; j < 10; j++ { //将tmp中的为准依次分配给每个桶
			count[j] = count[j-1] + count[j]
		}
		for j = length-1; j > 0; j-- {
			k = (numbers[j] / radix) % 10
			tmp[count[k]-1] = numbers[j]
			count[k]--
		}
		for j = 0; j < length; j++ {
			numbers[j] = tmp[j]
		}
		radix = radix * 10
	}
}


//获取数组的最大值的位数
func maxDigits(arr []uint64) (ret uint64) {
	ret = 1
	var key uint64 = 10
	for i := 0; i < len(arr); i++ {
		for arr[i] >= key {
			key = key * 10
			ret++
		}
	}
	return
}
