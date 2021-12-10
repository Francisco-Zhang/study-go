package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func sum() {
	sum := 0
	for i := 1; i <= 100; i++ {
		sum += i
	}
}

//省略初始条件
func convertToBin(n int) string {
	result := ""
	for ; n > 0; n /= 2 {
		lsb := n % 2
		result = strconv.Itoa(lsb) + result
	}
	return result
}

//省略递增条件
func readfile(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() { //相当于其他语言的while，go中没有while
		fmt.Println(scanner.Text())
	}

}

//省略结束条件
func forever() {
	for {
		fmt.Println("abc")
	}
}

func main() {
	fmt.Println(convertToBin(5), //101
		convertToBin(13))
	readfile("abc.txt")
	forever()
}
