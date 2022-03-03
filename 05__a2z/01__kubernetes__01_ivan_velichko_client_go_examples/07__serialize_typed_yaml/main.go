package main

import (
    "fmt"
    v1 "k8s.io/api/core/v1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/apimachinery/pkg/runtime"
    jsonserializer "k8s.io/apimachinery/pkg/runtime/serializer/json"
    "k8s.io/client-go/kubernetes/scheme"
)

func main() {
    obj := v1.ConfigMap{
        TypeMeta: metav1.TypeMeta{
            Kind:       "ConfigMap",
            APIVersion: "v1",
        },
        ObjectMeta: metav1.ObjectMeta{
            Name:      "my-configmap",
            Namespace: "default",
        },
        Data: map[string]string{
            "foo": "bar",
        },
    }

    // Serializer = Encoder + decoder
    serializer := jsonserializer.NewSerializerWithOptions(
        jsonserializer.DefaultMetaFactory,
        scheme.Scheme,
        scheme.Scheme,
        jsonserializer.SerializerOptions{
            Yaml:   true,
            Pretty: false,
            Strict: false,
        })
    // Runtime.Encode() is just a helper function to invoke Encoder.Encode()
    encodedYaml, err := runtime.Encode(serializer, &obj)
    if err != nil {
        panic(err.Error())
    }
    fmt.Println("Serialized:\n%s", string(encodedYaml))

    // YAML -> typed
    decoded, err := runtime.Decode(serializer, encodedYaml)
    if err != nil {
        panic(err.Error())
    }
    fmt.Printf("Deserialized: %#v\n", decoded)
}
