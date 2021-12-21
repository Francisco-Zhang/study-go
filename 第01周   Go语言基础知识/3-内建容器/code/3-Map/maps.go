package main

import "fmt"

func main() {
	m := map[string]string{
		"name":    "ccmouse",
		"course":  "golong",
		"site":    "imooc",
		"quality": "notbad",
	}
	m2 := make(map[string]int)

	fmt.Println(m, m2)
}
