package main

import (
	"fmt"
)

type uint64Slice []uint64

func main()  {
	numbers := []uint64{5,4,20,3,8,2,9}
	insertSort(numbers)
	fmt.Println(numbers)
}

func insertSort(numbers uint64Slice)  {
	for i:=1; i < len(numbers); i++{
		tmp := numbers[i]
		// 从待排序序列开始比较,找到比其小的数
		j:=i
		for j>0 && tmp<numbers[j-1] {
			numbers[j] = numbers[j-1]
			j--
		}
		// 存在比其小的数插入
		if j!=i{
			numbers[j] = tmp
		}
	}
}


// 交换方法
func (numbers uint64Slice)swap(i,j int)  {
	numbers[i],numbers[j] = numbers[j],numbers[i]
}