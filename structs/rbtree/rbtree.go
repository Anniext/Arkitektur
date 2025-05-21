package rbtree

import (
	"container/list"
	"fmt"
	"math"
	"os"
	"sync"
)

const ERBColorRed = true
const ERBColorBlack = false

type RBTreeNode[T any] struct {
	key        T
	left       *RBTreeNode[T]
	right      *RBTreeNode[T]
	parent     *RBTreeNode[T]
	color      bool
	childCount int32
}

type RBTree[T any] struct {
	root       *RBTreeNode[T]
	sentinel   *RBTreeNode[T]
	keyNodeMap map[any]*RBTreeNode[T]
	size       int32
	capacity   int32
	mutex      sync.RWMutex
	Less       func(a, b T) bool
	GetOnlyKey func(a T) any
	GetValue   func(a T) any
}

// NewRBTree 创建树
func NewRBTree[T any](capacity int32, f1 func(a, b T) bool, f2 func(a T) any, f3 func(a T) any) *RBTree[T] {
	var zero T
	nilNode := &RBTreeNode[T]{
		key:        zero,
		left:       nil,
		right:      nil,
		parent:     nil,
		color:      ERBColorBlack,
		childCount: 0,
	}
	nilNode.left = nilNode
	nilNode.right = nilNode
	nilNode.parent = nilNode

	//GetOnlyKey = f1
	//GetValue = f2
	//Less = f3

	if capacity == 0 {
		capacity = int32(math.MaxInt32)
	}

	return &RBTree[T]{
		root:       nilNode,
		sentinel:   nilNode,
		keyNodeMap: make(map[any]*RBTreeNode[T]),
		size:       0,
		capacity:   capacity,
		Less:       f1,
		GetOnlyKey: f2,
		GetValue:   f3,
	}
}

func (tree *RBTree[T]) UpdateLess(less func(a T, b T) bool) {

	tree.Less = less
}

// Size 当前大小
func (tree *RBTree[T]) Size() int32 {
	tree.mutex.RLock()
	defer tree.mutex.RUnlock()
	return tree.size
}

// Empty 是否为空
func (tree *RBTree[T]) Empty() bool {
	tree.mutex.RLock()
	defer tree.mutex.RUnlock()
	return tree.size == 0
}

// Keys 升序uid数组
func (tree *RBTree[T]) Keys() []interface{} {
	tree.mutex.RLock()
	defer tree.mutex.RUnlock()
	keys := make([]interface{}, tree.size)
	for it, i := tree.min(tree.root), 0; it != nil; it, i = tree.next(it), i+1 {
		keys[i] = tree.GetOnlyKey(it.key)
	}
	return keys
}

// ReverseKeys 降序uid数组
func (tree *RBTree[T]) ReverseKeys() []interface{} {
	tree.mutex.RLock()
	defer tree.mutex.RUnlock()
	keys := make([]interface{}, tree.size)
	for it, i := tree.max(tree.root), 0; it != nil; it, i = tree.prev(it), i+1 {
		keys[i] = tree.GetOnlyKey(it.key)
	}
	return keys
}

// Get 通过Key查询节点
func (tree *RBTree[T]) Get(key T) *RBTreeNode[T] {
	tree.mutex.RLock()
	defer tree.mutex.RUnlock()
	if node, ok := tree.keyNodeMap[tree.GetOnlyKey(key)]; ok {
		return node
	}
	return nil
}

// Put 放如一个key
func (tree *RBTree[T]) Put(key T) {
	tree.mutex.Lock()
	defer tree.mutex.Unlock()
	if node, ok := tree.keyNodeMap[tree.GetOnlyKey(key)]; ok {
		tree.deleteNode(node)
	}
	node := &RBTreeNode[T]{
		key:        key,
		left:       nil,
		right:      nil,
		parent:     nil,
		color:      ERBColorBlack,
		childCount: 0,
	}
	tree.insertNode(node)

	if tree.size > tree.capacity {
		minNode := tree.min(tree.root)
		tree.deleteNode(minNode)
	}
}

// FindNodeByKey 通过key查询节点和排名
func (tree *RBTree[T]) FindNodeByKey(key T) (*RBTreeNode[T], int32) {
	tree.mutex.RLock()
	defer tree.mutex.RUnlock()
	var temp, sentinel *RBTreeNode[T]
	var rank int32
	//var ret int8

	sentinel = tree.sentinel
	temp = tree.root
	for temp != sentinel {
		if tree.GetOnlyKey(key) == tree.GetOnlyKey(temp.key) {
			return temp, rank + temp.right.childCount + 1
		} else if tree.Less(key, temp.key) {
			rank += temp.right.childCount + 1
			temp = temp.left
		} else {
			temp = temp.right
		}
	}
	return nil, 0
}

// FindRankByKey 通过key查询排名
func (tree *RBTree[T]) FindRankByKey(key T) int32 {
	tree.mutex.RLock()
	defer tree.mutex.RUnlock()
	var temp, sentinel *RBTreeNode[T]
	var rank int32
	//var ret int8

	sentinel = tree.sentinel
	temp = tree.root
	for temp != sentinel {
		if tree.GetOnlyKey(key) == tree.GetOnlyKey(temp.key) {
			return rank + temp.right.childCount + 1
		} else if tree.Less(key, temp.key) {
			rank += temp.right.childCount + 1
			temp = temp.left
		} else {
			temp = temp.right
		}
	}
	return 0
}

func (tree *RBTree[T]) findNodeByRank(rank int32) *RBTreeNode[T] {
	var temp, sentinel *RBTreeNode[T]
	var baseRank, tempRank int32

	sentinel = tree.sentinel
	temp = tree.root
	for temp != sentinel {
		tempRank = baseRank + temp.right.childCount + 1
		if rank == tempRank {
			return temp
		} else if rank < tempRank {
			temp = temp.right
		} else {
			baseRank += temp.right.childCount + 1
			temp = temp.left
		}
	}
	return nil
}

// FindKeyByRank 通过排名查询节点
func (tree *RBTree[T]) FindKeyByRank(rank int32) (key T) {
	tree.mutex.RLock()
	defer tree.mutex.RUnlock()
	node := tree.findNodeByRank(rank)
	if node != nil {
		key = node.key
	}
	return
}

// FindKeysByRange 通过排名和查询数量, 查询节点数组
func (tree *RBTree[T]) FindKeysByRange(start, num int32) (keyList []T) {
	tree.mutex.RLock()
	defer tree.mutex.RUnlock()

	node := tree.findNodeByRank(start)
	if node == nil || node == tree.sentinel {
		return
	}
	keyList = append(keyList, node.key)
	for i := int32(1); i < num; i++ {
		node = tree.prev(node)
		if node == nil || node == tree.sentinel {
			break
		}
		keyList = append(keyList, node.key)
	}
	return
}

func (tree *RBTree[T]) floor(key T) *RBTreeNode[T] {
	found := false
	var node, temp, sentinel *RBTreeNode[T]
	sentinel = tree.sentinel
	temp = tree.root
	for temp != sentinel {
		if tree.GetOnlyKey(key) == tree.GetOnlyKey(temp.key) {
			return temp
		} else if tree.Less(key, temp.key) {
			temp = temp.left
		} else {
			node, found = temp, true
			temp = temp.right
		}
	}
	if found {
		return node
	}
	return nil
}

func (tree *RBTree[T]) ceiling(key T) *RBTreeNode[T] {
	found := false
	var node, temp, sentinel *RBTreeNode[T]
	sentinel = tree.sentinel
	temp = tree.root
	for temp != sentinel {
		if tree.GetOnlyKey(key) == tree.GetOnlyKey(temp.key) {
			return temp
		} else if tree.Less(key, temp.key) {
			node, found = temp, true
			temp = temp.left
		} else {
			temp = temp.right
		}
	}
	if found {
		return node
	}
	return nil
}

// FindKeysBigger 查询大于等于key的节点数组 默认降序
func (tree *RBTree[T]) FindKeysBigger(key T, num int32, reverse bool) (keyList []T) {
	tree.mutex.RLock()
	defer tree.mutex.RUnlock()
	if reverse {
		temp := tree.ceiling(key)
		if temp == nil || temp == tree.sentinel || tree.Less(temp.key, key) {
			return
		}
		for it, i := temp, int32(0); it != nil && i < num; it, i = tree.next(it), i+1 {
			keyList = append(keyList, it.key)
		}
	} else {
		for it, i := tree.max(tree.root), int32(0); it != nil && i < num && tree.Less(key, it.key); it, i = tree.prev(it), i+1 {
			keyList = append(keyList, it.key)
		}
	}
	return
}

// FindKeysSmaller 查询小于key的节点数组 默认降序
func (tree *RBTree[T]) FindKeysSmaller(key T, num int32, reverse bool) (keyList []T) {
	tree.mutex.RLock()
	defer tree.mutex.RUnlock()
	if reverse {
		for it, i := tree.min(tree.root), int32(0); it != nil && i < num && tree.Less(it.key, key); it, i = tree.next(it), i+1 {
			keyList = append(keyList, it.key)
		}
	} else {
		temp := tree.floor(key)
		if temp == nil || temp == tree.sentinel || tree.Less(key, temp.key) {
			return
		}
		for it, i := temp, int32(0); it != nil && i < num; it, i = tree.prev(it), i+1 {
			keyList = append(keyList, it.key)
		}
	}
	return
}

// FindKeysBetween 查询大于等于keyLeft, 小于keyRight的节点数组 默认降序
func (tree *RBTree[T]) FindKeysBetween(keyLeft, keyRight T, num int32, reverse bool) (keyList []T) {
	tree.mutex.RLock()
	defer tree.mutex.RUnlock()
	nodeLeft := tree.ceiling(keyLeft)
	if nodeLeft == nil || tree.Less(nodeLeft.key, keyLeft) {
		return
	}
	nodeRight := tree.floor(keyRight)
	if nodeRight == nil || tree.Less(keyRight, nodeRight.key) {
		return
	}
	if reverse {
		for it, i := nodeLeft, int32(0); it != nil && i < num && tree.Less(it.key, keyRight); it, i = tree.next(it), i+1 {
			keyList = append(keyList, it.key)
		}
	} else {
		for it, i := nodeRight, int32(0); it != nil && i < num && tree.Less(keyLeft, it.key); it, i = tree.prev(it), i+1 {
			keyList = append(keyList, it.key)
		}
	}
	return
}

func (tree *RBTree[T]) insertNode(node *RBTreeNode[T]) {
	originalNode := node
	var temp, root, sentinel *RBTreeNode[T]
	root, sentinel = tree.root, tree.sentinel

	tree.size += 1
	tree.keyNodeMap[tree.GetOnlyKey(node.key)] = node

	if root == sentinel {
		node.parent = nil
		node.left = sentinel
		node.right = sentinel
		node.color = ERBColorBlack
		node.childCount = 1
		tree.root = node
		return
	}
	var p **RBTreeNode[T]

	// insert value
	temp = root
	for {
		if tree.Less(node.key, temp.key) {
			p = &temp.left
		} else {
			p = &temp.right
		}
		if *p == sentinel {
			break
		}
		temp = *p
	}

	*p = node
	node.parent = temp
	node.left = sentinel
	node.right = sentinel
	node.color = ERBColorRed

	temp = nil
	// re-balance tree
	for node != root && node.parent.color == ERBColorRed {
		if node.parent == node.parent.parent.left {
			temp = node.parent.parent.right
			if temp.color == ERBColorRed {
				node.parent.color = ERBColorBlack
				temp.color = ERBColorBlack
				node.parent.parent.color = ERBColorRed
				node = node.parent.parent
			} else {
				if node == node.parent.right {
					node = node.parent
					tree.leftRotate(node)
				}
				node.parent.color = ERBColorBlack
				node.parent.parent.color = ERBColorRed
				tree.rightRotate(node.parent.parent)
			}
		} else {
			temp = node.parent.parent.left
			if temp.color == ERBColorRed {
				node.parent.color = ERBColorBlack
				temp.color = ERBColorBlack
				node.parent.parent.color = ERBColorRed
				node = node.parent.parent
			} else {
				if node == node.parent.left {
					node = node.parent
					tree.rightRotate(node)
				}
				node.parent.color = ERBColorBlack
				node.parent.parent.color = ERBColorRed
				tree.leftRotate(node.parent.parent)
			}
		}
	}
	tree.root.color = ERBColorBlack
	tempNode := *p
	tempNode = originalNode
	for tempNode != nil && tempNode != tree.sentinel {
		tempNode.childCount++
		tempNode = tempNode.parent
	}
}

// Remove 通过key删除节点
func (tree *RBTree[T]) Remove(key T) {
	tree.mutex.Lock()
	defer tree.mutex.Unlock()

	if node, ok := tree.keyNodeMap[tree.GetOnlyKey(key)]; ok {
		tree.deleteNode(node)
	}
}

func (tree *RBTree[T]) deleteNode(node *RBTreeNode[T]) {
	var red bool
	var sentinel, subst, temp, w *RBTreeNode[T]
	var root **RBTreeNode[T]
	root = &tree.root
	sentinel = tree.sentinel

	tree.size -= 1
	delete(tree.keyNodeMap, tree.GetOnlyKey(node.key))
	tempNode := node
	for tempNode != nil && tempNode != tree.sentinel {
		tempNode.childCount--
		tempNode = tempNode.parent
	}

	// delete
	if node.left == sentinel {
		temp = node.right
		subst = node
	} else if node.right == sentinel {
		temp = node.left
		subst = node
	} else {
		subst = tree.min(node.right)
		temp = subst.right
	}

	if subst == *root {
		*root = temp
		temp.color = ERBColorBlack
		node.left = nil
		node.right = nil
		node.parent = nil
		var zero T
		node.key = zero
		return
	}

	red = subst.color

	if subst == subst.parent.left {
		subst.parent.left = temp
	} else {
		subst.parent.right = temp
	}
	tempNode = subst
	for tempNode != node && tempNode != nil && tempNode != tree.sentinel {
		tempNode.childCount -= 1
		tempNode = tempNode.parent
	}

	if subst == node {
		temp.parent = subst.parent
	} else {
		if subst.parent == node {
			temp.parent = subst
		} else {
			temp.parent = subst.parent
		}
		subst.left = node.left
		subst.childCount, node.childCount = node.childCount, subst.childCount
		subst.right = node.right
		subst.parent = node.parent
		subst.color = node.color
		if node == *root {
			*root = subst
		} else {
			if node == node.parent.left {
				node.parent.left = subst
			} else {
				node.parent.right = subst
			}
		}
		if subst.left != sentinel {
			subst.left.parent = subst
		}
		if subst.right != sentinel {
			subst.right.parent = subst
		}
	}

	var zero T
	node.left = nil
	node.right = nil
	node.parent = nil
	node.key = zero
	if red {
		return
	}

	// delete fixup
	for temp != *root && temp.color == ERBColorBlack {
		if temp == temp.parent.left {
			w = temp.parent.right
			if w.color == ERBColorRed {
				w.color = ERBColorBlack
				temp.parent.color = ERBColorRed
				tree.leftRotate(temp.parent)
				w = temp.parent.right
			}
			if w.left.color == ERBColorBlack && w.right.color == ERBColorBlack {
				w.color = ERBColorRed
				temp = temp.parent
			} else {
				if w.right.color == ERBColorBlack {
					w.left.color = ERBColorBlack
					w.color = ERBColorRed
					tree.rightRotate(w)
					w = temp.parent.right
				}

				w.color = temp.parent.color
				temp.parent.color = ERBColorBlack
				w.right.color = ERBColorBlack
				tree.leftRotate(temp.parent)
				temp = *root
			}
		} else {
			w = temp.parent.left
			if w.color == ERBColorRed {
				w.color = ERBColorBlack
				temp.parent.color = ERBColorRed
				tree.rightRotate(temp.parent)
				w = temp.parent.left
			}
			if w.left.color == ERBColorBlack && w.right.color == ERBColorBlack {
				w.color = ERBColorRed
				temp = temp.parent
			} else {
				if w.left.color == ERBColorBlack {
					w.right.color = ERBColorBlack
					w.color = ERBColorRed
					tree.leftRotate(w)
					w = temp.parent.left
				}
				w.color = temp.parent.color
				temp.parent.color = ERBColorBlack
				w.left.color = ERBColorBlack
				tree.rightRotate(temp.parent)
				temp = *root
			}
		}
	}
	temp.color = ERBColorBlack
}

func (tree *RBTree[T]) leftRotate(node *RBTreeNode[T]) {
	var temp *RBTreeNode[T]
	temp = node.right
	parentCount := node.childCount
	node.right = temp.left
	node.childCount -= temp.childCount
	temp.childCount = parentCount
	if temp.left != tree.sentinel {
		temp.left.parent = node
		node.childCount += temp.left.childCount
	}
	temp.parent = node.parent
	if node == tree.root {
		tree.root = temp
	} else if node == node.parent.left {
		node.parent.left = temp
	} else {
		node.parent.right = temp
	}
	temp.left = node
	node.parent = temp
}

func (tree *RBTree[T]) rightRotate(node *RBTreeNode[T]) {
	var temp *RBTreeNode[T]
	temp = node.left
	parentCount := node.childCount
	node.left = temp.right
	node.childCount -= temp.childCount
	temp.childCount = parentCount
	if temp.right != tree.sentinel {
		temp.right.parent = node
		node.childCount += temp.right.childCount
	}
	temp.parent = node.parent
	if node == tree.root {
		tree.root = temp
	} else if node == node.parent.right {
		node.parent.right = temp
	} else {
		node.parent.left = temp
	}
	temp.right = node
	node.parent = temp
}

func (tree *RBTree[T]) next(node *RBTreeNode[T]) *RBTreeNode[T] {
	var root, sentinel, parent *RBTreeNode[T]
	sentinel = tree.sentinel
	if node.right != sentinel {
		return tree.min(node.right)
	}
	root = tree.root
	for {
		parent = node.parent
		if node == root {
			return nil
		}
		if node == parent.left {
			return parent
		}
		node = parent
	}
}

func (tree *RBTree[T]) prev(node *RBTreeNode[T]) *RBTreeNode[T] {
	var root, sentinel, parent *RBTreeNode[T]
	sentinel = tree.sentinel
	if node.left != sentinel {
		return tree.max(node.left)
	}
	root = tree.root
	for {
		parent = node.parent
		if node == root {
			return nil
		}
		if node == parent.right {
			return parent
		}
		node = parent
	}
}

func (tree *RBTree[T]) min(node *RBTreeNode[T]) *RBTreeNode[T] {
	sentinel := tree.sentinel
	if node == nil || node == sentinel {
		return nil
	}
	for node.left != sentinel {
		node = node.left
	}
	return node
}

func (tree *RBTree[T]) max(node *RBTreeNode[T]) *RBTreeNode[T] {
	sentinel := tree.sentinel
	if node == nil || node == sentinel {
		return nil
	}
	for node.right != sentinel {
		node = node.right
	}
	return node
}

func (tree *RBTree[T]) calcChildCount(node *RBTreeNode[T]) int32 {
	if node == nil || node == tree.sentinel {
		return 0
	}
	return 1 + tree.calcChildCount(node.left) + tree.calcChildCount(node.right)
}

// Left 树的最左边
func (tree *RBTree[T]) Left() *RBTreeNode[T] {
	tree.mutex.RLock()
	defer tree.mutex.RUnlock()
	return tree.min(tree.root)
}

// Right 树的最右边
func (tree *RBTree[T]) Right() *RBTreeNode[T] {
	tree.mutex.RLock()
	defer tree.mutex.RUnlock()
	return tree.max(tree.root)
}

// CalcChildCount 树的孩子数量
func (tree *RBTree[T]) CalcChildCount(node *RBTreeNode[T]) int32 {
	tree.mutex.RLock()
	defer tree.mutex.RUnlock()
	return tree.calcChildCount(node)
}

// GenDot 生成树的dot文件
func (tree *RBTree[T]) GenDot(fileName string) {
	tree.mutex.RLock()
	defer tree.mutex.RUnlock()

	getNodeName := func(node *RBTreeNode[T]) string {
		return fmt.Sprintf("n%d_%d_%d", tree.GetOnlyKey(node.key), tree.GetValue(node.key), node.childCount)
	}
	getEdgeName := func(node1, node2 string, left bool) string {
		if left {
			return fmt.Sprintf("%s:f0->%s:f1", node1, node2)
		}
		return fmt.Sprintf("%s:f2->%s:f1", node1, node2)
	}
	const EInvisEdge, EInvis, ERed, EBlack, ENil = -1, 0, 1, 2, 3
	/*
		digraph binaryTree{
		    node[shape=circle,color=red,fontcolor=blue,fontsize=10];
		    n4_3->n2_2;
		    n4_3->n5_1;
		    n2_2->n1_0;
		    n2_2->n3_0;
		    n5_1->n2_2_left[style=invis];
		    n5_1->n6_0;

		    n1_0[color=black];
		    n2_2[color=red];
		    n3_0[color=red];
		    n4_3[color=black];
		    n5_1[color=black];
		    n6_0[color=black];
		    n2_2_left[style=invis];
		}
	*/
	// nodeEdgeMap = map[nodeEdgeString]color  0:invis  1:red  2:black  3:nil
	nodeEdgeMap := map[string]int8{}
	var nodeEdgeList []string
	var node, nodeLeft, nodeRight, edgeLeft, edgeRight string

	queue := list.New()
	queue.PushBack(tree.root)
	for {
		if queue.Len() <= 0 {
			break
		}
		it := queue.Remove(queue.Front()).(*RBTreeNode[T])
		if it.left != tree.sentinel {
			queue.PushBack(it.left)
		}
		if it.right != tree.sentinel {
			queue.PushBack(it.right)
		}
		node = getNodeName(it)
		if it.color {
			nodeEdgeMap[node] = ERed
		} else {
			nodeEdgeMap[node] = EBlack
		}
		nodeEdgeList = append(nodeEdgeList, node)
		if it.left != tree.sentinel {
			edgeLeft = getEdgeName(node, getNodeName(it.left), true)
			nodeEdgeMap[edgeLeft] = ENil
		} else {
			nodeLeft = node + "_left"
			nodeEdgeMap[nodeLeft] = EInvis
			nodeEdgeList = append(nodeEdgeList, nodeLeft)

			edgeLeft = getEdgeName(node, nodeLeft, true)
			nodeEdgeMap[edgeLeft] = EInvisEdge
		}
		nodeEdgeList = append(nodeEdgeList, edgeLeft)
		if it.right != tree.sentinel {
			edgeRight = getEdgeName(node, getNodeName(it.right), false)
			nodeEdgeMap[edgeRight] = ENil
		} else {
			nodeRight = node + "_right"
			nodeEdgeMap[nodeRight] = EInvis
			nodeEdgeList = append(nodeEdgeList, nodeRight)

			edgeRight = getEdgeName(node, nodeRight, false)
			nodeEdgeMap[edgeRight] = EInvisEdge
		}
		nodeEdgeList = append(nodeEdgeList, edgeRight)

	}

	f, err := os.Create(fileName)
	if err != nil {
		return
	}
	defer f.Close()
	f.WriteString("digraph RBTree{\n")
	f.WriteString("    bgcolor=\"beige\"\n")
	f.WriteString("    node[shape=record, height=.1, color=red, fontcolor=blue, fontsize=10];\n")
	var color int8
	var ok bool
	for _, nodeStr := range nodeEdgeList {
		color, ok = nodeEdgeMap[nodeStr]
		if !ok {
			continue
		}
		delete(nodeEdgeMap, nodeStr)
		f.WriteString("    ")
		f.WriteString(nodeStr)
		switch color {
		case EInvisEdge:
			f.WriteString(fmt.Sprintf("[style=invis];\n"))
		case EInvis:
			f.WriteString(fmt.Sprintf("[style=invis, label=\"<f0> | <f1> %s | <f2>\"];\n", nodeStr))
		case ERed:
			f.WriteString(fmt.Sprintf("[color=red, label=\"<f0> | <f1> %s | <f2>\"];\n", nodeStr))
		case EBlack:
			f.WriteString(fmt.Sprintf("[color=black, label=\"<f0> | <f1> %s | <f2>\"];\n", nodeStr))
		case ENil:
			f.WriteString(";\n")
		}
	}
	f.WriteString("}\n")
}
