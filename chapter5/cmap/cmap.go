package cmap

import (
	"math"
	"sync/atomic"
)

// ConcurrentMap 代表并发安全的字典的接口。
type ConcurrentMap interface {
	// Concurrency 会返回并发量。
	Concurrency() int
	// Put 会推送一个键-元素对。
	// 注意！参数element的值不能为nil。
	// 第一个返回值表示是否新增了键-元素对。
	// 若键已存在，新元素值会替换旧的元素值。
	Put(key string, element interface{}) (bool, error)
	// Get 会获取与指定键关联的那个元素。
	// 若返回nil，则说明指定的键不存在。
	Get(key string) interface{}
	// Delete 会删除指定的键-元素对。
	// 若结果值为true则说明键已存在且已删除，否则说明键不存在。
	Delete(key string) bool
	// Len 会返回当前字典中键-元素对的数量。
	Len() uint64
}

// myConcurrentMap 代表ConcurrentMap接口的实现类型。
type myConcurrentMap struct {
	concurrency int  //并发量
	segments    []Segment  //散列段
	total       uint64
}

// NewConcurrentMap 会创建一个ConcurrentMap类型的实例。
// 参数pairRedistributor可以为nil。
func NewConcurrentMap(
	concurrency int,
	pairRedistributor PairRedistributor) (ConcurrentMap, error) {
	if concurrency <= 0 {
		return nil, newIllegalParameterError("concurrency is too small")
	}
	if concurrency > MAX_CONCURRENCY {
		return nil, newIllegalParameterError("concurrency is too large")
	}
	cmap := &myConcurrentMap{}
	cmap.concurrency = concurrency
	cmap.segments = make([]Segment, concurrency) //散列段
	for i := 0; i < concurrency; i++ {
		cmap.segments[i] =
			newSegment(DEFAULT_BUCKET_NUMBER, pairRedistributor) //一个散列段默认有16个散列桶
	}
	return cmap, nil
}
//返回当前并发量
func (cmap *myConcurrentMap) Concurrency() int {
	return cmap.concurrency
}

func (cmap *myConcurrentMap) Put(key string, element interface{}) (bool, error) {
	p, err := newPair(key, element)
	if err != nil {
		return false, err
	}
	s := cmap.findSegment(p.Hash()) //根据key的hash值找到散列段
	ok, err := s.Put(p)
	if ok {
		//原子增加1
		atomic.AddUint64(&cmap.total, 1)
	}
	return ok, err
}

func (cmap *myConcurrentMap) Get(key string) interface{} {
	keyHash := hash(key)
	s := cmap.findSegment(keyHash) //根据key的hash值找到散列段
	pair := s.GetWithHash(key, keyHash) //从散列段里面找到Pair元素
	if pair == nil {
		return nil
	}
	return pair.Element()
}

func (cmap *myConcurrentMap) Delete(key string) bool {
	s := cmap.findSegment(hash(key))  //根据key的hash值找到散列段
	if s.Delete(key) {
		//原子减1
		atomic.AddUint64(&cmap.total, ^uint64(0))
		return true
	}
	return false
}

func (cmap *myConcurrentMap) Len() uint64 {
	//原子Load 键值对总数
	return atomic.LoadUint64(&cmap.total)
}

// findSegment 会根据给定参数寻找并返回对应散列段。
func (cmap *myConcurrentMap) findSegment(keyHash uint64) Segment {
	//核心思想是：使用高位的几个字节来决定散列段的索引，这样可以让K-V元素在segments中分布更广一些，更均匀一些
	if cmap.concurrency == 1 {
		//如果并发量为1，则只有一个散列段
		return cmap.segments[0]
	}
	var keyHash32 uint32
	if keyHash > math.MaxUint32 {
		keyHash32 = uint32(keyHash >> 32)
	} else {
		keyHash32 = uint32(keyHash)
	}
	return cmap.segments[int(keyHash32>>16)%(cmap.concurrency-1)]
}
