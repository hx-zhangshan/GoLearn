package main

import (
	"fmt"
)

//并发的应用
func main() {
	ch := make(chan string)
	for i := 0; i < 5; i++ {
		go printhello(i, ch)
	}
	//开启 goroutine
	//time.Sleep(time.Millisecond)
	for {
		fmt.Println(<-ch)
	}
}
func printhello(i int, ch chan string) {
	for {
		ch <- fmt.Sprintf("hello world from goroutine %d\n", i)
	}

}
