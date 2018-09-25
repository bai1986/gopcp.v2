package main

import (
	"fmt"
	"time"
)

// Counter 代表计数器的类型。
type Counter struct {
	count int
}

func (counter *Counter) String() string {
	return fmt.Sprintf("{count=:%d}", counter.count)
}

var mapChan = make(chan map[string]*Counter, 1)

func main() {
	syncChan := make(chan struct{}, 2)
	go func() { // 用于演示接收操作。
		for {
			if elem, ok := <-mapChan; ok {
				//counter是一个Counter结构体类型
				counter := elem["count"]
				counter.count++
			} else {
				break
			}
		}
		fmt.Println("Stopped. [receiver]")
		syncChan <- struct{}{}
	}()
	go func() { // 用于演示发送操作。
		countMap := map[string]*Counter{
			"count": &Counter{},
		}
		for i := 0; i < 5; i++ {
			mapChan <- countMap
			time.Sleep(time.Millisecond)
			fmt.Printf("The count map: %v. [sender]\n", countMap)
		}
		close(mapChan)
		syncChan <- struct{}{}
	}()
	<-syncChan
	<-syncChan
}
//演示结果
//The count map: map[count:{count=:1}]. [sender]
//The count map: map[count:{count=:2}]. [sender]
//The count map: map[count:{count=:3}]. [sender]
//The count map: map[count:{count=:4}]. [sender]
//The count map: map[count:{count=:5}]. [sender]
//Stopped. [receiver]


