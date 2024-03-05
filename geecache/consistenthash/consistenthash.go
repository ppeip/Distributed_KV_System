package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

// Hash maps bytes to uint32
type Hash func(data []byte) uint32

// Map constains all hashed keys
type ConsistentHash struct {
	hash         Hash  //hash函数
	replicas     int   //虚拟节点倍数
	virtualNodes []int // Sorts
	hashMap      map[int]string
}

// 实例哈希环
func New(replicas int, fn Hash) *ConsistentHash {
	// 默认hash函数为crc32.ChecksumIEEE
	if fn == nil {
		fn = crc32.ChecksumIEEE
	}
	ch := &ConsistentHash{
		replicas: replicas,
		hash:     fn,
		hashMap:  make(map[int]string),
	}

	return ch
}

// 添加真实节点
func (ch *ConsistentHash) AddTrueNode(nodes ...string) {
	for _, node := range nodes {
		for i := 0; i < ch.replicas; i++ {
			// 虚拟节点的名称是： key + strconv.Itoa(i)，在计算哈希值，然后添加到环上
			hash := int(ch.hash([]byte(node + strconv.Itoa(i))))
			ch.virtualNodes = append(ch.virtualNodes, hash)
			ch.hashMap[hash] = node //记录虚拟节点对真实节点的映射
		}
	}
	sort.Ints(ch.virtualNodes)
}

// 选择真实节点
func (ch *ConsistentHash) GetTrueNode(key string) string {
	if len(ch.virtualNodes) == 0 {
		return ""
	}

	hash := int(ch.hash([]byte(key)))
	// Binary search for appropriate replicas.
	idx := sort.Search(len(ch.virtualNodes), func(i int) bool {
		return ch.virtualNodes[i] >= hash
	})

	return ch.hashMap[ch.virtualNodes[idx%len(ch.virtualNodes)]]
}

// 删除节点
