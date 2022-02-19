package main

import (
    "fmt"
)

const (
    spanish            = "Spanish"
    french             = "French"
    frenchHelloPrefix  = "Bonjour, "
    englishHelloPrefix = "Hello, "
    spanishHelloPrefix = "Hola, "
)

func Hello(name, language string) string {

    if name == "" {
        name = "World"
    }
    helloPrefix := ""
    helloPrefix = greetingPrefix(language)
    return fmt.Sprintf("%s%s", helloPrefix, name)
}

func greetingPrefix(language string) string {
    prefix := ""
    switch language {

    case spanish:
        prefix = spanishHelloPrefix
    case french:
        prefix = frenchHelloPrefix
    default:
        prefix = englishHelloPrefix
    }
    return prefix
}

func main() {
    fmt.Println(Hello("world", ""))
}
