package main

import (
	"fmt"
	"time"
)

var mapChan = make(chan map[string]int, 1)

func main() {
	syncChan := make(chan struct{}, 2)
	go func() { // 用于演示接收操作。
		for {
			if elem, ok := <-mapChan; ok {
				elem["count"]++
			} else {
				break
			}
		}
		fmt.Println("Stopped. [receiver]")
		syncChan <- struct{}{}
	}()
	go func() { // 用于演示发送操作。
		countMap := make(map[string]int)
		for i := 0; i < 5; i++ {
			//接收方会持有这个map，因为map是引用类型，所以相当于直接持有发送发的map
			mapChan <- countMap
			time.Sleep(time.Millisecond)
			fmt.Printf("The count map: %v. [sender]\n", countMap)
		}
		//已经关闭了的通道还是可以从里面接收到值，只是不能给关闭的通道发送值
		close(mapChan)
		syncChan <- struct{}{}
	}()
	<-syncChan
	<-syncChan
}

//演示结果
//The count map: map[count:1]. [sender]
//The count map: map[count:2]. [sender]
//The count map: map[count:3]. [sender]
//The count map: map[count:4]. [sender]
//The count map: map[count:5]. [sender]
//Stopped. [receiver]

