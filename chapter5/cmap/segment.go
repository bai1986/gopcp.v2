package cmap

import (
	"fmt"
	"sync"
	"sync/atomic"
)

// Segment 代表并发安全的散列段的接口。
type Segment interface {
	// Put 会根据参数放入一个键-元素对。
	// 第一个返回值表示是否新增了键-元素对。
	Put(p Pair) (bool, error)
	// Get 会根据给定参数返回对应的键-元素对。
	// 该方法会根据给定的键计算哈希值。
	Get(key string) Pair
	// GetWithHash 会根据给定参数返回对应的键-元素对。
	// 注意！参数keyHash应该是基于参数key计算得出哈希值。
	GetWithHash(key string, keyHash uint64) Pair
	// Delete 会删除指定键的键-元素对。
	// 若返回值为true则说明已删除，否则说明未找到该键。
	Delete(key string) bool
	// Size 用于获取当前段的尺寸（其中包含的散列桶的数量）。
	Size() uint64
}

// segment 代表并发安全的散列段的类型。
type segment struct {
	// buckets 代表散列桶切片。
	buckets []Bucket
	// bucketsLen 代表散列桶切片的长度。
	bucketsLen int
	// pairTotal 代表键-元素对总数。
	pairTotal uint64
	// pairRedistributor 代表键-元素对的再分布器。
	pairRedistributor PairRedistributor
	lock              sync.Mutex
}

// NewSegment 会创建一个Segment类型的实例。
func newSegment(
	bucketNumber int, pairRedistributor PairRedistributor) Segment {

	if bucketNumber <= 0 {
		bucketNumber = DEFAULT_BUCKET_NUMBER //16
	}
	if pairRedistributor == nil {
		pairRedistributor =
			newDefaultPairRedistributor(
				DEFAULT_BUCKET_LOAD_FACTOR, bucketNumber) //默认散列桶加载因子0.75，散列桶数量
	}
	buckets := make([]Bucket, bucketNumber)
	//初始化所有散列桶
	for i := 0; i < bucketNumber; i++ {
		buckets[i] = newBucket()
	}
	return &segment{
		buckets:           buckets,
		bucketsLen:        bucketNumber,
		pairRedistributor: pairRedistributor,
	}
}

func (s *segment) Put(p Pair) (bool, error) {
	// 对散列段上锁
	s.lock.Lock()
	// 根据Pair哈希值找到一个散列桶
	// 放入数据时随便分布，依据是key的哈希值取余桶的数量
	b := s.buckets[int(p.Hash()%uint64(s.bucketsLen))]
	// 把Pair放入散列桶
	ok, err := b.Put(p, nil)
	if ok {
		//pairTotal 总数原子加1
		newTotal := atomic.AddUint64(&s.pairTotal, 1)
		//触发散列段里面的散列桶K-V再分布，分布参数是散列段内键值对总数，桶的大小
		s.redistribute(newTotal, b.Size())
	}
	//散列段解锁
	s.lock.Unlock()
	return ok, err
}

func (s *segment) Get(key string) Pair {
	//内部调用hash取值方法
	return s.GetWithHash(key, hash(key))
}

func (s *segment) GetWithHash(key string, keyHash uint64) Pair {
	s.lock.Lock()
	// 根据哈希值取余散列桶数量
	b := s.buckets[int(keyHash%uint64(s.bucketsLen))]
	s.lock.Unlock()
	return b.Get(key)
}

func (s *segment) Delete(key string) bool {
	s.lock.Lock()
	// 根据哈希值取余散列桶数量
	b := s.buckets[int(hash(key)%uint64(s.bucketsLen))]
	//lock参数为nil时，表示外部自行保证并发安全
	ok := b.Delete(key, nil)
	if ok {
		//pairTotal原子减1
		newTotal := atomic.AddUint64(&s.pairTotal, ^uint64(0))
		//触发散列段里面的散列桶K-V再分布
		s.redistribute(newTotal, b.Size())
	}
	s.lock.Unlock()
	return ok
}

func (s *segment) Size() uint64 {
	return atomic.LoadUint64(&s.pairTotal)
}

// redistribute 会检查给定参数并设置相应的阈值和计数，
// 并在必要时重新分配所有散列桶中的所有键-元素对。
// 注意！必须在互斥锁的保护下调用本方法！
func (s *segment) redistribute(pairTotal uint64, bucketSize uint64) (err error) {
	defer func() {
		if p := recover(); p != nil {
			if pErr, ok := p.(error); ok {
				err = newPairRedistributorError(pErr.Error())
			} else {
				err = newPairRedistributorError(fmt.Sprintf("%s", p))
			}
		}
	}()
	//根据K-V总数   和  散列桶总数计算更新阀值
	s.pairRedistributor.UpdateThreshold(pairTotal, s.bucketsLen)
	//根据K-V总数和散列桶数量检查   散列桶状态
	bucketStatus := s.pairRedistributor.CheckBucketStatus(pairTotal, bucketSize)
	//根据散列桶状态和散列桶切片重新分布K-V
	newBuckets, changed := s.pairRedistributor.Redistribe(bucketStatus, s.buckets)
	//changed表示是否重新分布
	if changed {
		s.buckets = newBuckets //将重新分布后的散列桶赋值给当前散列段
		s.bucketsLen = len(s.buckets) //重新计算散列桶数量
	}
	return nil
}
