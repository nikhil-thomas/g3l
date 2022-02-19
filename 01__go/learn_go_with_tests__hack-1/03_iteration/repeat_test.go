package iteration_test

import (
    "fmt"
    "github.com/nikhil-thomas/lgwt/iteration"
    "testing"
)

func TestRepeeat(t *testing.T) {
    repeated := iteration.Repeat("a", 10)
    expected := "aaaaaaaaaa"
    if repeated != expected {
        t.Errorf("expected %q but got %q", expected, repeated)
    }
}

func BenchmarkRepeat(b *testing.B) {
    for i := 0; i < b.N; i++ {
        iteration.Repeat("a", 10)
    }
}

func ExampleRepeat() {
    repeated := iteration.Repeat("a", 3)
    fmt.Println(repeated)
    // Output: aaa
}
