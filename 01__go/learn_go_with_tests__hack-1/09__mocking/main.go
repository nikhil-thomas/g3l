package main

import (
	"os"
	"time"
)

type defaultSleeper struct {
}

func (ds defaultSleeper) Sleep() {
	time.Sleep(1 * time.Second)
}

func main() {
	Countdown(os.Stdout, defaultSleeper{})
	Countdown(os.Stdout, &ConfigurableSleeper{
		duration: 3 * time.Second,
		sleep:    time.Sleep,
	})
}
