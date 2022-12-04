package main

import (
	"fmt"
)

func main()  {
	numbers := []uint64{5,3,4,7,4,3,4,7}
	sortBucket(numbers)
	fmt.Println(numbers)
}

func sortBucket(numbers []uint64) {
	num := len(numbers) // 桶数量
	max := getMaxValue(numbers)
	buckets := make([][]uint64,num)
	var index uint64
	for _,v := range numbers{
		// 分配桶 index = value * (n-1)/k
		index = v * uint64(num-1) / max

		buckets[index] = append(buckets[index],v)
	}

	// 桶内排序
	tmpPos := 0
	for k:=0; k < num; k++ {
		bucketLen := len(buckets[k])
		if bucketLen>0{
			sortUseInsert(buckets[k])
			copy(numbers[tmpPos:],buckets[k])
			tmpPos +=bucketLen
		}
	}
}

func sortUseInsert(bucket []uint64)  {
	length := len(bucket)
	if length == 1 {return}
	for i := 1; i < length; i++ {
		backup := bucket[i]
		j := i -1
		for  j >= 0 && backup < bucket[j] {
			bucket[j+1] = bucket[j]
			j --
		}
		bucket[j + 1] = backup
	}
}

//获取数组最大值
func getMaxValue(numbers []uint64) uint64{
	max := numbers[0]
	for i := 1; i < len(numbers); i++ {
		if numbers[i] > max{ max = numbers[i]}
	}
	return max
}
