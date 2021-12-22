package main

import "fmt"

func main() {
	m := map[string]string{
		"name":    "ccmouse",
		"course":  "golong",
		"site":    "imooc",
		"quality": "notbad",
	}
	m2 := make(map[string]int) //m2 == empty map
	var m3 map[string]int      //m3 == nill

	fmt.Println(m, m2, m3)

	fmt.Println("Traversing map")

	for k, v := range m {
		fmt.Println(k, v)
	}

	fmt.Println("Getting values")
	courseName, ok := m["course"] //key不存在，返回空串,ok==false
	fmt.Println(courseName, ok)
	if cname, ok := m["cname"]; ok {
		fmt.Println(cname)
	} else {
		fmt.Println("key does not exist")
	}

	fmt.Println("Deleting values")
	delete(m, "name")
	name, ok := m["name"]
	fmt.Println(name, ok)

}
