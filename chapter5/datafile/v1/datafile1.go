package v1

import (
	"errors"
	"io"
	"os"
	"sync"
)

// Data 代表数据的类型。
//别名类型
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
	woffset int64        // 写操作需要用到的偏移量。
	roffset int64        // 读操作需要用到的偏移量。
	wmutex  sync.Mutex   // 写操作需要用到的互斥锁。对写操作偏移量进行锁定
	rmutex  sync.Mutex   // 读操作需要用到的互斥锁。对读操作偏移量进行锁定
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
	//可以把指针类型当做一个接口类型返回（因为该指针类型已经实现了该接口类型）
	df := &myDataFile{f: f, dataLen: dataLen}
	return df, nil
}


func (df *myDataFile) Read() (rsn int64, d Data, err error) {
	// 读取并更新读偏移量。
	//一个读操作进来先更新偏移量，相当于占一个坑，那么其他go也调用Read时就不会和当前go冲突，从而去读下一个数据块
	var offset int64
	df.rmutex.Lock()
	offset = df.roffset
	df.roffset += int64(df.dataLen)
	df.rmutex.Unlock()

	//读取一个数据块。
	//根据当前Read要读取的偏移量获取区段号
	rsn = offset / int64(df.dataLen)
	bytes := make([]byte, df.dataLen)
	for {
		//对读写锁进行读锁定
		//在每次的for循环中对读写锁进行读锁定，在结束时对读写锁进行读解锁是为了让其他goroutine可以进行写锁定
		df.fmutex.RLock()
		//根据偏移量读取文件内容
		_, err = df.f.ReadAt(bytes, offset)
		if err != nil {
			if err == io.EOF {
				//对读写锁进行读解锁
				df.fmutex.RUnlock()
				//遇到文件边界，继续尝试读取，这么做是为了防止读写文件的goroutine数量不一致
				continue
			}
			//对读写锁进行读解锁
			df.fmutex.RUnlock()
			return
		}
		d = bytes
		//对读写锁进行读解锁
		df.fmutex.RUnlock()
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

	//写入一个数据块。
	wsn = offset / int64(df.dataLen)
	var bytes []byte
	if len(d) > int(df.dataLen) {
		//截取多余部分
		bytes = d[0:df.dataLen]
	} else {
		bytes = d
	}
	//对读写锁进行写锁定
	df.fmutex.Lock()
	//对读写锁进行写解锁
	defer df.fmutex.Unlock()
	//将数据在互斥锁环境下写入文件
	_, err = df.f.Write(bytes)
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
