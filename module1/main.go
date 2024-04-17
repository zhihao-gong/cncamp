package main

import (
	"fmt"
	"time"
)

func consumer(messages <-chan int, done <-chan bool) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for _ = range ticker.C {
		select {
		case <-done:
			fmt.Println("child process interrupt...")
			return
		default:
			fmt.Printf("receive message: %d\n", <-messages)
		}
	}
}

func producer(messages chan<- int, done chan bool) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for _ = range ticker.C {
		select {
		case <-done:
			fmt.Println("parent process interrupt...")
			return
		default:
			fmt.Printf("send message: %d\n", 1)
			messages <- 1
		}
	}
}

func main() {
	messages := make(chan int, 10)
	done := make(chan bool)

	defer close(messages)

	go consumer(messages, done)
	go producer(messages, done)

	time.Sleep(10 * time.Second)

	fmt.Println("main process exit!")
	close(done)
}
