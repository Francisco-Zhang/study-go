package main

import (
	"fmt"
	"math"
	"math/cmplx"
)

func euler() {
	c := 3 + 4i //4i 会被识别为复数，4*i会把i识别为变量
	fmt.Println(cmplx.Abs(c))
	//欧拉公式, 1i可以防止i被当作变量
	fmt.Println(cmplx.Pow(math.E, 1i*math.Pi) + 1)
	//更简单的写法，Exp直接代表e的指数
	fmt.Println(cmplx.Exp(1i*math.Pi) + 1)
	//float精度不够，只取小数点后三位，就可以输出0
	fmt.Printf("%.3f\n", cmplx.Exp(1i*math.Pi)+1)
}

func triangle() {
	var a, b int = 3, 4
	var c int
	//float浮点数在任何语言都是不准的。有可能 math.sqrt算出的结果是 4.9999,所以严谨一点写法是向上取整
	c = int(math.Sqrt(float64(a*a + b*b)))
	fmt.Println(c)
}

func main() {
	euler()
	triangle()
}
