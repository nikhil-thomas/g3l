package main

import (
    "fmt"
    _ "github.com/tektoncd/catalog/task/tkn/0.1"
    _ "github.com/tektoncd/catalog/task/aws-cli/0.1"
)

func main() {
    fmt.Println("hello")
}
