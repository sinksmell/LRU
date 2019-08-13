package main

import (
	"fmt"
	"sync"
)

// Cache interface
type Cacher interface {
	Get(key interface{}) interface{}
	Put(key, value interface{})
}

type Node struct {
	Key   interface{}
	Value interface{}
	Pre   *Node
	Next  *Node
}

// LRUCache
type LRUCache struct {
	cap     int // capacity of cache
	head    *Node
	tail    *Node
	nodeMap map[interface{}]*Node
	mutex   sync.Mutex
}

// Get value from cache by key
func (this *LRUCache) Get(key interface{}) interface{} {
	var (
		node *Node
		exit bool
	)
	this.mutex.Lock()
	defer this.mutex.Unlock()
	if node, exit = this.nodeMap[key]; !exit {
		return nil
	}
	// 调整节点到链表头部
	this.remove(node)
	this.addFirst(node)
	return node.Value
}

// Put key value into cache
func (this *LRUCache) Put(key, value interface{}) {
	var (
		node *Node
		exit bool
	)
	this.mutex.Lock()
	defer this.mutex.Unlock()
	if node, exit = this.nodeMap[key]; !exit {
		// key-value 不存在
		if len(this.nodeMap) >= this.cap {
			this.removeLast()
		}
		node = NewNode(key, value)
		this.addFirst(node)
		return
	}
	// key-value 已经存在 则调整节点至链表首部
	node.Value = value
	this.remove(node)
	this.addFirst(node)
}

// NewNode generate a Node
func NewNode(key, value interface{}) *Node {
	return &Node{Key: key, Value: value}
}

// NewLRUCache generate a LRUCache
func NewLRUCache(cap int) *LRUCache {
	cache := &LRUCache{}
	cache.cap = cap
	cache.nodeMap = make(map[interface{}]*Node)
	return cache
}

// addFirst add node to the first of list
func (this *LRUCache) addFirst(node *Node) {
	// 传入的 node 并非都是新建的，也可能是复用之前的 node
	// 所以要把前驱指针清空
	node.Pre = nil
	if this.head == nil {
		this.head = node
		this.tail = node
	} else {
		node.Next = this.head
		this.head.Pre = node
		this.head = node
	}
	this.nodeMap[node.Key] = node
}

// remove delete the node  in list
func (this *LRUCache) remove(node *Node) {
	pre := node.Pre
	next := node.Next
	if pre != nil {
		pre.Next = next
	} else {
		this.head = next
	}
	if next != nil {
		next.Pre = pre
	} else {
		this.tail = pre
	}
	delete(this.nodeMap, node.Key)
}

// removeLast   remove the last node of list
func (this *LRUCache) removeLast() {
	if this.tail == nil {
		return
	}
	// 删除key-value
	delete(this.nodeMap, this.tail.Key)
	pre := this.tail.Pre
	if pre == nil {
		this.tail = nil
		this.head = nil
	} else {
		pre.Next = nil
		this.tail = pre
	}
}

func main() {
	// 测试数据 两个值的切片代表 put kv
	// 一个值的则是根据 key 来查找
	inputs := [][]int{
		{10, 13}, {3, 17}, {6, 11}, {10, 5}, {9, 10},
		{13}, {2, 19}, {2}, {3}, {5, 25},
		{8}, {9, 22}, {5, 5}, {1, 30}, {11},
		{9, 12}, {7}, {5}, {8}, {9},
		{4, 30}, {9, 3}, {9}, {10}, {10},
		{6, 14}, {3, 1}, {3}, {10, 11}, {8},
		{2, 14}, {1}, {5}, {4}, {11, 4},
		{12, 24}, {5, 18}, {13}, {7, 23}, {8},
		{12}, {3, 27}, {2, 12}, {5}, {2, 9},
		{13, 4}, {8, 18}, {1, 7}, {11, 7}, {5, 2},
	}
	var  cache Cacher= NewLRUCache(10)
	for _, input := range inputs {
		if len(input) == 1 {
			res := cache.Get(input[0])
			fmt.Printf("Get key: %d, value is %+v\t", input[0], res)
			if res == nil {
				fmt.Println("miss!")
			} else {
				fmt.Println()
			}
		} else {
			cache.Put(input[0], input[1])
			fmt.Printf("Put kv : %d %d\n", input[0], input[1])
		}
	}

}
