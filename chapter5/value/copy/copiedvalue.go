package main

import (
	"sync/atomic"
	"fmt"
)

func main() {
	var countVal atomic.Value
	countVal.Store([]int{1, 3, 5, 7})
	//Value是一个结构体类型，副本并不会受并发安全保护
	//但是*Value类型，副本依然指向原地址，因此可以并发安全
	anotherStore(countVal)
	fmt.Printf("The anotherStore value: %+v \n", countVal.Load())
	anotherStoreByPointer(&countVal)
	fmt.Printf("The anotherStoreByPointer value: %+v \n", countVal.Load())
}

func anotherStore(countVal atomic.Value) {
	//Store()方法不能传入nil
	//Store()方法传入的值必须和上一次传入的值类型相同（如果有上一次的话）
	countVal.Store([]int{2, 4, 6, 8})
}

func anotherStoreByPointer(countVal *atomic.Value) {
	countVal.Store([]int{2, 4, 6, 8})
}
