package main

import "fmt"

func swap(a, b *int) {
	*a, *b = *b, *a
}

//更简单，更普遍的写法
func swap2(a, b int) (int, int) {
	return b, a
}

func main() {
	a, b := 3, 4
	swap(&a, &b)
	fmt.Println(a, b)

	fmt.Println(swap2(5, 6))
}
