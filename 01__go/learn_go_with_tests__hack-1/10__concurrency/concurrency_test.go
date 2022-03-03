package concurrency

import (
    "reflect"
    "testing"
    "time"
)

func mockWebsiteChecker(url string) bool {
    if url == "waat://furhurterwe.geds" {
        return false
    }
    return true
}

func slowStubWebsiteChecker(_ string) bool {
    time.Sleep(20 * time.Millisecond)
    return true
}

func TestCheckWebsites(t *testing.T) {
    websites := []string{
        "http://google.com",
        "http://blog.gypsydave5.com",
        "waat://furhurterwe.geds",
    }

    want := map[string]bool{
        "http://google.com":          true,
        "http://blog.gypsydave5.com": true,
        "waat://furhurterwe.geds":    false,
    }
    got := CheckWebsites(mockWebsiteChecker, websites)
    if !reflect.DeepEqual(want, got) {
        t.Errorf("want %v, got %v\n", want, got)
    }
}

func BenchmarkCheckWebsites(b *testing.B) {
    urls := make([]string, 100)
    for i := 0; i < len(urls); i++ {
        urls[i] = "a url"
    }
    for i := 0; i < b.N; i++ {
        CheckWebsites(slowStubWebsiteChecker, urls)

    }
}

func TestConcurrentCheckWebsites(t *testing.T) {
    websites := []string{
        "http://google.com",
        "http://blog.gypsydave5.com",
        "waat://furhurterwe.geds",
    }

    want := map[string]bool{
        "http://google.com":          true,
        "http://blog.gypsydave5.com": true,
        "waat://furhurterwe.geds":    false,
    }
    got := ConcurrentCheckWebsites(mockWebsiteChecker, websites)
    if !reflect.DeepEqual(want, got) {
        t.Errorf("want %v, got %v\n", want, got)
    }

}

func BenchmarkConcurrentCheckWebsites(b *testing.B) {
    urls := make([]string, 100)
    for i := 0; i < len(urls); i++ {
        urls[i] = "a url"
    }
    for i := 0; i < b.N; i++ {
        ConcurrentCheckWebsites(slowStubWebsiteChecker, urls)
    }
}
