package main

import (
	"bufio"
	"context"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/workqueue"
	"os"
	"path/filepath"
	"time"
)

var (
	namespace, ConfigMapResource = "default", schema.GroupVersionResource{
		Group:    "",
		Version:  "v1",
		Resource: "configmaps",
	}
)

func main() {
	client := createClientOrDie()

	// The work queue has the following properties:
	// - Fair: items processed in the order in which they are added.
	// - Stingy: a single item will not be processed multiple times concurrently,
	//   and if an item is added multiple times before it can be processed, it
	//   will only be processed once.
	// - Multitenant: Multiple consumers and producers. In particular, it is allowed for an
	//   item to be requeued while it is being processed
	queue := workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())
	defer queue.ShutDown()

	// The queue is typically populated by one or more informers watching events on Kubernetes resources.
	// An "idiomatic" way to get an informer is via
	// a SharedInformerFactory
	// - A factory is essentially a struct keeping a map (type -> informer).
	// - 5*time.Second is a default resync period (for all informers).
	// - namespace makes the informers watch only the specified namespace.
	// - an extra func allows to tweak other listing options like label- or field- selectors.
	factory := dynamicinformer.NewFilteredDynamicSharedInformerFactory(
		client, 5*time.Second, namespace, func(*metav1.ListOptions) {},
	)
	dynamicInformer := factory.ForResource(ConfigMapResource)

	// Informer watches a resource (ConfigMap in this particular example)
	// and simply pushes object keys to the queue.
	dynamicInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			// key is a string <namespace>/<name> (or just <name> for cluster-wide objects)
			key, err := cache.MetaNamespaceKeyFunc(obj)
			if err != nil {
				fmt.Printf("New eventL ADD %s\n", key)
				queue.Add(key)
			}

		},
		UpdateFunc: func(old, new interface{}) {
			key, err := cache.MetaNamespaceKeyFunc(new)
			if err == nil {
				fmt.Printf("New event: UPDATE %s\n", key)
				queue.Add(key)
			}
		},
		DeleteFunc: func(obj interface{}) {
			// much like cache.MetaNamesapceKeyFunc + some extra check
			key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
			if err == nil {
				fmt.Printf("New event: DELETE %s\n", key)
			}
		},
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start the informer's machinery
	factory.Start(ctx.Done())

	// factory.Start() releases the execution flow without waiting for all the
	// internal machinery to warm up
	for gvr, ok := range factory.WaitForCacheSync(ctx.Done()) {
		if !ok {
			panic(fmt.Sprintf("Failed to sync cache for resource %v", gvr))
		}
	}

	// Consuming the work queue with N=3 parallel worker go routines
	for i := 0; i < 1; i++ {
		// a better way is to use wait.Until() from "k8s.io/apimachinery/pkg/util/wait"
		// for every worker
		fmt.Printf("Starting worker %d\n", i)

		// worker()
		go func(n int) {
			for {
				// check done channel
				select {
				case <-ctx.Done():
					fmt.Printf("Controller's done! Worker %d exiting...\n", n)
					return
				default:
				}

				// Obtain a piece of work.
				fmt.Println("worker:", i, "queue len:", queue.Len())
				key, quit := queue.Get()
				fmt.Println("worker:", i, "queue len:", queue.Len())
				if quit {
					fmt.Printf("Work queue has been shut down! Worker %d exiting...\n", n)
					return
				}
				fmt.Printf("Worker %d is about to start process new item %s.\n", n, key)

				// processSingleItem() - scoped to utilize defer and premature returns.
				func() {
					// Tell the queue that we are done with processing this key.
					// This unblocks the key for other workers and allows safe parallel
					// processing because two objects with the same key are never processed
					// in parallel.
					defer queue.Done(key)

					// Your controller's business logic goes here
					obj, err := dynamicInformer.Lister().Get(key.(string))
					if err == nil {
						fmt.Printf("Worker %d found ConfigMap object in informer's cache %#v.\n", n, obj.(*unstructured.Unstructured).GetName())
						// RECONCILE THE OBJECT - PUT YOU BUSINESS LOGIC HERE
						if n == 1 {
							err = fmt.Errorf("worker %d is a chronic failure", n)
						}
					} else {
						fmt.Printf("Worker %d got error %v while looking up ConfigMap object in informer's cache.\n", n, err)
					}

					if err == nil {
						// The key has been handled successfully - forget about it. In partivular, it
						// ensures that future processing of updates for this key won't be rate limited
						// because of errors on previous attempts.
						fmt.Printf("Worker %d reconciled ConfigMap %s successfully, Removing it from the queue.\n", n, key)
						queue.Forget(key)
						return
					}
					// we try no morethan K=5 times
					if queue.NumRequeues(key) >= 5 {
						fmt.Printf("Worker %d gave up on processing %s. Removing it from the queue.\n", n, key)
						queue.Forget(key)
						return
					}

					// Re-enqueue the key rate to be (re-)processed later again
					// Notice that deferred queue.Done(key) call above knows how to deal with re-enqueuing
					// it marks the key as done and then re-appends it again
					fmt.Printf("Worker %d failed to process %s. Putting it back to the queue to retry later.\n", n, key)
					queue.AddRateLimited(key)
					fmt.Println("------------")
				}()
			}
		}(i)
	}
	// Create some kubernetes objects to make the above program actually process something
	cm1 := createConfigMap(client)
	cm2 := createConfigMap(client)
	cm3 := createConfigMap(client)
	cm4 := createConfigMap(client)
	cm5 := createConfigMap(client)

	deleteConfigMap(client, cm1)
	deleteConfigMap(client, cm2)
	deleteConfigMap(client, cm3)
	deleteConfigMap(client, cm4)
	deleteConfigMap(client, cm5)

	// Stay a couple more seconds to let the program finish
	//prompt()
	time.Sleep(300 * time.Second)
	fmt.Println("timer ran our 30")
	queue.ShutDown()
	cancel()
	time.Sleep(1 * time.Second)

}

func prompt() {
	fmt.Printf("Press Return key to continuee...")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		break
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	return
}

func deleteConfigMap(client dynamic.Interface, cm *unstructured.Unstructured) {
	err := client.Resource(ConfigMapResource).
		Namespace(namespace).Delete(context.Background(), cm.GetName(), metav1.DeleteOptions{})
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("Deleted ConfigMap %s/%s\n", cm.GetNamespace(), cm.GetName())
}

func createConfigMap(client dynamic.Interface) *unstructured.Unstructured {
	cm := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "ConfigMap",
			"metadata": map[string]interface{}{
				"namespace":    namespace,
				"generateName": "workqueue-",
			},
			"data": map[string]interface{}{
				"foo": "bar",
			},
		},
	}
	cm, err := client.Resource(ConfigMapResource).Namespace(namespace).Create(context.Background(), cm, metav1.CreateOptions{})
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("Created ConfigMap %s/%s\n", cm.GetNamespace(), cm.GetName())
	return cm
}

func createClientOrDie() dynamic.Interface {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	config, err := clientcmd.BuildConfigFromFlags("", filepath.Join(home, ".kube", "config"))
	if err != nil {
		panic(err.Error())
	}
	client, err := dynamic.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	return client
}
