package main

import (
	"fmt"
	"k8s.io/client-go/util/workqueue"
	"time"
)

// This function will simple fill the queue with five
// words and then exit
func fillQueue(queue workqueue.Interface) {
	time.Sleep(time.Second)
	queue.Add("word-1")
	fmt.Println("add word-1, queueLength:", queue.Len())
	queue.Add("word-2")
	fmt.Println("add word-2, queueLength:", queue.Len())
	queue.Add("word-3")
	fmt.Println("add word-3, queueLength:", queue.Len())
	queue.Add("word-4")
	fmt.Println("add word-4, queueLength:", queue.Len())
	queue.Add("word-5")
	fmt.Println("add word-5, queueLength:", queue.Len())
	fmt.Println("fillQueue completed, sending shutdown")
	time.Sleep(5 * time.Second)
	queue.ShutDown()
}

// Read from queue and print results
func readFromQueue(queue workqueue.Interface, stop chan int) {
	time.Sleep(3 * time.Second)
	for {
		item, shutdown := queue.Get()
		fmt.Printf("Got items[shutdown = %t]: %s, - remaining Queue length: %d\n", shutdown, item, queue.Len())
		if shutdown {
			// signal that we are gone
			stop <- -1
			return
		}
		queue.Done(item)
	}
}

func main() {
	fmt.Println("Starting main")
	// Create a channel
	stop := make(chan int)
	// Create a queue
	myQueue := workqueue.New()
	// Create our first worker thread. This goroutine will
	// simply put five words into the queue after one second
	// has passed
	go fillQueue(myQueue)
	// Create second thread that will read from the queue
	go readFromQueue(myQueue, stop)
	fmt.Println("Goroutines started, now waiting for reader to complete")
	<-stop
	fmt.Println("Reader signaled completion, exiting")
}
