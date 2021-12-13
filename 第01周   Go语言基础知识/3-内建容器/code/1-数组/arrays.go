package main

import "fmt"

func main() {
	var arr1 [5]int
	arr2 := [3]int{3, 4, 5}   //自动推算类型
	arr3 := [...]int{3, 4, 5} //让编译器计算数组长度

	var grid [4][5]int
	fmt.Println(arr1, arr2, arr3)
	fmt.Println(grid)

	for i, v := range arr3 {
		fmt.Println(i, v)
	}
}
