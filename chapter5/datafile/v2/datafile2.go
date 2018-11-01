package v2

import (
	"errors"
	"io"
	"os"
	"sync"
)

// Data 代表数据的类型。
type Data []byte

// DataFile 代表数据文件的接口类型。
type DataFile interface {
	// Read 会读取一个数据块。
	Read() (rsn int64, d Data, err error)
	// Write 会写入一个数据块。
	Write(d Data) (wsn int64, err error)
	// RSN 会获取最后读取的数据块的序列号。
	RSN() int64
	// WSN 会获取最后写入的数据块的序列号。
	WSN() int64
	// DataLen 会获取数据块的长度。
	DataLen() uint32
	// Close 会关闭数据文件。
	Close() error
}

// myDataFile 代表数据文件的实现类型。
type myDataFile struct {
	f       *os.File     // 文件。
	fmutex  sync.RWMutex // 被用于文件的读写锁。
	rcond   *sync.Cond   //读操作需要用到的条件变量
	woffset int64        // 写操作需要用到的偏移量。
	roffset int64        // 读操作需要用到的偏移量。
	wmutex  sync.Mutex   // 写操作需要用到的互斥锁。
	rmutex  sync.Mutex   // 读操作需要用到的互斥锁。
	dataLen uint32       // 数据块长度。
}

// NewDataFile 会新建一个数据文件的实例。
func NewDataFile(path string, dataLen uint32) (DataFile, error) {
	f, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	if dataLen == 0 {
		return nil, errors.New("Invalid data length!")
	}
	df := &myDataFile{f: f, dataLen: dataLen}
	//针对读写锁的读操作实施信号通知
	//sync.NewCond(&df.wmutex)
	//RLocker()获取的是读锁
	df.rcond = sync.NewCond(df.fmutex.RLocker())
	return df, nil
}

func (df *myDataFile) Read() (rsn int64, d Data, err error) {
	// 读取并更新读偏移量。
	var offset int64
	df.rmutex.Lock()
	offset = df.roffset
	df.roffset += int64(df.dataLen)
	df.rmutex.Unlock()

	//读取一个数据块。
	rsn = offset / int64(df.dataLen)
	bytes := make([]byte, df.dataLen)
	//rcond.Wait()之前要先对读写锁进行读锁定
	df.fmutex.RLock()
	//及时对读写锁进行读解锁
	defer df.fmutex.RUnlock()
	for {
		_, err = df.f.ReadAt(bytes, offset)
		if err != nil {
			for err == io.EOF {
				//一定要在调用rcond.Wait()方法之前锁定与之关联的读锁，否则会引发不可恢复的panic
				//放弃尝试读锁定，等待通知
				df.rcond.Wait() //Wait()方法返回之前会重新锁定与之关联的那个读锁
				//收到通知，并且成功进行了读锁定
				continue
			}
			return
		}
		d = bytes
		return
	}
}

func (df *myDataFile) Write(d Data) (wsn int64, err error) {
	// 读取并更新写偏移量。
	var offset int64
	df.wmutex.Lock()
	offset = df.woffset
	df.woffset += int64(df.dataLen)
	df.wmutex.Unlock()

	//写入一个数据块。wsn表示当开始写入的区段号
	wsn = offset / int64(df.dataLen)
	var bytes []byte
	if len(d) > int(df.dataLen) {
		bytes = d[0:df.dataLen]
	} else {
		bytes = d
	}
	//与下面的Signal()方法没有关联
	df.fmutex.Lock()
	//与下面的Signal()方法没有关联
	defer df.fmutex.Unlock()
	_, err = df.f.Write(bytes)
	//由于一个数据块只能有某个读操作读取，所以只是使用条件变量的Signal方法去通知某个为此等待的Wait()方法，并以此唤醒某一个相关的goroutine
	df.rcond.Signal()
	return
}

func (df *myDataFile) RSN() int64 {
	df.rmutex.Lock()
	defer df.rmutex.Unlock()
	return df.roffset / int64(df.dataLen)
}

func (df *myDataFile) WSN() int64 {
	df.wmutex.Lock()
	defer df.wmutex.Unlock()
	return df.woffset / int64(df.dataLen)
}

func (df *myDataFile) DataLen() uint32 {
	return df.dataLen
}

func (df *myDataFile) Close() error {
	if df.f == nil {
		return nil
	}
	return df.f.Close()
}
