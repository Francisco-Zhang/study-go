package main

import (
	"fmt"

	"me.com/tree"
)

//通过Embedding（内嵌）的方式扩展
type myTreeNode struct {
	*tree.Node //语法糖，相当于 Node *tree.Node，体现 Embedding 内嵌。
}

func (myNodde *myTreeNode) postOrder() {
	if myNodde == nil || myNodde.Node == nil {
		return
	}
	left := myTreeNode{myNodde.Left} //相当于 myNode.Node.Left
	left.postOrder()
	right := myTreeNode{myNodde.Right}
	right.postOrder()
	myNodde.Print() //也可以直接调用内嵌 Node 的方法
}

//类似于继承，但不是继承，只是一种语法糖。
func (myNodde *myTreeNode) Traverse() {
	fmt.Println("this method is shadowed.")
}

func main() {

	root := myTreeNode{&tree.Node{Value: 3}}
	root.Left = &tree.Node{} // root变为 myTreeNode 后，下面的 tree.Node 原有方法不需要修改。
	root.Right = &tree.Node{5, nil, nil}

	root.Right.Left = new(tree.Node)
	root.Left.Right = tree.CreateNode(2)
	root.Right.Left.SetValue(4)

	root.Traverse()      //调用父类方法，一种语法糖，在父类不存在Trverse方法时，自动转换为子类方法
	root.Node.Traverse() //调用子类方法

	root.postOrder()

	var bassRoot *tree.Node
	//baseRoot =  &root
	//与继承根本区别，其他语言可以这样调用。但是Go语言不行，对Go来说两个类型没有关系。
	//Go语言是通过接口来实现 对父类的引用的。
}
