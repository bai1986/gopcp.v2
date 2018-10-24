package buffer

import (
	"fmt"
	"sync"
	"sync/atomic"

	"gopcp.v2/chapter6/webcrawler/errors"
)

// Buffer 代表FIFO的缓冲器的接口类型。
type Buffer interface {
	// Cap 用于获取本缓冲器的容量。
	Cap() uint32
	// Len 用于获取本缓冲器中的数据数量。
	Len() uint32
	// Put 用于向缓冲器放入数据。
	// 注意！本方法应该是非阻塞的。
	// 若缓冲器已关闭则会直接返回非nil的错误值。
	Put(datum interface{}) (bool, error)
	// Get 用于从缓冲器获取器。
	// 注意！本方法应该是非阻塞的。
	// 若缓冲器已关闭则会直接返回非nil的错误值。
	Get() (interface{}, error)
	// Close 用于关闭缓冲器。
	// 若缓冲器之前已关闭则返回false，否则返回true。
	Close() bool
	// Closed 用于判断缓冲器是否已关闭。
	Closed() bool
}

//buffer 代表FIFO的缓冲器的接口类型
type Bufferr interface {
	//获取本缓冲器的容量
	Cap() uint32
	//获取本缓冲器的数据数量
	Len() uint32
	//向缓冲器放入数据，本方法是非阻塞的，若缓冲器已关闭则直接返回非nil错误
	Put(datum interface{}) (bool, error)
	//用于从缓冲器获取数据，本方法是非阻塞的
	Get() (interface{},error)
	//关闭缓冲器，若缓冲器之前已关闭则返回false，否则返回true
	Close() bool
	//用于判断缓冲器是否已经关闭
	Closed() bool
}

// myBuffer 代表缓冲器接口的实现类型。
type myBuffer struct {
	// ch 代表存放数据的通道。
	ch chan interface{}
	// closed 代表缓冲器的关闭状态：0-未关闭；1-已关闭。
	closed uint32
	// closingLock 代表为了消除因关闭缓冲器而产生的竞态条件的读写锁。
	closingLock sync.RWMutex
}
//缓冲器的实现类型
type myBufferr struct {
	//ch 代表存放数据的通道
	ch chan interface{}
	//缓冲器关闭状态
	closed uint32
	//读写锁用于消除因缓冲器关闭而产生竞态条件的读写锁
	closingLock sync.RWMutex
}

// NewBuffer 用于创建一个缓冲器。
// 参数size代表缓冲器的容量。
func NewBuffer(size uint32) (Buffer, error) {
	if size == 0 {
		errMsg := fmt.Sprintf("illegal size for buffer: %d", size)
		return nil, errors.NewIllegalParameterError(errMsg)
	}
	return &myBuffer{
		ch: make(chan interface{}, size),
	}, nil
}

func NewBufferr(size uint32) (Bufferr,error) {
	if size == 0 {
		errMsg := fmt.Sprintf("illegal size for buffer: %d",size)
		return nil,errors.NewIllegalParameterErrorr(errMsg)
	}
	return &myBufferr{ch:make(chan interface{},size)},nil
}

func (buf *myBuffer) Cap() uint32 {
	return uint32(cap(buf.ch))
}

func (buf *myBufferr) Cap() uint32 {
	return uint32(cap(buf.ch))
}

func (buf *myBuffer) Len() uint32 {
	return uint32(len(buf.ch))
}

func (buf *myBufferr) Len() uint32 {
	return uint32(len(buf.ch))
}

func (buf *myBuffer) Put(datum interface{}) (ok bool, err error) {
	//发送操作上读锁，表示可以并发发送
	//按照约定关闭操作上写写，保证同一时间只有一个GO能执行关闭函数，并且没有GO发送
	buf.closingLock.RLock()
	defer buf.closingLock.RUnlock()
	//一定要在GO拿到读锁后采取检测状态，而不是先拿状态再上读锁
	//这样会存在什么问题呢
	//如果在发送操作执行时，拿到状态为未关闭，由于这个时候并没有上锁，另一个GO在执行关闭操作，但是在发送操作拿到状态后才执行的操作
	//并且很快关闭逻辑执行完毕了，这时发送操作拿到了读锁，由于发送操作已经拿到状态为未关闭
	//那么这个时候会出现通道已经关闭或者正在关闭，发送操作发送数据到通道时会panic异常
	if buf.Closed() {
		return false, ErrClosedBuffer
	}
	select {
	case buf.ch <- datum:
		ok = true
	default:
		ok = false
	}
	return
}

func (buf *myBufferr) Put(data interface{}) (ok bool, err error) {
	buf.closingLock.RLock()
	defer buf.closingLock.RUnlock()
	if buf.Closed() {
		return false, ErrClosedBufferr
	}
	select {
	case buf.ch <- data:
		ok = true
	default:
		ok = false
	}
	return
}

func (buf *myBuffer) Get() (interface{}, error) {
	select {
	case datum, ok := <-buf.ch:
		if !ok {
			return nil, ErrClosedBuffer
		}
		return datum, nil
	default:
		return nil, nil
	}
}

func (buf *myBufferr) Get() (interface{},error) {
	select {
	case data ,ok := <- buf.ch:
		if !ok {
			return nil, ErrClosedBufferr
		}
		return data,nil
	default:
		return nil,nil
	}
}

func (buf *myBuffer) Close() bool {
	//先设置关闭状态
	if atomic.CompareAndSwapUint32(&buf.closed, 0, 1) {
		//用写锁，阻止其他go向ch发送数据
		buf.closingLock.Lock()
		close(buf.ch)
		buf.closingLock.Unlock()
		return true
	}
	return false
}

//这里为什么要在设置关闭状态成功的逻辑里面，再上写锁去关闭通道呢
//比如发送go进入执行逻辑，检测到通道状态为未关闭，进入发送逻辑部分，这是另一个go执行关闭逻辑，进去检测关闭状态，发现没有关闭，然后改变状态
//也成功了，当然设置了状态后可以阻止后面的go在进入发送逻辑，但是在关闭之前已经进入的go单单靠状态是阻止不了的
//所以要上写锁，然后去关闭，那么发送函数应该在上读锁之后检测状态，如果状态为已关闭则不要再发送
func (buf *myBufferr) Close() bool {
	//先设置关闭状态
	if atomic.CompareAndSwapUint32(&buf.closed,0,1) {
		//在关闭执行过程中，不允许发送
		buf.closingLock.Lock()
		close(buf.ch)
		buf.closingLock.Unlock()
		return true
	}
	return false
}

func (buf *myBuffer) Closed() bool {
	if atomic.LoadUint32(&buf.closed) == 0 {
		return false
	}
	return true
}

func (buf *myBufferr) Closed() bool {
	if atomic.LoadUint32(&buf.closed) == 0 {
		return false
	}
	return true
}
