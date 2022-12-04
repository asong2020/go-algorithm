package main

import (
	"fmt"
)

type uint64Slice []uint64

func main()  {

	numbers := []uint64{5,4,2,3,8}
	sortBubble(numbers)
	fmt.Println(numbers)
}

func sortBubble(numbers uint64Slice)  {
	length := len(numbers)
	if length == 0{
		return
	}
	flag := true

	for i:=0;i<length && flag;i++{
		flag = false
		for j:=length-1;j>i;j--{
			if numbers[j-1] > numbers[j] {
				numbers.swap(j-1,j)
				flag = true // 有数据才交换
			}
		}
	}
}

// 交换方法
func (numbers uint64Slice)swap(i,j int)  {
	numbers[i],numbers[j] = numbers[j],numbers[i]
}