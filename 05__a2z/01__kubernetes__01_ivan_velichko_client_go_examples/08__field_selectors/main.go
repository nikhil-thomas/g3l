package main

import (
    "fmt"
    "github.com/pkg/errors"
    "k8s.io/apimachinery/pkg/fields"
)

var (
    ErrNoFieldSetMatch = errors.New("selector should have matched the field set")
)

func main() {
    flds := fields.Set{
        "foo": "bar",
        "baz": "qux",
    }
    // Selector matching existing field seth
    sel := fields.SelectorFromSet(flds)
    if sel.Matches(flds) {
        fmt.Printf("Selector %v matches field %v\n", sel, flds)
    } else {
        panic(ErrNoFieldSetMatch)
    }

}
