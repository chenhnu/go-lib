package ptcollection

type treeNode struct {
	lChild *treeNode
	rChild *treeNode
	val    interface{}
}

func compare(a, b interface{}) bool {
	switch a.(type) {
	case int:
		return a.(int) > b.(int)
	case string:
		return a.(string) > b.(string)
	default:
		return false
	}
}
