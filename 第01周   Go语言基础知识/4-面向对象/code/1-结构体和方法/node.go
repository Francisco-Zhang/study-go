package main

import "fmt"

type treeNode struct {
	value       int
	left, right *treeNode
}

func (node treeNode) print() {
	fmt.Println(node.value)
}

func print2(node treeNode) {
	fmt.Println(node.value)
}

func (node treeNode) setValue(value int) {
	node.value = value
}

func (node *treeNode) setValue2(value int) {
	if node == nil { //非指针类型变量不会是nil
		fmt.Println("Setting value to nil node，Ignored.")
		return
	}
	node.value = value //都是使用 . 访问成员变量
}

//中序遍历
func (node *treeNode) traverse() {
	if node == nil {
		return
	}
	node.left.traverse() //其他语言 需要判断 node.left == nil,go不需要。
	node.print()
	node.right.traverse()
}

//使用工厂函数，自己控制构建
func createNode(value int) *treeNode {
	return &treeNode{value: value} //返回局部变量的地址，在c++中是非常典型的错误。go语言可以正常使用。
}

func main() {
	//定义方式一
	var root treeNode //会进行初始化
	//方式二
	root = treeNode{value: 3}
	//方式三
	root.left = &treeNode{}
	//方式四
	root.right = &treeNode{5, nil, nil}
	//方式五  new是内置函数，返回值就是地址
	root.right.left = new(treeNode)

	root.left = createNode(4)

	nodes := []treeNode{
		{value: 2},
		{},
		{6, nil, &root},
	}

	fmt.Println(nodes)

	//下面这两种方式是一样的。
	root.print()
	print2(root)

	//go语言所有的参数都是传值
	root.right.left.setValue(9)
	root.right.left.print()

	root.right.left.setValue2(9)
	root.right.left.print()

	pRoot := &root
	pRoot.setValue2(10)
	pRoot.print()

	var pRootNil *treeNode
	pRootNil.setValue2(30)
	pRootNil.setValue(2) //指针也可以调用值接收者方法
	pRootNil = pRoot
	pRootNil.setValue2(40)
	pRootNil.print()
}
