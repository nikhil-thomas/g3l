package main

import (
    "fmt"
    "k8s.io/apimachinery/pkg/labels"
    "k8s.io/apimachinery/pkg/selection"
)

func main() {
    lbls := labels.Set{"foo": "bar", "baz": "qux", "foo2": "bar2"}
    sel := labels.NewSelector()
    req, err := labels.NewRequirement("foo", selection.Equals, []string{"bar"})
    if err != nil {
        panic(err)
    }
    sel = sel.Add(*req)
    req, _ = labels.NewRequirement("foo2", selection.Equals, []string{"bar2"})
    sel = sel.Add(*req)
    fmt.Printf("\n\n%v\n\n", sel)
    if sel.Matches(lbls) {
        fmt.Printf("Selector %v matched label set %v\n", sel, lbls)
    } else {
        panic("Selector should have matched labels")
    }

    sel, err = labels.Parse("foo=bar")
    if err != nil {
        panic(err)
    }
    if sel.Matches(lbls) {
        fmt.Printf("Selector %v matched label set %v\n", sel, lbls)
    } else {
        panic("Selector should have matched labels")
    }
}
