package main

import (
	"fmt"
	"time"
)

func main() {
	timer := time.NewTimer(2 * time.Second)
	fmt.Printf("Present time: %v.\n", time.Now())
	//停止定时器之后，再去接收会造成死锁
	//fmt.Printf("Stop timer: %v.\n", timer.Stop())
	expirationTime := <-timer.C
	fmt.Printf("Expiration time: %v.\n", expirationTime)
	//返回false，因为定时器早已过期
	fmt.Printf("Stop timer: %v.\n", timer.Stop())
}
