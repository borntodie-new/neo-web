package neo

type node struct {
	pattern  string  // 叶子节点唯一标识 例如：/study/:lang
	part     string  // 单个节点的唯一标识 例如：study、:lang
	children []*node // 子节点 例如：[study, :lang]
	isWild   bool    // 是否精确匹配，part含有 : 或 * 为true
}

// 查询子节点是否含有part节点
// [/ , user, home]
func (n *node) search(part string) *node {
	if n.children == nil {
		n.children = make([]*node, 0)
	}
	for _, child := range n.children {
		// 精确匹配，优先级高
		if child.part == part {
			return child
		}
		// 模糊匹配，优先级低
		if child.isWild {
			return child
		}
	}
	return nil
}
