package consistent

import (
	"errors"
	"hash/crc32"
	"sort"
	"strconv"
)

// SortedKeys 排序后的key
type SortedKeys []uint32

// 下面这几个函数都是vscode自动生成的，Less好像必须写
func (x SortedKeys) Len() int { return len(x) }

func (x SortedKeys) Less(i, j int) bool { return x[i] < x[j] }

func (x SortedKeys) Swap(i, j int) { x[i], x[j] = x[j], x[i] }

// ConsistentHashing  一致性哈希数据结构
type ConsistentHashing struct {
	NumOfVirtualNode int               // 虚拟节点个数，这个有啥用？
	hashSortedKeys   SortedKeys        // 排序后的所有服务器Hash值
	circleRing       map[uint32]string // map key:value  =>  服务器ID:服务器名称
	serverSet        map[string]bool   // map key:value  =>  服务器名称:bool 这里直接改成[]string不可以吗？
}

// NewConsistentHashing 新建一个一致性Hash数据结构
func NewConsistentHashing() *ConsistentHashing {
	return &ConsistentHashing{
		NumOfVirtualNode: 20,
		circleRing:       make(map[uint32]string),
		serverSet:        make(map[string]bool),
	}
}

// Get 入参：客户端名称，出参：服务器名
func (c *ConsistentHashing) Get(obj string) (string, error) {
	if len(c.serverSet) == 0 {
		return "", errors.New("Empty struct")
	}
	nearServer, _ := c.circleRing[c.hashSortedKeys[c.searchNearRingIndex(obj)]]
	return nearServer, nil
}

// Add add a node to this consistent hash ring，往哈希环中加新的服务器
func (c *ConsistentHashing) Add(node string) {
	if _, find := c.serverSet[node]; find {
		return
	}
	c.serverSet[node] = true
	// 计算服务器id
	key := c.hashKey(node)
	c.circleRing[key] = node
	// Add virtual node for "balance"
	for i := 0; i < c.NumOfVirtualNode; i++ {
		vk := c.getVirtualNodeKey(i, node)
		c.circleRing[vk] = node
	}
	c.updateSortHashKeys()
}

// Remove remove a node form this consistent hashing ring 后端移除一台服务器
func (c *ConsistentHashing) Remove(node string) {
	if _, find := c.serverSet[node]; !find {
		return // not in our serverSet
	}
	delete(c.serverSet, node)
	key := c.hashKey(node)
	delete(c.circleRing, key)

	// delete virtual node
	for i := 0; i < c.NumOfVirtualNode; i++ {
		vk := c.getVirtualNodeKey(i, node)
		delete(c.circleRing, vk)
	}
	c.updateSortHashKeys()
}

// ListNodes list whole nodes in consistent hashing ring 列出环中所有的服务器
func (c *ConsistentHashing) ListNodes() []string {
	var retList []string
	for server := range c.serverSet {
		retList = append(retList, server)
	}
	return retList
}

func (c *ConsistentHashing) getVirtualNodeKey(index int, obj string) uint32 {
	newObjStr := strconv.Itoa(index) + "-" + obj
	return c.hashKey(newObjStr)
}

// searchNearRingIndex 找客户端对应的服务器，入参：ojb 客户端名称，出参：最近的服务器ID
func (c *ConsistentHashing) searchNearRingIndex(client string) int {
	targetKey := c.hashKey(client)
	// 这个sort.Search算法是核心，返回一个int类型
	targetIndex := sort.Search(len(c.hashSortedKeys), func(i int) bool {
		// **客户端与服务端处于同一哈希空间中**
		// 此时所有的服务端id都已经排序了，找客户端c1对应的服务器的过程就是在服务器ID空间，顺时针找到第一个比客户端c1ID大的，这个服务器就是c1对应的服务器
		return c.hashSortedKeys[i] >= targetKey
	})

	if targetIndex >= len(c.hashSortedKeys) {
		targetIndex = 0
	}
	return targetIndex
}

// updateSortHashKeys 更新所有服务器ID
func (c *ConsistentHashing) updateSortHashKeys() {
	c.hashSortedKeys = nil
	for server := range c.serverSet {
		key := c.hashKey(server)
		c.hashSortedKeys = append(c.hashSortedKeys, key)
	}
	sort.Sort(c.hashSortedKeys)
}

// hashKey 求客户端（string）类型的crc32校验和，CRC32、MD5、SHA1 是一类算法，用来校验文件完整性
// 这里用相对简单的 crc32
func (c *ConsistentHashing) hashKey(node string) uint32 {
	var scratch [64]byte
	if len(node) < 64 {
		copy(scratch[:], node)
	}
	return crc32.ChecksumIEEE(scratch[:len(node)])
}
