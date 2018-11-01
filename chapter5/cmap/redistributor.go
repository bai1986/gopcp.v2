package cmap

import "sync/atomic"

// BucketStatus 代表散列桶状态的类型。
type BucketStatus uint8

const (
	// BUCKET_STATUS_NORMAL 代表散列桶正常。
	BUCKET_STATUS_NORMAL BucketStatus = 0
	// BUCKET_STATUS_UNDERWEIGHT 代表散列桶过轻。
	BUCKET_STATUS_UNDERWEIGHT BucketStatus = 1
	// BUCKET_STATUS_OVERWEIGHT 代表散列桶过重。
	BUCKET_STATUS_OVERWEIGHT BucketStatus = 2
)

// PairRedistributor 代表针对键-元素对的再分布器。
// 用于当散列段内的键-元素对分布不均时进行重新分布。
type PairRedistributor interface {
	//  UpdateThreshold 会根据键-元素对总数和散列桶总数计算并更新阈值。
	UpdateThreshold(pairTotal uint64, bucketNumber int)
	// CheckBucketStatus 用于检查散列桶的状态。
	CheckBucketStatus(pairTotal uint64, bucketSize uint64) (bucketStatus BucketStatus)
	// Redistribe 用于实施键-元素对的再分布。
	Redistribe(bucketStatus BucketStatus, buckets []Bucket) (newBuckets []Bucket, changed bool)
}

// myPairRedistributor 代表PairRedistributor的默认实现类型。
type myPairRedistributor struct {
	// loadFactor 代表装载因子。
	loadFactor float64
	// upperThreshold 代表散列桶重量的上阈限。
	// 当某个散列桶的尺寸增至此值时会触发再散列。
	upperThreshold uint64
	// overweightBucketCount 代表过重的散列桶的计数。
	overweightBucketCount uint64
	// emptyBucketCount 代表空的散列桶的计数。
	emptyBucketCount uint64
}

// newDefaultPairRedistributor 会创建一个PairRedistributor类型的实例。
// 参数loadFactor代表散列桶的负载因子。
// 参数bucketNumber代表散列桶的数量。
func newDefaultPairRedistributor(loadFactor float64, bucketNumber int) PairRedistributor {
	//  散列桶加载因子
	if loadFactor <= 0 {
		loadFactor = DEFAULT_BUCKET_LOAD_FACTOR
	}
	pr := &myPairRedistributor{}
	pr.loadFactor = loadFactor
	pr.UpdateThreshold(0, bucketNumber)
	return pr
}

// bucketCountTemplate 代表调试用散列桶状态信息模板。
var bucketCountTemplate = `Bucket count: 
    pairTotal: %d
    bucketNumber: %d
    average: %f
    upperThreshold: %d
    emptyBucketCount: %d

`

func (pr *myPairRedistributor) UpdateThreshold(pairTotal uint64, bucketNumber int) {
	var average float64
	average = float64(pairTotal / uint64(bucketNumber))
	//如果每个桶里面平均存储的K-V对数量小于100，那么就将平均值设定为100
	if average < 100 {
		average = 100
	}
	// defer func() {
	// 	fmt.Printf(bucketCountTemplate,
	// 		pairTotal,
	// 		bucketNumber,
	// 		average,
	// 		atomic.LoadUint64(&pr.upperThreshold),
	// 		atomic.LoadUint64(&pr.emptyBucketCount))
	// }()
	//更新散列桶重分布阀值 === 平均值 * 加载因子0.75
	atomic.StoreUint64(&pr.upperThreshold, uint64(average*pr.loadFactor))
}

// bucketStatusTemplate 代表调试用散列桶状态信息模板。
var bucketStatusTemplate = `Check bucket status: 
    pairTotal: %d
    bucketSize: %d
    upperThreshold: %d
    overweightBucketCount: %d
    emptyBucketCount: %d
    bucketStatus: %d
	
`

func (pr *myPairRedistributor) CheckBucketStatus(pairTotal uint64, bucketSize uint64) (bucketStatus BucketStatus) {
	// defer func() {
	// 	fmt.Printf(bucketStatusTemplate,
	// 		pairTotal,
	// 		bucketSize,
	// 		atomic.LoadUint64(&pr.upperThreshold),
	// 		atomic.LoadUint64(&pr.overweightBucketCount),
	// 		atomic.LoadUint64(&pr.emptyBucketCount),
	// 		bucketStatus)
	// }()
	//如果散列桶大小 已经大于预定义最大值，或者大于更新阀值限制
	if bucketSize > DEFAULT_BUCKET_MAX_SIZE ||
		bucketSize >= atomic.LoadUint64(&pr.upperThreshold) {
		//那么就将 "过重桶的数量"加1
		atomic.AddUint64(&pr.overweightBucketCount, 1)
		//然后将散列桶状态标记为过重
		bucketStatus = BUCKET_STATUS_OVERWEIGHT
		return
	}
	//当前桶的大小为0，则将闲置桶的数量加1
	if bucketSize == 0 {
		atomic.AddUint64(&pr.emptyBucketCount, 1)
	}
	return
}

// redistributionTemplate 代表重新分配信息模板。
var redistributionTemplate = `Redistributing: 
    bucketStatus: %d
    currentNumber: %d
    newNumber: %d

`

func (pr *myPairRedistributor) Redistribe(
	bucketStatus BucketStatus, buckets []Bucket) (newBuckets []Bucket, changed bool) {
	//当前散列桶的数量
	currentNumber := uint64(len(buckets))
	newNumber := currentNumber
	// defer func() {
	// 	fmt.Printf(redistributionTemplate,
	// 		bucketStatus,
	// 		currentNumber,
	// 		newNumber)
	// }()
	switch bucketStatus {
	//如果散列桶状态为超重
	case BUCKET_STATUS_OVERWEIGHT:
		//如果超重的散列桶数量小于当前散列桶总数四分之一
		if atomic.LoadUint64(&pr.overweightBucketCount)*4 < currentNumber {
			return nil, false
		}
		//如果超重散列桶数量大于等于当前散列桶总数四分之一
		//开始扩容
		//那么新的散列桶数量就是当前散列桶数量的2倍
		newNumber = currentNumber << 1
		//如果散列桶状态为过轻
	case BUCKET_STATUS_UNDERWEIGHT:
		//如果当前散列桶数量小于100或者闲置桶数量小于当前散列桶总数的四分之一,则返回
		if currentNumber < 100 ||
			atomic.LoadUint64(&pr.emptyBucketCount)*4 < currentNumber {
			return nil, false
		}
		//散列桶缩容
		newNumber = currentNumber >> 1
		//最低散列桶数量为2
		if newNumber < 2 {
			newNumber = 2
		}
	default:
		return nil, false
	}
	//如果计算出来新散列桶数量 和 当前散列桶数量相等
	if newNumber == currentNumber {
		//重置过重散列桶计数
		atomic.StoreUint64(&pr.overweightBucketCount, 0)
		//重置闲置散列桶计数
		atomic.StoreUint64(&pr.emptyBucketCount, 0)
		return nil, false
	}
	var pairs []Pair
	//把散列桶里面所有的K-V取出来
	for _, b := range buckets {
		for e := b.GetFirstPair(); e != nil; e = e.Next() {
			pairs = append(pairs, e)
		}
	}
	//开始规划新散列桶
	//如果重新计算的新散列桶数量大于当前散列桶数量
	if newNumber > currentNumber {
		//清空所有散列桶
		for i := uint64(0); i < currentNumber; i++ {
			buckets[i].Clear(nil)
		}
		//补齐散列桶（新补齐的桶都是空桶）,扩容散列桶
		for j := newNumber - currentNumber; j > 0; j-- {
			buckets = append(buckets, newBucket())
		}
	} else {
		//如果重新计算的散列桶小于当前散列桶，收缩散列桶
		buckets = make([]Bucket, newNumber)
		for i := uint64(0); i < newNumber; i++ {
			buckets[i] = newBucket()
		}
	}
	//重新计算后的散列桶都是空桶
	var count int
	//将所有K-V对根据哈希取余散列桶数量（重新计算后的）再均匀分布到上面构建的空桶中
	for _, p := range pairs {
		index := int(p.Hash() % newNumber)
		b := buckets[index]
		b.Put(p, nil)
		count++
	}
	//重置过重散列桶计数
	atomic.StoreUint64(&pr.overweightBucketCount, 0)
	//重置闲置散列桶计数
	atomic.StoreUint64(&pr.emptyBucketCount, 0)
	return buckets, true
}
