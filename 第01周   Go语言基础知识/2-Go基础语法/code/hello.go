package main

import "fmt"

//函数外变量 不能使用 ：= 必须使用var定义
// var aa = 1
// var bb = true
// var ss = "kkk"

var (
	aa = 1
	bb = true
	ss = "kkk"
)

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

//多个不同类型变量类型推断
func variableTypeDeduction() {
	var a, b, c, s = 3, 4, true, "def"
	fmt.Println(a, b, c, s)
}

func variableShorter() {
	a, b, c, s := 3, 4, true, "def"
	b = 5 //第一次用 ：=   之后赋值用 =
	fmt.Println(a, b, c, s)
}

func main() {
	fmt.Println("Hello World")
	variableZeroValue()
	variableInitialValue()
	variableTypeDeduction()
	variableShorter()
	fmt.Println(aa, bb, ss)
}
