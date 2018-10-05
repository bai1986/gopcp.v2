package cow

import (
	"errors"
	"fmt"
	"sync/atomic"
)

// ConcurrentArray 代表并发安全的整数数组接口。
type ConcurrentArray interface {
	// Set 用于设置指定索引上的元素值。
	Set(index uint32, elem int) (err error)
	// Get 用于获取指定索引上的元素值。
	Get(index uint32) (elem int, err error)
	// Len 用于获取数组的长度。
	Len() uint32
}

// intArray 代表ConcurrentArray接口的实现类型。
type intArray struct {
	length uint32
	val    atomic.Value
}

// NewConcurrentArray 会创建一个ConcurrentArray类型值。
func NewConcurrentArray(length uint32) ConcurrentArray {
	array := intArray{}
	array.length = length
	//make([]int, array.length)返回一个slice指针
	array.val.Store(make([]int, array.length))
	return &array
}

func (array *intArray) Set(index uint32, elem int) (err error) {
	//检查index索引是否合法
	if err = array.checkIndex(index); err != nil {
		return
	}
	//如果val中没有一个值，直接返回（在Set()中说明实例没有进行初始化）
	if err = array.checkValue(); err != nil {
		return
	}

	// 不要这样做！否则会形成竞态条件！
	//array.val.Load().([]int)返回的是slice，在并发读写时会出现竞态
	//无论在原子值Value中存储什么类型的值，只要新值需要根据旧值计算得出，那么并发写的时候就可能出现问题
	//oldArray := array.val.Load().([]int)
	//oldArray[index] = elem
	//array.val.Store(oldArray)

	//利用原子值实现了COW（copy-on-write）写时复制算法
	//当要修改值时生成并修改副本，然后再用副本完全替换原值
	newArray := make([]int, array.length)
	copy(newArray, array.val.Load().([]int))
	newArray[index] = elem
	array.val.Store(newArray)
	return
}

func (array *intArray) Get(index uint32) (elem int, err error) {
	//检查index索引是否合法
	if err = array.checkIndex(index); err != nil {
		return
	}
	//如果val中没有一个值，直接返回
	if err = array.checkValue(); err != nil {
		return
	}
	//取出val并类型断言，然后取出指定索引位置的值
	elem = array.val.Load().([]int)[index]
	return
}

func (array *intArray) Len() uint32 {
	return array.length
}

// checkIndex 用于检查索引的有效性。
func (array *intArray) checkIndex(index uint32) error {
	if index >= array.length {
		return fmt.Errorf("Index out of range [0, %d)!", array.length)
	}
	return nil
}

// checkValue 用于检查原子值中是否已存有值。
func (array *intArray) checkValue() error {
	v := array.val.Load()
	if v == nil {
		return errors.New("Invalid int array!")
	}
	return nil
}
