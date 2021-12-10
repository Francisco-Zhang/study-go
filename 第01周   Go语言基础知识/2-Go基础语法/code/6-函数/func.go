package main

import "fmt"

//函数可以返回多个值 	13 /3 = 4 ...1
func div(a, b int) (int, int) {
	return a / b, a % b
}

//函数可以返回多个值 	13 /3 = 4 ...1
func div2(a, b int) (q, r int) {
	return a / b, a % b
}

//还可以这样写，直接使用返回值的变量
func div3(a, b int) (q, r int) {
	q = a / b
	r = a % b
	return
}

func main() {
	fmt.Println(div(13, 3))
	q, r := div2(13, 3)
	fmt.Println(q, r)
}
