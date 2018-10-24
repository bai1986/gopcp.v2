package buffer

import (
	"fmt"
	"sync"
	"sync/atomic"
	"gopcp.v2/chapter6/webcrawler/errors"
)

// Pool 代表数据缓冲池的接口类型。
type Pool interface {
	// BufferCap 用于获取池中缓冲器的统一容量。
	BufferCap() uint32
	// MaxBufferNumber 用于获取池中缓冲器的最大数量。
	MaxBufferNumber() uint32
	// BufferNumber 用于获取池中缓冲器的数量。
	BufferNumber() uint32
	// Total 用于获取缓冲池中数据的总数。
	Total() uint64
	// Put 用于向缓冲池放入数据。
	// 注意！本方法应该是阻塞的。
	// 若缓冲池已关闭则会直接返回非nil的错误值。
	Put(datum interface{}) error
	// Get 用于从缓冲池获取数据。
	// 注意！本方法应该是阻塞的。
	// 若缓冲池已关闭则会直接返回非nil的错误值。
	Get() (datum interface{}, err error)
	// Close 用于关闭缓冲池。
	// 若缓冲池之前已关闭则返回false，否则返回true。
	Close() bool
	// Closed 用于判断缓冲池是否已关闭。
	Closed() bool
}

//代表数据缓冲池的接口类型
type Pooll interface {
	//获取缓冲池中缓冲器的统一容量
	BufferCap() uint32
	//用于获取缓冲池中缓冲器的最大数量
	MaxBufferNumber() uint32
	//用于获取缓冲池中缓冲器的数量
	BufferNumber() uint32
	//用于获取缓冲池中数据的总数
	Total() uint32
	//向缓冲池中放入数据
	//此方法是非阻塞的
	//若缓冲池已关闭则会直接返回非nil的错误值
	Put(data interface{}) error
	//从缓冲池中获取数据
	//此方法是非阻塞的
	//若缓冲池已关闭则会直接返回非nil的错误值
	Get() (dataum interface{}, err error)
	//关闭缓冲池
	//若在关闭之前已关闭则返回false，否则返回true
	Close() bool
	//判断缓冲池是否已关闭
	Closed() bool
}


// myPool 代表数据缓冲池接口的实现类型。
type myPool struct {
	// bufferCap 代表缓冲器的统一容量。
	bufferCap uint32
	// maxBufferNumber 代表缓冲器的最大数量。
	maxBufferNumber uint32
	// bufferNumber 代表缓冲器的实际数量。
	bufferNumber uint32
	// total 代表池中数据的总数。
	total uint64
	// bufCh 代表存放缓冲器的通道。
	bufCh chan Buffer
	// closed 代表缓冲池的关闭状态：0-未关闭；1-已关闭。
	closed uint32
	// lock 代表保护内部共享资源的读写锁。
	rwlock sync.RWMutex
}
//缓冲池接口的实现类型
type myPooll struct {
	//代表缓冲器的统一容量，在初始化缓冲器时有用
	bufferCap uint32
	//缓冲器的最大数量
	maxBufferNumber uint32
	//缓冲器的实际数量
	bufferNumber uint32
	//池中数据总数
	total uint64
	//缓冲器通道
	bufCh chan Bufferr
	//缓冲池的状态，0表示未关闭
	closed uint32
	//保护内部资源的读写锁
	rwlock sync.RWMutex
}

// NewPool 用于创建一个数据缓冲池。
// 参数bufferCap代表池内缓冲器的统一容量。
// 参数maxBufferNumber代表池中最多包含的缓冲器的数量。
func NewPool(
	bufferCap uint32,
	maxBufferNumber uint32) (Pool, error) {
	if bufferCap == 0 {
		errMsg := fmt.Sprintf("illegal buffer cap for buffer pool: %d", bufferCap)
		return nil, errors.NewIllegalParameterError(errMsg)
	}
	if maxBufferNumber == 0 {
		errMsg := fmt.Sprintf("illegal max buffer number for buffer pool: %d", maxBufferNumber)
		return nil, errors.NewIllegalParameterError(errMsg)
	}
	bufCh := make(chan Buffer, maxBufferNumber)
	//缓冲池预热
	buf, _ := NewBuffer(bufferCap)
	bufCh <- buf
	return &myPool{
		bufferCap:       bufferCap,
		maxBufferNumber: maxBufferNumber,
		bufferNumber:    1,
		bufCh:           bufCh,
	}, nil
}
//用于创建一个数据缓冲池
//buffercap代表缓冲器的统一容量
//maxbuffernumber代表池中最多包含的缓冲器数量

func NewPooll(
	bufferCap uint32,
	maxBufferNumber uint32) (Pool, error) {
	if bufferCap == 0 {
		errMsg := fmt.Sprintf("illegal buffer cap for buffer pool:%d", bufferCap)
		return nil,errors.NewIllegalParameterErrorr(errMsg)
	}
	if maxBufferNumber == 0 {
		errMsg := fmt.Sprintf("illegal max buffer number for buffer pool :%d",maxBufferNumber)
		return nil, errors.NewIllegalParameterErrorr(errMsg)
	}
	bufCh := make(chan Bufferr, maxBufferNumber)
	//缓冲池预热
	buf , _ := NewBufferr(bufferCap)
	bufCh <- buf
	return &myPooll{
		bufferCap:bufferCap,
		maxBufferNumber:maxBufferNumber,
		bufferNumber:1,
		bufCh:bufCh,
	},nil
}

func (pool *myPool) BufferCap() uint32 {
	//bufferCap在初始化之后不会变，所以不用枷锁
	return pool.bufferCap
}

func (pool *myPooll) BufferCap() uint32 {
	return pool.bufferCap
}

func (pool *myPool) MaxBufferNumber() uint32 {
	//最大buffer数量在初始化后不会改变，也就是说只能读，没有写入
	return pool.maxBufferNumber
}

func (pool *myPooll) MaxBufferNumber() uint32 {
	return pool.maxBufferNumber
}

func (pool *myPool) BufferNumber() uint32 {
	//bufferNumber是实时的
	return atomic.LoadUint32(&pool.bufferNumber)
}

func (pool *myPooll) BufferNumber() uint32 {
	return atomic.LoadUint32(&pool.bufferNumber)
}

func (pool *myPool) Total() uint64 {
	return atomic.LoadUint64(&pool.total)
}

func (pool *myPooll) Total() uint64 {
	return atomic.LoadUint64(&pool.total)
}

func (pool *myPool) Put(datum interface{}) (err error) {
	if pool.Closed() {
		return ErrClosedBufferPool
	}
	var count uint32
	maxCount := pool.BufferNumber() * 5
	var ok bool
	for buf := range pool.bufCh {
		ok, err = pool.putData(buf, datum, &count, maxCount)
		if ok || err != nil {
			break
		}
	}
	return
}

func (pool *myPooll) Put(data interface{}) (err error) {
	if pool.Closed() {
		return ErrClosedBufferPooll
	}
	var count uint32
	maxCount := pool.BufferNumber() * 5
	var ok bool
	for buf := range pool.bufCh {
		ok , err = pool.putData(buf,data,&count,maxCount)
		if ok || err != nil {
			break
		}
	}
	return
}

func (pool *myPooll) putData(
	buf Bufferr,data interface{},count *uint32,maxCount uint32) (ok bool,err error) {
	if pool.Closed() {
		return false,ErrClosedBufferPooll
	}
	defer func() {
		pool.rwlock.RLock()
		if pool.Closed() {
			atomic.AddUint32(&pool.bufferNumber,^uint32(0))
			err = ErrClosedBufferPooll
		} else {
			pool.bufCh <- buf
		}
	}()

	//尝试将数据放入缓冲器
	ok , err = buf.Put(data)
	//如果成功放入，则返回
	if ok {
		atomic.AddUint64(&pool.total,1)
		return
	}
	//如果出现错误,返回
	if err != nil {
		return
	}
	//如果没有成功放入，但是也没报错（put方法是异步的）

	//若因缓冲器已满而未被放入数据就递增计数
	(*count)++
	//若尝试向缓冲器放入数据失败的次数达到最大值,并且池中缓冲器的数量未达到最大值，那么就尝试创建一个新的缓冲器，先把数据放入缓冲器，然后把缓冲器放入缓冲池中
	if *count >= maxCount && pool.BufferNumber() < pool.MaxBufferNumber() {
		pool.rwlock.Lock()
		if pool.BufferNumber() < pool.MaxBufferNumber() {
			//如果缓冲池已经被关闭了
			if pool.Closed() {
				pool.rwlock.Unlock()
				return
			}
			newBuf ,_ := NewBufferr(pool.bufferCap)
			newBuf.Put(data)
			pool.bufCh <- newBuf
			atomic.AddUint32(&pool.bufferNumber,1)
			atomic.AddUint64(&pool.total,1)
			ok = true
		}
		pool.rwlock.Unlock()
		*count = 0
	}
	return
}

// putData 用于向给定的缓冲器放入数据，并在必要时把缓冲器归还给池。
func (pool *myPool) putData(
	buf Buffer, datum interface{}, count *uint32, maxCount uint32) (ok bool, err error) {
	if pool.Closed() {
		return false, ErrClosedBufferPool
	}
	defer func() {
		pool.rwlock.RLock()
		if pool.Closed() {
			atomic.AddUint32(&pool.bufferNumber, ^uint32(0))
			err = ErrClosedBufferPool
		} else {
			pool.bufCh <- buf
		}
		pool.rwlock.RUnlock()
	}()
	ok, err = buf.Put(datum)
	if ok {
		atomic.AddUint64(&pool.total, 1)
		return
	}
	if err != nil {
		return
	}
	// 若因缓冲器已满而未放入数据就递增计数。
	(*count)++
	// 如果尝试向缓冲器放入数据的失败次数达到阈值，
	// 并且池中缓冲器的数量未达到最大值，
	// 那么就尝试创建一个新的缓冲器，先放入数据再把它放入池。
	if *count >= maxCount &&
		pool.BufferNumber() < pool.MaxBufferNumber() {
		pool.rwlock.Lock()
		if pool.BufferNumber() < pool.MaxBufferNumber() {
			if pool.Closed() {
				pool.rwlock.Unlock()
				return
			}
			newBuf, _ := NewBuffer(pool.bufferCap)
			newBuf.Put(datum)
			pool.bufCh <- newBuf
			atomic.AddUint32(&pool.bufferNumber, 1)
			atomic.AddUint64(&pool.total, 1)
			ok = true
		}
		pool.rwlock.Unlock()
		*count = 0
	}
	return
}

func (pool *myPool) Get() (datum interface{}, err error) {
	if pool.Closed() {
		return nil, ErrClosedBufferPool
	}
	var count uint32
	maxCount := pool.BufferNumber() * 10
	for buf := range pool.bufCh {
		datum, err = pool.getData(buf, &count, maxCount)
		if datum != nil || err != nil {
			break
		}
	}
	return
}

func (pool *myPooll) Get() (data interface{}, err error) {
	if pool.Closed() {
		return nil,ErrClosedBufferPooll
	}
	var count uint32
	maxCount := pool.BufferNumber() * 10
	for buf := range pool.bufCh {
		data, err = pool.getData(buf,&count,maxCount)
		if data != nil || err != nil {
			break
		}
	}
	return
}

func (pool *myPooll) getData(
	buf Bufferr,count *uint32,maxCount uint32) (data interface{},err error) {
	if pool.Closed() {
		return nil, ErrClosedBufferPooll
	}
	defer func() {
		//如果尝试从缓冲器获取数据失败次数达到最大值
		//同事当前缓冲器已空，且缓冲器数量大于1，为什么要大于1,要保证缓冲池中至少有一个缓冲器
		//那么直接关掉当前缓冲器，并且不归还
		if *count >= maxCount && buf.Len() ==0 && pool.BufferNumber() > 1 {
			buf.Close()
			atomic.AddUint32(&pool.bufferNumber,^uint32(0))
			*count = 0
			return
		}
		pool.rwlock.RLock()
		if pool.Closed() {
			atomic.AddUint32(&pool.bufferNumber, ^uint32(0))
			err = ErrClosedBufferPooll
		} else {
			pool.bufCh <- buf
		}
		pool.rwlock.RUnlock()
	}()


	data, err = buf.Get()
	if data != nil {
		atomic.AddUint64(&pool.total, ^uint64(0))
		return
	}
	if err != nil {
		return
	}
	// 若因缓冲器已空未取出数据就递增计数。
	(*count) ++
	return
}

// getData 用于从给定的缓冲器获取数据，并在必要时把缓冲器归还给池。
func (pool *myPool) getData(
	buf Buffer, count *uint32, maxCount uint32) (datum interface{}, err error) {
	if pool.Closed() {
		return nil, ErrClosedBufferPool
	}
	defer func() {
		// 如果尝试从缓冲器获取数据的失败次数达到阈值，
		// 同时当前缓冲器已空且池中缓冲器的数量大于1，
		// 那么就直接关掉当前缓冲器，并不归还给池。
		if *count >= maxCount &&
			buf.Len() == 0 &&
			pool.BufferNumber() > 1 {
			buf.Close()
			atomic.AddUint32(&pool.bufferNumber, ^uint32(0))
			*count = 0
			return
		}
		pool.rwlock.RLock()
		if pool.Closed() {
			atomic.AddUint32(&pool.bufferNumber, ^uint32(0))
			err = ErrClosedBufferPool
		} else {
			pool.bufCh <- buf
		}
		pool.rwlock.RUnlock()
	}()
	datum, err = buf.Get()
	if datum != nil {
		atomic.AddUint64(&pool.total, ^uint64(0))
		return
	}
	if err != nil {
		return
	}
	// 若因缓冲器已空未取出数据就递增计数。
	(*count)++
	return
}

func (pool *myPool) Close() bool {
	if !atomic.CompareAndSwapUint32(&pool.closed, 0, 1) {
		return false
	}
	pool.rwlock.Lock()
	defer pool.rwlock.Unlock()
	close(pool.bufCh)
	for buf := range pool.bufCh {
		buf.Close()
	}
	return true
}

func (pool *myPooll) Close() bool {
	//如果设置缓冲池关闭状态失败则返回false
	if !atomic.CompareAndSwapUint32(&pool.closed,0,1) {
		return false
	}
	//先将缓冲池状态设置为关闭，是为了阻止后面的GO进来
	//这里上写锁是为了避免已经进入put流程的GO往池中发送数据
	pool.rwlock.Lock()
	defer pool.rwlock.Unlock()
	//先关闭pool缓冲池，是为了防止继续往池中增加缓冲器
	close(pool.bufCh)
	for buf := range pool.bufCh {
		buf.Close()
	}
	return true
}

func (pool *myPool) Closed() bool {
	if atomic.LoadUint32(&pool.closed) == 1 {
		return true
	}
	return false
}

func (pool *myPooll) Closed() bool {
	if atomic.LoadUint32(&pool.closed) == 1 {
		return true
	}
	return false
}