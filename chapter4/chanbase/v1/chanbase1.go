package main

import (
	"fmt"
	"time"
)

var strChan = make(chan string, 3)

func main() {
	//syncChan1用来协调发送go和接收go之间的协调
	syncChan1 := make(chan struct{}, 1)
	//syncChan2用来协调其他go和主go的同步
	syncChan2 := make(chan struct{}, 2)
	// 用于演示接收操作。
	go func() {
		<-syncChan1
		fmt.Println("Received a sync signal and wait a second... [receiver]")
		time.Sleep(time.Second)
		for {
			if elem, ok := <-strChan; ok {
				fmt.Println("Received:", elem, "[receiver]")
			} else {
				break
			}
		}
		fmt.Println("Stopped. [receiver]")
		syncChan2 <- struct{}{}
	}()
	// 用于演示发送操作。
	go func() {
		for _, elem := range []string{"a", "b", "c", "d"} {
			strChan <- elem
			fmt.Println("Sent:", elem, "[sender]")
			if elem == "c" {
				syncChan1 <- struct{}{}
				fmt.Println("Sent a sync signal. [sender]")
			}
		}
		fmt.Println("Wait 2 seconds... [sender]")
		time.Sleep(time.Second * 2)
		close(strChan)
		syncChan2 <- struct{}{}
	}()
	<-syncChan2
	<-syncChan2
}
