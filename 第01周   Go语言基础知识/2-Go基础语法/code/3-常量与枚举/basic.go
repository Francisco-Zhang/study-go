package main

import (
	"fmt"
	"math"
)

func consts() {
	const filename = "abc.txt"
	const a, b = 3, 4
	var c int
	c = int(math.Sqrt(a*a + b*b)) //const 可以看作单纯的文本替换，此处不需要做类型转换。如果 声明了类型，此处才需要做类型转换。
	c = int(math.Sqrt(3*3 + 4*4)) //在编译阶段会自动确定类型
	fmt.Println(c, filename)
}

func consts1() {
	const (
		filename = "abc.txt"
		a, b     = 3, 4
	)
	var c int
	c = int(math.Sqrt(a*a + b*b))
	fmt.Println(c, filename)
}

// func enums() {
// 	const (
// 		cpp    = 0
// 		java   = 1
// 		python = 2
// 		golong = 3
// 	)
// 	fmt.Println(cpp, java, python, golong)
// }
// func enums() {
// 	const (
// 		cpp = iota
// 		java
// 		python
// 		golong
// 	)
// 	fmt.Println(cpp, java, python, golong)
// }

// func enums() {
// 	const (
// 		cpp = iota
// 		_
// 		python
// 		golong
// 		javascript
// 	)
// 	fmt.Println(cpp, javascript, python, golong)
// }

//iota 参与运算
func enums() {
	//b,kb,mb,gb,tb,pb
	const (
		b = 1 << (10 * iota)
		kb
		mb
		gb
		tb
		pb
	)
	fmt.Println(b, kb, mb, gb, tb, pb)
}

func main() {
	enums()
}
