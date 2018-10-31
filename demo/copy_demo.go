package main

import "fmt"

func main() {
	var slice1 = make([]int, 3)
	var slice2 = make([]int, 3)
	slice1 = append(slice1, 1)
	slice1 = append(slice1, 2)
	slice1 = append(slice1, 3)
	fmt.Println("slice1:",slice1)
	slice2 = append(slice2, 7)
	slice2 = append(slice2, 8)
	slice2 = append(slice2, 9)
	fmt.Println("slice2:",slice2)
	//copy(slice2, slice1[2:])
	slice2 = slice1
	fmt.Println("copy slice1 to --> slice2")
	fmt.Println("")
	fmt.Println("slice1:",slice1)
	fmt.Println("slice2:",slice2)

	fmt.Println("nil == nil")
	fmt.Println(nil == nil )
}
