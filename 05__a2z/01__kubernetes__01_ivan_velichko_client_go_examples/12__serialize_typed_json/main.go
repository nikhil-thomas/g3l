package main

import (
    "encoding/json"
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

    // Typed -> JOSN (option 1)
    // - Serializer = Decoder + Encoder. Since we need only Encoder functionality
    //   in this example, we can pass nil instead of MetaFactory, Creator, and
    //   Typer arguments as they are used only by the Decoder
    encoder := jsonserializer.NewSerializerWithOptions(
        nil,
        nil,
        nil,
        jsonserializer.SerializerOptions{
            Yaml:   false,
            Pretty: false,
            Strict: false,
        })
    encoded, err := runtime.Encode(encoder, &obj)
    if err != nil {
        panic(err.Error())
    }
    fmt.Println("serializer: option 1: ", string(encoded))

    // Typed -> JSON (Option 2)
    // The implementation of Encoder.Encode() in the case of JSON
    // boils down to calling the stdlib encoding/json.Marshal() with optional
    // pretty-printing and converting JSON to YAML
    encoded2, err := json.Marshal(obj)
    if err != nil {
        panic(err.Error())
    }
    fmt.Println("Serialized: Option 2:", string(encoded2))

    // JSON -> Typed
    // Serializer = Decoder + Encoder
    // - jsonserizlizer.Metafactor is a simple partial JOSN unmarshaller that
    //   looks for APIGroup/Version and Kind attributes in the supplied
    //   piece of JSON and parses them into a schema.GroupVersionKind{} object
    // - runtime.ObjectCreator is used to create an empty typed runtime.Object
    //   eg. Deployment, Pod, and ObjectTyper is used to make sure
    //   the MetaFatory's GroupVersionKind matches the on from the `into` argument

    decoder := jsonserializer.NewSerializerWithOptions(
        jsonserializer.DefaultMetaFactory,
        scheme.Scheme,
        scheme.Scheme,
        jsonserializer.SerializerOptions{
            Yaml:   false,
            Pretty: false,
            Strict: false,
        })

    decoded, err := runtime.Decode(decoder, encoded)
    if err != nil {
        panic(err.Error())
    }
    fmt.Printf("Deserialized %#v\n", decoded)
}
