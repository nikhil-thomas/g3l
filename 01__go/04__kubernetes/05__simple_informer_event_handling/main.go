package main

import (
	"fmt"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
)

// A generic event handler for pod events
func handleEvent(pod *v1.Pod, eventType string) {
	fmt.Printf("Recorded event of type %s on pod %s/%s\n", eventType, pod.Namespace, pod.Name)
}

// Our event handler for adding a pod
func onAdd(obj interface{}) {
	pod, ok := obj.(*v1.Pod)
	if !ok {
		fmt.Printf("Error during conversion, this does not seen tom be a pod\n")
		return
	}
	handleEvent(pod, "ADD")
}

// Our event handler for deletion of a pod
func onDelete(obj interface{}) {
	pod, ok := obj.(*v1.Pod)
	if !ok {
		fmt.Printf("Error during conversion, this does not seem to be a pod\n")
		return
	}
	handleEvent(pod, "DEL")
}

// our event handler for modification of a pod
func onUpdate(old interface{}, new interface{}) {
	pod, ok := old.(*v1.Pod)
	if !ok {
		fmt.Printf("Error during conversion, this does not seem to be a pod\n")
		return
	}
	handleEvent(pod, "MOD")
}

// Create channel that will be closed when a signal is received
func createSignalHandler() (stopCh <-chan struct{}) {
	stop := make(chan struct{})
	c := make(chan os.Signal, 2)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-c
		fmt.Printf("Signal handler: reveived signal %s\n", sig.String())
		close(stop)
	}()
	return stop
}

func main() {
	// create a clientSet
	home := homedir.HomeDir()
	kubeconfig, err := clientcmd.BuildConfigFromFlags("", filepath.Join(home, ".kube", "config"))
	if err != nil {
		panic(err.Error())
	}
	clientset, err := kubernetes.NewForConfig(kubeconfig)
	if err != nil {
		panic(err.Error())
	}
	// implement signal handler - this will return a channel which will be closed if a signal is received
	stopCh := createSignalHandler()

	// to create informer, we will use a factory. This factory expects two argument
	// a clientset and a resync time (after this time,
	// the cache will be rebuilt from scratch
	factory := informers.NewSharedInformerFactory(clientset, 30*time.Second)

	// we can ask our factory to create a pod informer

	podInformer := factory.Core().V1().Pods().Informer()
	fmt.Println("Starting informer")

	// starting a factory will start all informers created by the factory
	factory.Start(stopCh)
	fmt.Println("informer running")

	// wait for informer to sync. 

	if ok := cache.WaitForCacheSync(stopCh, podInformer.HasSynced); !ok {
		panic("error while waiting for informer to sync")
	}

	// the informer main loop is now running. We can now add our event handlers
	// and wait form the stopCh to be closed
	fmt.Println("Informer synced, now adding event handlers and waiting for stop channel")
	podInformer.AddEventHandler(
		&cache.ResourceEventHandlerFuncs{
			AddFunc:    onAdd,
			UpdateFunc: onUpdate,
			DeleteFunc: onDelete,
		})
	<-stopCh
}
