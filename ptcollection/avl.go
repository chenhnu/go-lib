package ptcollection

type avl struct {
	root *treeNode
}

func NewAVL(args ...interface{}) *avl {
	if args == nil {
		return &avl{}
	}
	res := avl{
		root: &treeNode{val: args[0]},
	}
	args = args[1:]
	for _, arg := range args {
		res.addNode(arg)
	}
	return &res
}
func (a *avl) addNode(arg interface{}) {
	temp := treeNode{val: arg}
	p := a.root
	for {
		if compare(p.val, arg) {
			if p.lChild != nil {
				p = p.lChild
			} else {
				p.lChild = &temp
				break
			}
		} else {
			if p.rChild != nil {
				if p.lChild == nil {
					if compare(p.rChild.val, arg) {
						temp.val = p.val
						p.val = arg
					} else {
						temp.val = p.val
						p.val = p.rChild.val
						p.rChild.val = arg
					}
					p.lChild = &temp
					break
				}
				p = p.rChild
			} else {
				p.rChild = &temp
				break
			}
		}
	}
}
