package main

import (
    "bufio"
    "flag"
    "fmt"
    "k8s.io/client-go/discovery"
    "k8s.io/client-go/tools/clientcmd"
    "k8s.io/client-go/util/homedir"
    "os"
    "path/filepath"
)

func main() {
    kubeconfigPath := getKubeconfigPath()

    kubeconfig, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
    if err != nil {
        panic(err.Error())
    }
    //clientSet, err := kubernetes.NewForConfig(kubeconfig)
    //if err != nil {
    //    panic(err.Error())
    //}
    discoveryClient, err := discovery.NewDiscoveryClientForConfig(kubeconfig)
    groups, resourceList, err := discoveryClient.ServerGroupsAndResources()
    fmt.Println("Groups supported by server")
    for i, group := range groups {
        fmt.Printf("%d. %s,%s\n", i+1, group.Kind, group.Name)
        for j, version := range group.Versions {
            fmt.Printf("\t%d. %s\n", j+1, version.GroupVersion)
        }
        //prompt()
    }
    fmt.Println("Resources supported by server")
    for i, resource := range resourceList {
        fmt.Printf("%d. %s,%s\n", i+1, (*resource).Kind, (*resource).GroupVersion)
        for j, r := range resource.APIResources {
            fmt.Printf("\t%d. name: %s, namespaced: %t, kind: %s, group: %s\n", j+1, r.Name, r.Namespaced, r.Kind, resource.GroupVersion)
        }
        prompt()
    }

    discoveryClient.ServerPreferredResources()
}

func prompt() {
    scanner := bufio.NewScanner(os.Stdin)
    fmt.Printf("press return to continue: ")
    for scanner.Scan() {
        break
    }
    if err := scanner.Err(); err != nil {
        panic(err.Error())
    }
}

func getKubeconfigPath() string {
    defaultKubeconfigPath := ""
    defaultMessage := "absolute path to kubecofig"
    if homeDir := homedir.HomeDir(); homeDir != "" {
        defaultKubeconfigPath = filepath.Join(homeDir, ".kube", "config")
        defaultMessage = "(optional) " + defaultMessage
    }
    kubeconfigPath := flag.String("kubeconfig", defaultKubeconfigPath, defaultMessage)
    flag.Parse()
    return *kubeconfigPath
}
