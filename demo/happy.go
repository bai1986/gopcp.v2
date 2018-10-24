package main

import (
	"context"
	"math/rand"
	"fmt"
	"time"
)

var maxNumber = 2018
var wish = 1024

func main () {
	var ctx, cancelFunc = context.WithCancel(context.Background())
	HAPPY:
		for  {
			happy := randNumber()
			checkHappy(happy, cancelFunc)
			select {
			case <- ctx.Done():
				fmt.Println("1024 happy ")
				break HAPPY
			default:
				fmt.Println("I have a Dream")
			}
			time.Sleep(500 * time.Millisecond)
		}
}

func randNumber() int {
	var number = rand.Intn(maxNumber)
	return number
}

func checkHappy(source int, do func())  {
	if source == wish {
		do()
	}
	return
}
