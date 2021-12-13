package main

import (
	"fmt"
	"math"
	"reflect"
	"runtime"
)

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

func eval(a, b int, op string) int {
	switch op {
	case "+":
		return a + b
	case "-":
		return a - b
	case "*":
		return a * b
	case "/":
		q, _ := div2(13, 3) //多个返回值只使用一个，变量用_
		return q
	default:
		panic("unsurported operator:" + op)
	}
}

func eval2(a, b int, op string) (int, error) {
	switch op {
	case "+":
		return a + b, nil
	case "-":
		return a - b, nil
	case "*":
		return a * b, nil
	case "/":
		q, _ := div2(13, 3) //多个返回值只使用一个，变量用_
		return q, nil
	default:
		return 0, fmt.Errorf("unsurported operator: %s", op)
	}
}

//体现函数是一等公民
func apply(op func(int, int) int, a, b int) int {
	p := reflect.ValueOf(op).Pointer()
	opName := runtime.FuncForPC(p).Name()
	fmt.Printf("Calling function %s with args "+
		"(%d,%d)\n", opName, a, b)
	return op(a, b)
}

func pow(a, b int) int {
	return int(math.Pow(float64(a), float64(b)))
}

func sum(numbers ...int) int {
	s := 0
	for i := range numbers {
		s += numbers[i]
	}
	return s
}

func main() {
	fmt.Println(div(13, 3))
	q, r := div2(13, 3)
	fmt.Println(q, r)

	fmt.Println(eval2(1, 2, "n"))

	if result, err := eval2(1, 2, "n"); err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println(result)
	}

	fmt.Println(apply(func(a, b int) int {
		return int(math.Pow(float64(a), float64(b)))
	}, 3, 4))

	fmt.Println(sum(1, 2, 3))
}
