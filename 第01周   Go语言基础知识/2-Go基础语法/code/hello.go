package main

import "fmt"

//变量定义
func variableZeroValue() {
	var a int
	var b string
	fmt.Println(a, b)
	fmt.Printf("%d %q\n", a, b) //打印空串标记
}

func variableInitialValue() {
	var a, b int = 3, 4 //定义多个变量
	var s string = "abc"
	fmt.Println(a, b, s)
}

func main() {
	fmt.Println("Hello World")
	variableZeroValue()
	variableInitialValue()
}
