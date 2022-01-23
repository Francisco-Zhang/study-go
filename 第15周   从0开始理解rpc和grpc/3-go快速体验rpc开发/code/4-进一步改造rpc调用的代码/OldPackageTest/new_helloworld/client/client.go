package main

import (
	"OldPackageTest/new_helloworld/client_proxy"
	"fmt"
)

func main() {
	client := client_proxy.NewHelloServiceClient("tcp", "localhost:1234")
	var reply string
	err := client.Hello("tom", &reply)
	if err != nil {
		panic("调用失败")
	}
	fmt.Println(reply)
}
