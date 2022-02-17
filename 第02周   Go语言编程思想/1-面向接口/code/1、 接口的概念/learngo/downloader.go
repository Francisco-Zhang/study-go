package main

import (
	"fmt"
	"io/ioutil"
	"learngo/infra"
	"net/http"
)

func retrieve(url string) string {
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	bytes, _ := ioutil.ReadAll(resp.Body)
	return string(bytes)
}

func getRetriever() retriever {
	//return testing.Retriever{}
	return infra.Retriever{}
}

//类型应该是不固定的，something can Get
type retriever interface {
	Get(url string) string
}

func main() {
	//调用 retrieve 函数 耦合性比较大
	//fmt.Println(retrieve("https://www.baidu.com"))

	//这种写法，数据类型固定，只能是 infra.Retriever。如果测试人员想更换为自己写的 testing.Retriever。代码改动非常大。

	retriever := getRetriever()
	fmt.Println(retriever.Get("https://www.baidu.com"))
}
