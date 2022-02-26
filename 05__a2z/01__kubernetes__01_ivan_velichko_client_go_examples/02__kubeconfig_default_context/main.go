package main

import (
    "fmt"
    "k8s.io/client-go/discovery"
    _ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
    "k8s.io/client-go/tools/clientcmd"
    "os"
    "path"
)

func main() {
    home, err := os.UserHomeDir()
    if err != nil {
        panic(err)
    }
    config, err := clientcmd.BuildConfigFromFlags("", path.Join(home, ".kube", "config"))
    if err != nil {
        panic(err)
    }
    client, err := discovery.NewDiscoveryClientForConfig(config)
    if err != nil {
        panic(err)
    }
    ver, err := client.ServerVersion()
    if err != nil {
        panic(err)
    }
    fmt.Println(ver.String())
}
