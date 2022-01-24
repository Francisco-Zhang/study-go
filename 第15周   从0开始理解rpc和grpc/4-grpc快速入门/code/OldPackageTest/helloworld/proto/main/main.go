package main

import (
	hello "OldPackageTest/helloworld/proto"
	"fmt"
)
import "github.com/golang/protobuf/proto"

type Hello struct {
	Name    string   `json:"name"`
	Age     int      `json:"age"`
	Courses []string `json:"courses"`
}

func main() {
	req := hello.HelloRequest{
		Name:    "Tom",
		Age:     18,
		Courses: []string{"go", "gin", "微服务"},
	}
	rsp, _ := proto.Marshal(&req)
	fmt.Println(string(rsp))

	//jsonStruct := Hello{Name: "Tom", Age: 18, Courses: []string{"go", "gin", "微服务"}}
	//jsonRsp, _ := json.Marshal(jsonStruct)
	//fmt.Println(string(jsonRsp))

	newReq := hello.HelloRequest{}
	_ = proto.Unmarshal(rsp, &newReq)
	fmt.Println(newReq.Name, newReq.Age, newReq.Courses)
}
