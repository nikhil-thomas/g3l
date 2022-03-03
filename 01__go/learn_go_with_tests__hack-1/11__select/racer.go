package main

import (
    "fmt"
    "net/http"
    "time"
)

func Racer(a, b string) string {
    aDuration := measureResponseTime(a)

    bDuration := measureResponseTime(b)

    if aDuration < bDuration {
        return a
    }
    return b
}

func RacerWithSelect(a, b string, timeout time.Duration) (string, error) {
    select {
    case <-ping(a):
        return a, nil
    case <-ping(b):
        return b, nil
    case <-time.After(timeout):
        return "", fmt.Errorf("timed out waiting for %s and %s", a, b)
    }
}

func ping(url string) chan struct{} {
    out := make(chan struct{})
    go func() {
        http.Get(url)
        //out <- struct{}{}
        close(out)
    }()

    return out
}

func measureResponseTime(url string) time.Duration {
    start := time.Now()
    http.Get(url)
    return time.Since(start)
}
