package main

import (
    "fmt"
    "k8s.io/client-go/tools/clientcmd"
    "os"
    "path"
)

func main() {
    home, err := os.UserHomeDir()
    if err != nil {
        panic(err)
    }
    config, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
        &clientcmd.ClientConfigLoadingRules{ExplicitPath: path.Join(home, ".kube", "config")},
        &clientcmd.ConfigOverrides{},
    ).RawConfig()
    if err != nil {
        panic(err.Error())
    }
    for name := range config.Contexts {
        fmt.Printf("Found context %s\n", name)
    }
}
