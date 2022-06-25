package gee

import "strings"

type node struct {
	// pattern 路径, /p/:lang
	pattern string
	// part 路径的一部分 :lang
	part string
	// children 子节点 [doc, tutorial, intro]
	children []*node
	// 是否精确匹配 contain [':', '*'] will be true
	isWild bool
}

// matchChild 第一个匹配成功的结点，用于插入
func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		// TODO: 如果此处都是[':', '*']，就会都匹配第一个？
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

// matchChildren 匹配当前结点的孩子，用于查找
func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

// height: this node in pattern index
func (n *node) insert(pattern string, parts []string, height int) {
	if len(parts) == height {
		n.pattern = pattern
		return
	}
	part := parts[height]
	child := n.matchChild(part)
	// 没找到就新增一个结点存储
	if child == nil {
		child = &node{
			pattern:  "",
			part:     part,
			children: nil,
			isWild:   part[0] == ':' || part[0] == '*',
		}
		n.children = append(n.children, child)
	}
	// 递归插入
	child.insert(pattern, parts, height+1)
}

func (n *node) search(parts []string, height int) *node {
	// 匹配查找
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		// 防止出现空串
		if n.pattern == "" {
			return nil
		}
		return n
	}
	part := parts[height]
	// get this node children
	children := n.matchChildren(part)
	for _, child := range children {
		// 递归
		res := child.search(parts, height+1)
		if res != nil {
			return res
		}
	}
	return nil
}
