package dependency_injection

import (
	"bytes"
	"testing"
)

func TestGreet(t *testing.T) {
	want := "Hello World"
	buffer := bytes.Buffer{}
	Greet(&buffer)
	got := buffer.String()
	if want != got {
		t.Errorf("want %q, got %v", want, got)
	}
}
