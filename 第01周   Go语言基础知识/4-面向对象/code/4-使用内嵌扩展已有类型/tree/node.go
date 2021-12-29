package tree

import "fmt"

type Node struct {
	Value       int
	Left, Right *Node
}

func (node Node) Print() {
	fmt.Println(node.Value)
}

func (node *Node) SetValue(value int) {
	if node == nil { //非指针类型变量不会是nil
		fmt.Println("Setting value to nil node，Ignored.")
		return
	}
	node.Value = value //都是使用 . 访问成员变量
}

//使用工厂函数，自己控制构建
func CreateNode(value int) *Node {
	return &Node{Value: value} //返回局部变量的地址，在c++中是非常典型的错误。go语言可以正常使用。
}
