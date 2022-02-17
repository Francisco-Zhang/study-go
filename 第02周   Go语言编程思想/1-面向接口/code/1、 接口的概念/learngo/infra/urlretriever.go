package infra

import (
	"io/ioutil"
	"net/http"
)

type Retriever struct{}

//因为结构体内为空，不需要访问内部属性，所以函数中 Retriever 可以不写变量名
func (Retriever) Get(url string) string {
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	bytes, _ := ioutil.ReadAll(resp.Body)
	return string(bytes)
}
