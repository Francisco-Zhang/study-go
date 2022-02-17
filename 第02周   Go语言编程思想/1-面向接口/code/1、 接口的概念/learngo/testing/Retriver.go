package testing

type Retriever struct{}

//因为结构体内为空，不需要访问内部属性，所以函数中 Retriever 可以不写变量名
func (Retriever) Get(url string) string {

	return "fake content"
}
