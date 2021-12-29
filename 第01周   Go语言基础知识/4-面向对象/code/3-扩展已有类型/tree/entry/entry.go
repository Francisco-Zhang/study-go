package main

import "me.com/tree"

//通过组合的方式
type myTreeNode struct {
	node *tree.Node
}

func (myNodde *myTreeNode) postOrder() {
	if myNodde == nil || myNodde.node == nil {
		return
	}
	left := myTreeNode{myNodde.node.Left}
	left.postOrder()
	right := myTreeNode{myNodde.node.Right}
	right.postOrder()
	myNodde.node.Print()
}

func main() {
	var root tree.Node
	root = tree.Node{Value: 3}
	root.Left = &tree.Node{}
	root.Right = &tree.Node{5, nil, nil}

	root.Right.Left = new(tree.Node)

	root.Left.Right = tree.CreateNode(2)
	root.Right.Left.SetValue(4)

	myRoot := myTreeNode{&root}
	myRoot.postOrder()

}
