package main

import (
	"fmt"
)

type uint64Slice []uint64

func main()  {
	numbers := []uint64{2, 3, 8, 7, 1, 2, 2, 2, 7, 3, 9, 8, 2}

	countSort(numbers,getMaxValue(numbers))
	fmt.Println(numbers)
}

func countSort(numbers uint64Slice,maxValue uint64) {
	bucketLen := maxValue + 1
	bucket := make(uint64Slice,bucketLen) // 初始都是0的数组
	sortedIndex := 0

	for _,v:= range numbers{
		bucket[v] +=1
	}
	var j uint64
	for j=0;j<bucketLen;j++{
		for bucket[j]>0{
			numbers[sortedIndex] = j
			sortedIndex +=1
			bucket[j] -= 1
		}
	}
}


func getMaxValue(numbers uint64Slice) uint64{
   maxValue := numbers[0]
   for _,v:=range numbers {
	   if maxValue < v {
		   maxValue = v
	   }
   }
   	return maxValue
}