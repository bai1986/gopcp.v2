package main

import (
	"fmt"
	"io"
	"os"
	"time"
)

func main() {
	fileBasedPipe()
	inMemorySyncPipe()
}

func fileBasedPipe() {
	//系统级别管道,命名管道
	reader, writer, err := os.Pipe()
	if err != nil {
		fmt.Printf("Error: Couldn't create the named pipe: %s\n", err)
	}
	//命名管道操作是阻塞的
	//下面对同一管道的输入和输出操作,会永久阻塞
	// n,err := writer.Write(input)
	// m,err := reader.Read(output)
	go func() {
		output := make([]byte, 100)
		n, err := reader.Read(output)
		if err != nil {
			fmt.Printf("Error: Couldn't read data from the named pipe: %s\n", err)
		}
		fmt.Printf("Read %d byte(s). [file-based pipe]\n", n)
	}()
	input := make([]byte, 26)
	for i := 65; i <= 90; i++ {
		input[i-65] = byte(i)
	}
	//管道可以多路复用
	//操作系统级别管道不保证并发安全
	n, err := writer.Write(input)
	if err != nil {
		fmt.Printf("Error: Couldn't write data to the named pipe: %s\n", err)
	}
	fmt.Printf("Written %d byte(s). [file-based pipe]\n", n)
	//等待管道的读取操作完成
	time.Sleep(200 * time.Millisecond)
}

func inMemorySyncPipe() {
	//基于内存的命名管道
	reader, writer := io.Pipe()
	//命名管道默认阻塞
	go func() {
		output := make([]byte, 100)
		n, err := reader.Read(output)
		if err != nil {
			fmt.Printf("Error: Couldn't read data from the named pipe: %s\n", err)
		}
		fmt.Printf("Read %d byte(s). [in-memory pipe]\n", n)
	}()
	input := make([]byte, 26)
	for i := 65; i <= 90; i++ {
		input[i-65] = byte(i)
	}
	//内存管道支持并发安全性
	n, err := writer.Write(input)
	if err != nil {
		fmt.Printf("Error: Couldn't write data to the named pipe: %s\n", err)
	}
	fmt.Printf("Written %d byte(s). [in-memory pipe]\n", n)
	//为了等待对管道的读取操作完成
	time.Sleep(200 * time.Millisecond)
}
