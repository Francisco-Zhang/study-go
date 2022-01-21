package main

import (
	"fmt"
	"net/rpc"
)

func main() {
	client, err := rpc.Dial("tcp", "localhost:1234")
	if err != nil {
		panic("连接失败")
	}
	//var reply *string = new(string)
	var reply string
	err = client.Call("HelloService.Hello", "tom", &reply)
	if err != nil {
		panic("调用失败")
	}
	//fmt.Println(*reply)
	fmt.Println(reply)
}
