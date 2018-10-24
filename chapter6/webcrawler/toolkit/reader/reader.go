package reader

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
)

// MultipleReader 代表多重读取器的接口。
type MultipleReader interface {
	// Reader 用于获取一个可关闭读取器的实例。
	// 后者会持有本多重读取器中的数据。
	Reader() io.ReadCloser
}

//多重读取器接口
type MultipleReaderr interface {
	//reader 用于获取一个可关闭的读取器的实例
	//后者会持有本多重读取器中的数据
	Reader() io.ReadCloser
}

// myMultipleReader 代表多重读取器的实现类型。
type myMultipleReader struct {
	data []byte
}

//多重读取器的实现类型
type myMultipleReaderr struct {
	data []byte
}

// NewMultipleReader 用于新建并返回一个多重读取器的实例。
func NewMultipleReader(reader io.Reader) (MultipleReader, error) {
	var data []byte
	var err error
	if reader != nil {
		data, err = ioutil.ReadAll(reader)
		if err != nil {
			return nil, fmt.Errorf("multiple reader: couldn't create a new one: %s", err)
		}
	} else {
		data = []byte{}
	}
	return &myMultipleReader{
		data: data,
	}, nil
}

//创建一个多重读取器的实例
func NewMultipleReaderr(reader io.Reader) (MultipleReaderr, error) {
	var data []byte //uint8的别名类型
	var err error
	if reader != nil {
		data , err = ioutil.ReadAll(reader)
		if err != nil {
			//fmt.Errorf会返回一个实现了error接口的实际类型
			return nil, fmt.Errorf("multiple reader: couldnt create a new one: %s", err)
		}
	} else {
		data = []byte{}
	}
	return &myMultipleReaderr{
		data:data,
	},nil
}

func (rr *myMultipleReader) Reader() io.ReadCloser {
	//NopCloser根据一个reader接口封装成一个ReaderClose接口
	//bytes.NewReader()根据一个[]byte类型生成一个Reader接口的实现类型
	return ioutil.NopCloser(bytes.NewReader(rr.data))
}

func (rr *myMultipleReaderr) Reader() io.ReadCloser {
	return ioutil.NopCloser(bytes.NewReader(rr.data))
}