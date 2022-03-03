package concurrency

type WebsiteChecker func(string) bool

type result struct {
    string
    bool
}

func CheckWebsites(wc WebsiteChecker, urls []string) map[string]bool {
    results := make(map[string]bool)
    for _, url := range urls {
        results[url] = wc(url)
    }
    return results
}

func ConcurrentCheckWebsites(wc WebsiteChecker, urls []string) map[string]bool {
    results := map[string]bool{}
    resultChannel := make(chan result)
    for _, url := range urls {
        go func(u string) {
            //results[u] = wc(u)
            resultChannel <- result{u, wc(u)}
        }(url)
    }
    //time.Sleep(2 * time.Second)
    for i := 0; i < len(urls); i++ {
        result := <-resultChannel
        results[result.string] = result.bool
    }
    return results
}
