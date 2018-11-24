package main

import (
	"fmt"
	"time"
)

func main() {

	time1 := time.Now().Unix()
	time2 := time.Now().UnixNano()
	time3 := time.Now().String()
	fmt.Printf("time1: %v\n",time1)
	fmt.Printf("time2: %v\n",time2)
	fmt.Printf("time3: %v\n",time3)
}

