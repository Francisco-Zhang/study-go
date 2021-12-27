package tree

//定义在不同的文件中，使用起来还是一样的
func (node *Node) Traverse() {
	if node == nil {
		return
	}
	node.Left.Traverse() //其他语言 需要判断 node.left == nil,go不需要。
	node.Print()
	node.Right.Traverse()
}
