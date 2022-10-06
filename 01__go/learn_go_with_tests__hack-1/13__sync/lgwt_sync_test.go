package lgwt_sync

import (
    "sync"
    "testing"
)

func TestCounter(t *testing.T) {
    counter := Counter{}
    counter.Inc()
    counter.Inc()
    counter.Inc()

    assertCount(t, &counter, 3)

    t.Run("it runs safely concurrently", func(t *testing.T) {
        wantedCount := 1000
        counter := Counter{}
        var wg sync.WaitGroup
        wg.Add(wantedCount)

        for i := 0; i < wantedCount; i++ {
            go func() {
                counter.Inc()
                wg.Done()
            }()
        }
        wg.Wait()
        assertCount(t, &counter, wantedCount)
    })

}

func assertCount(t testing.TB, got *Counter, want int) {
    t.Helper()
    if got.Value() != want {
        t.Errorf("got %d, want %d", got.Value(), want)
    }
}
