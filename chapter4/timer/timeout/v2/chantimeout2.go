package main

import (
	"fmt"
	"time"
)

func main() {
	intChan := make(chan int, 1)
	go func() {
		for i := 0; i < 20; i++ {
			time.Sleep(time.Second * 2)
			intChan <- i
		}
		close(intChan)
	}()
	//如果定时器的时间间隔比上面发送的时间间隔长，select中定时器case始终不会触发，会被一直阻塞(第一个case有return的情况下)，
	// 如果去掉第一个select的return那么当第一个case接收完消息后，定时器才有机会接收
	timeout := time.Second * 2
	var timer *time.Timer
	//另一种创建Timer.C的方法
	//这种方式创建Timer.C会阻塞当前goroutie
	timeAfterC := <- time.After(100 * time.Millisecond)
	fmt.Println("timeAfterC:",timeAfterC)
	//过2秒执行函数，AfterFunc不会阻塞当前goroutine，因为定时器内部会新开一个go去执行回调函数
	time.AfterFunc(2 * time.Second, func() {
		//在新的goroutine中执行
		fmt.Println("2s exec AfterFunc")
	})
	for {
		if timer == nil {
			timer = time.NewTimer(timeout)
		} else {
			//重置定时器
			timer.Reset(timeout)
		}
		select {
		case e, ok := <-intChan:
			if !ok {
				fmt.Println("Received End.")
				return
			} else {
				fmt.Printf("Received: %v\n", e)
			}
		case c := <-timer.C:
			fmt.Println("Timer time:", c)
		}
	}

}
