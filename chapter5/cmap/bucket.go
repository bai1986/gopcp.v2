package cmap

import (
	"bytes"
	"sync"
	"sync/atomic"
)

// Bucket 代表并发安全的散列桶的接口。
type Bucket interface {
	// Put 会放入一个键-元素对。
	// 第一个返回值表示是否新增了键-元素对。
	// 若在调用此方法前已经锁定lock，则不要把lock传入！否则必须传入对应的lock！
	Put(p Pair, lock sync.Locker) (bool, error)
	// Get 会获取指定键的键-元素对。
	Get(key string) Pair
	// GetFirstPair 会返回第一个键-元素对。
	GetFirstPair() Pair
	// Delete 会删除指定的键-元素对。
	// 若在调用此方法前已经锁定lock，则不要把lock传入！否则必须传入对应的lock！
	Delete(key string, lock sync.Locker) bool
	// Clear 会清空当前散列桶。
	// 若在调用此方法前已经锁定lock，则不要把lock传入！否则必须传入对应的lock！
	Clear(lock sync.Locker)
	// Size 会返回当前 散列桶的 尺寸。
	Size() uint64
	// String 会返回当前 散列桶的 字符串表示形式。
	String() string
}

// bucket 代表并发安全的散列桶的类型。
type bucket struct {
	// firstValue 存储的是键-元素对列表的表头。
	firstValue atomic.Value
	size       uint64
}

// 占位符。
// 由于原子值不能存储nil，所以当散列桶空时用此符占位。
var placeholder Pair = &pair{}

// newBucket 会创建一个 Bucket 类型的实例。
func newBucket() Bucket {
	b := &bucket{}
	//存储占位值，也确立原子值存储类型
	b.firstValue.Store(placeholder)
	return b
}

func (b *bucket) Put(p Pair, lock sync.Locker) (bool, error) {
	if p == nil {
		return false, newIllegalParameterError("pair is nil")
	}
	if lock != nil {
		lock.Lock()
		defer lock.Unlock()
	}
	//获取桶中第一个K-V对
	firstPair := b.GetFirstPair()
	//桶中没有元素值
	if firstPair == nil {
		//把传进来的元素存入firstValue
		b.firstValue.Store(p)
		atomic.AddUint64(&b.size, 1)
		return true, nil
	}
	var target Pair
	key := p.Key()
	//根据参数p遍历匹配是否有匹配的pair
	for v := firstPair; v != nil; v = v.Next() {
		if v.Key() == key {
			target = v
			break
		}
	}
	//找到了，替换
	if target != nil {
		target.SetElement(p.Element())
		return false, nil
	}
	//没有找到
	//把当前的第一个K-V对指定为参数值得单链目标
	//把旧firstPair连接到当前pair后面，当前pair成为新firstPair
	p.SetNext(firstPair)
	//用当前pair替换桶中的firstValue
	b.firstValue.Store(p)
	atomic.AddUint64(&b.size, 1)
	return true, nil
}

//根据key遍历单链表，找到pair
func (b *bucket) Get(key string) Pair {
	firstPair := b.GetFirstPair()
	if firstPair == nil {
		return nil
	}
	for v := firstPair; v != nil; v = v.Next() {
		if v.Key() == key {
			return v
		}
	}
	return nil
}

func (b *bucket) GetFirstPair() Pair {
	//查找firstValue里面是否真正存储有元素值，占位值不算
	if v := b.firstValue.Load(); v == nil {
		return nil
	} else if p, ok := v.(Pair); !ok || p == placeholder {
		return nil
	} else {
		return p
	}
}

func (b *bucket) Delete(key string, lock sync.Locker) bool {
	if lock != nil {
		lock.Lock()
		defer lock.Unlock()
	}
	firstPair := b.GetFirstPair()
	//桶里一个pair都没有
	if firstPair == nil {
		return false
	}
	//前导K-V列表，前导pair列表
	var prevPairs []Pair
	//目标pair
	var target Pair
	//目标pair的下一个pair，也就是要断开的pair
	var breakpoint Pair
	for v := firstPair; v != nil; v = v.Next() {
		if v.Key() == key {
			target = v
			//找到了目标pair，那么目标pair的下一个pair就是需要断开的pair
			breakpoint = v.Next()
			break
		}
		//将不匹配的pair放入前导pair列表
		prevPairs = append(prevPairs, v)
	}
	//没有找到目标pair，说明当前bucket中没有这个K，则无需删除
	if target == nil {
		return false
	}
	//将前导pair列表的最后一个和后续pair连接起来
	// -->1-->2-->3-->4-->5 其中3是目标pair，现在要把3删除掉
	//重新连接后的
	//-->1-->2-->4-->5
	newFirstPair := breakpoint
	//是否可以再优化
	for i := len(prevPairs) - 1; i >= 0; i-- {
		//依次（从最后一个pair开始）取出前导pair列表中的pair
		pairCopy := prevPairs[i].Copy()
		//把断裂后的头pair链接到前导pari的最后一个pair
		pairCopy.SetNext(newFirstPair)
		//链接了断裂的头pair的pair将成为新的断头pair，依次和前导pair链接起来
		newFirstPair = pairCopy
		//最后一个（前导pair列表的第一个pair将成为新的firstValue）
	}
	//将新的firstValue存入桶中firstValue中
	if newFirstPair != nil {
		//原子存储
		b.firstValue.Store(newFirstPair)
	} else {
		//原子的存入站位值
		b.firstValue.Store(placeholder)
	}
	//原子的减1
	//原子减法 atomic.AddUint64(&b.size, ^uint64(-N-1)) //其中N表示负数，总结起来就是对负数求绝对值再减1之后取其补码
	//原子的减5 ==> atomic.AddUint64(&b.size, ^uint64(4))
	atomic.AddUint64(&b.size, ^uint64(0))
	return true
}

func (b *bucket) Clear(lock sync.Locker) {
	if lock != nil {
		lock.Lock()
		defer lock.Unlock()
	}
	atomic.StoreUint64(&b.size, 0)
	//占位，保证下一次重新使用不会panic
	b.firstValue.Store(placeholder)
}

func (b *bucket) Size() uint64 {
	return atomic.LoadUint64(&b.size)
}

func (b *bucket) String() string {
	var buf bytes.Buffer
	buf.WriteString("[ ")
	for v := b.GetFirstPair(); v != nil; v = v.Next() {
		buf.WriteString(v.String())
		buf.WriteString(" ")
	}
	buf.WriteString("]")
	return buf.String()
}
