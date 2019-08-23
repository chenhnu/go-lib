package ptcollection

type bst struct {
	root *node
}

func NewBST(args ...interface{}) *bst {
	if args == nil {
		return &bst{}
	}
	res := bst{}
	res.root = &node{
		val: args[0],
	}
	args = args[1:]
	for _, arg := range args {
		res.addNode(arg)
	}
	return &res
}

//添加一个节点到BST中
func (b *bst)addNode(arg interface{})  {
	p:=b.root
	for{
		tmp:=node{val:arg}
		if compare(p.val,arg){
			if p.lChild==nil{
				p.lChild=&tmp
				break
			}else{
				p=p.lChild
			}
		}else{
			if p.rChild==nil{
				p.rChild=&tmp
				break
			}else{
				p=p.rChild
			}
		}
	}
}
