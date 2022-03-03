package main

import (
    "fmt"
    "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
    "k8s.io/apimachinery/pkg/runtime"
    "k8s.io/apimachinery/pkg/util/yaml"
)

func main() {
    yConfigMap := `---
apiVersion: v1
kind: ConfigMap
metadata:
  creationTimestamp:
  name: my-configmap
  namespace: default
data:
  foo: bar
`
    // yaml -> Unstructured (through JSON)
    jConfigMap, err := yaml.ToJSON([]byte(yConfigMap))
    fmt.Printf("%s\n", jConfigMap)
    if err != nil {
        panic(err)
    }
    object, err := runtime.Decode(unstructured.UnstructuredJSONScheme, jConfigMap)
    if err != nil {
        panic(err)
    }
    uConfigMap, ok := object.(*unstructured.Unstructured)
    if !ok {
        panic("unstructured.Unstructured expected")
    }
    if uConfigMap.GetName() != "my-configmap" {
        panic("unexpected configmap data")
    }
    fmt.Println("Pass")
}
