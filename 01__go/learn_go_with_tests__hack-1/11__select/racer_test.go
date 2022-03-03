package main

import (
    "net/http"
    "net/http/httptest"
    "testing"
    "time"
)

func XTestRacer(t *testing.T) {
    slowURL := "http://www.facebook.com"
    fastURL := "http://www.quii.dev"

    want := fastURL
    got := Racer(slowURL, fastURL)
    if got != want {
        t.Errorf("got %q, want %q", got, want)
    }
}

func TestRacer(t *testing.T) {
    slowServer := makeDelayedServer(20 * time.Millisecond)

    fastServer := makeDelayedServer(0 * time.Millisecond)

    defer slowServer.Close()
    defer fastServer.Close()

    slowUrl := slowServer.URL
    fastUrl := fastServer.URL

    want := fastUrl
    got := Racer(fastUrl, slowUrl)

    if got != want {
        t.Errorf("got %q, want %q", got, want)
    }
}

func TestRacerWithSelect(t *testing.T) {
    slowServer := makeDelayedServer(20 * time.Millisecond)

    fastServer := makeDelayedServer(0 * time.Millisecond)

    defer slowServer.Close()
    defer fastServer.Close()

    slowUrl := slowServer.URL
    fastUrl := fastServer.URL
    timeout := 10 * time.Millisecond

    want := fastUrl
    got, err := RacerWithSelect(fastUrl, slowUrl, timeout)
    if err != nil {
        t.Errorf("expected no error: %v", err)
    }
    if got != want {
        t.Errorf("got %q, want %q", got, want)
    }

}

func TestRacerWithSelectTimeout(t *testing.T) {
    serverA := makeDelayedServer(11 * time.Millisecond)
    serverB := makeDelayedServer(12 * time.Millisecond)
    defer serverB.Close()
    defer serverA.Close()
    timeout := 1 * time.Millisecond

    _, err := RacerWithSelect(serverA.URL, serverB.URL, timeout)
    if err == nil {
        t.Errorf("expected an error")
    }
}

func makeDelayedServer(delay time.Duration) *httptest.Server {
    return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        time.Sleep(delay)
        w.WriteHeader(http.StatusOK)
    }))
}
