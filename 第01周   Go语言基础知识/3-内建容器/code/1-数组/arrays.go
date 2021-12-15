package main

import "fmt"

// []int 代表切片，[5]int 才是数组
func printArray(arr [5]int) {
	for i, v := range arr {
		fmt.Println(i, v)
	}
	arr[0] = 100
}

//改用指针，简化操作
func printArray2(arr *[5]int) { //使用切片，操作arr[0]会更方便。
	for i, v := range arr {
		fmt.Println(i, v)
	}
	arr[0] = 100
}

func main() {
	var arr1 [5]int
	arr2 := [3]int{3, 4, 5}         //自动推算类型
	arr3 := [...]int{3, 4, 5, 7, 8} //让编译器计算数组长度

	var grid [4][5]int
	fmt.Println(arr1, arr2, arr3)
	fmt.Println(grid)

	//原始写法
	for i := 0; i < len(arr3); i++ {
		fmt.Println(i, arr3[i])
	}

	for i, v := range arr3 {
		fmt.Println(i, v)
	}

	printArray(arr3) //传值只能是5个元素的数组，不是5个的被认为是类型不同。
	printArray(arr3) //数组是值类型，函数内部改变数组的值，不会影响外部的数组

	printArray2(&arr3)
	fmt.Println("指针赋值后")
	fmt.Println(arr3)
}
