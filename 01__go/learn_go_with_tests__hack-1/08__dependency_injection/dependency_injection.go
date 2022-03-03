package dependency_injection

import (
	"fmt"
	"io"
)

func Greet(out io.Writer) {
	fmt.Fprintf(out, "Hello World")
}
