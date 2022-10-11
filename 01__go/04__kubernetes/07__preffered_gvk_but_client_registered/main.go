package main

import (
	"bufio"
	"fmt"
	"golang.stackrox.io/kube-linter/pkg/objectkinds"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"strings"

	"os"
	"path/filepath"
)

func main() {
	fmt.Println("hello")
	homeDir := homedir.HomeDir()
	kubeconfigPath := filepath.Join(homeDir, ".kube", "config")
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		panic(err)
	}
	disCli, err := discovery.NewDiscoveryClientForConfig(config)
	fmt.Println(disCli)
	if err != nil {
		panic(err)
	}
	//expServerPrefferedVsKubeSchemeResources(disCli)
	reconcileResourceList(disCli)
}

func reconcileResourceList(disCli discovery.DiscoveryInterface) ([]metav1.APIResource, error) {
	apiResources := []metav1.APIResource{}

	apiResourceSet := map[string]*metav1.APIResource{}

	disCli.ServerPreferredResources()

	// get all resources

	apiGroups, apiResourceLists, err := disCli.ServerGroupsAndResources()
	if err != nil {
		return nil, err
	}
	rscCount := 0
	klAccept := 0
	schemeAccept := 0
	for _, apiResourceList := range apiResourceLists {
		gv, err := schema.ParseGroupVersion(apiResourceList.GroupVersion)
		if err != nil {
			panic(err)
		}
		for _, apiResource := range apiResourceList.APIResources {
			//fmt.Printf("%d/%d, GroupVersion: %s, G: %s, V: %s, Kind: %s, storageversionhash: %s\n", i+1, j+1, apiResourceList.GroupVersion, apiResource.Group, apiResource.Version, apiResource.Kind, apiResource.StorageVersionHash)
			//fmt.Printf("%d/%d. string: %s\n", i+1, j+1, apiResource.String())
			if isSubResource(&apiResource) {
				//fmt.Println("skip", apiResource.Kind, apiResource.Name)
				continue
			}
			apiResource.Group = gv.Group
			apiResource.Version = gv.Version
			rscCount++
			canValidate, err := isRegisteredKubeLinterKind(apiResource)
			if err != nil {
				return nil, err
			}

			if !canValidate {
				//fmt.Println("kubelinter discard", apiResource.Name, apiResource.Kind)
				continue
			}
			klAccept++

			kubeScheme := scheme.Scheme
			if !isRegisteredInScheme(&apiResource, kubeScheme) {
				fmt.Println("not registered", apiResource.Kind, apiResource.Name, apiResource.Group, apiResource.Version)
				continue
			}
			schemeAccept++
			fmt.Println("accept", apiResource.Kind, apiResource.Name, apiResource.Group, apiResource.Version, rscCount, klAccept, schemeAccept)

			//apiResources = append(apiResources, apiResource)
			apiResourceSet[apiResource.Kind] = &apiResource

		}
	}

	// de-duplicate
	// serverPrefferedResources
	serverPR, err := disCli.ServerPreferredResources()
	if err != nil {
		panic(err)
	}
	for kind, apiResource := range apiResourceSet {
		gv := prefferedGroupVersionForKind(kind)
		for _, r := serverPR {

		}
	}

	for i := range apiResources {
		fmt.Println(apiResources[i].Kind)
		for j := i + 1; j < len(apiResources); j++ {
			if apiResources[i].Kind == apiResources[j].Kind {
				if apiResources[i].Group != apiResources[j].Group || apiResources[i].Version != apiResources[j].Version {

					fmt.Println(apiResources[j].Kind, "duplicate")
				}
			}
		}
	}


	return apiResources, nil
}

func prefferedGroupVersionForKind(kind string, groups metav1.APIGroupList) metav1.GroupVersion {
	for _, group := range groups {

	}
}

func isSubResource(apiResource *metav1.APIResource) bool {
	return strings.Contains(apiResource.Name, "/")
}

func isRegisteredInScheme(apiResource *metav1.APIResource, schemes ...*runtime.Scheme) bool {
	gvk := schema.GroupVersionKind{
		Group:   apiResource.Group,
		Version: apiResource.Version,
		Kind:    apiResource.Kind,
	}
	for _, sch := range schemes {
		_, err := sch.New(gvk)
		if err == nil {
			return true
		}
		if runtime.IsNotRegisteredError(err) {
			continue
		}
		return false

	}
	return false
}

func expServerPrefferedVsKubeSchemeResources(disCli discovery.DiscoveryInterface) {
	disCli.ServerPreferredResources()
	apiGroups, apiResourceList, err := disCli.ServerGroupsAndResources()
	for i, apiGroup := range apiGroups {
		fmt.Printf("%d. GroupName: %s, Preffered Version: %s\n", i+1, apiGroup.Name, apiGroup.PreferredVersion)
	}
	fmt.Println("----------------------")
	for i, rl := range apiResourceList {
		fmt.Printf("%d. groupversion: %s\n", i+1, rl.GroupVersion)
		for j, resource := range rl.APIResources {
			fmt.Printf("%d. Group: %s, Version: %s, Kind: %s\n", j+1, resource.Group, resource.Version, resource.Kind)
		}
		fmt.Println("::::::::::::::::::::")
	}
	fmt.Println("--------------------")
	prompt()
	kubeScheme := scheme.Scheme
	for i, version := range kubeScheme.PreferredVersionAllGroups() {
		fmt.Printf("%d. %s, %s\n", i+1, version.Group, version.Version)

	}
	prompt()
	for i, version := range kubeScheme.PrioritizedVersionsAllGroups() {
		fmt.Printf("%d. %s, %s\n", i+1, version.Group, version.Version)
		prompt()

	}
	prompt()

	gvk := schema.GroupVersionKind{
		Group:   "batch",
		Version: "v1",
		Kind:    "CronJob",
	}
	typedObj, err := kubeScheme.New(gvk)
	fmt.Println(err, runtime.IsNotRegisteredError(err))
	if err == nil {

		fmt.Println(typedObj.GetObjectKind().GroupVersionKind())
	}
}

func prompt() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		break
	}
	err := scanner.Err()
	if err != nil {
		panic(err)
	}
}

func isRegisteredKubeLinterKind(rsrc metav1.APIResource) (bool, error) {
	// Construct the gvks for objects to watch.  Remove the Any
	// kind or else all objects kinds will be watched.
	kubeLinterKinds := getKubeLinterKinds()
	kubeLinterMatcher, err := objectkinds.ConstructMatcher(kubeLinterKinds...)
	if err != nil {
		return false, err
	}

	gvk := gvkFromMetav1APIResource(rsrc)
	if kubeLinterMatcher.Matches(gvk) {
		return true, nil
	}
	return false, nil
}

func getKubeLinterKinds() []string {
	kubeLinterKinds := objectkinds.AllObjectKinds()
	for i := range kubeLinterKinds {
		if kubeLinterKinds[i] == objectkinds.Any {
			kubeLinterKinds = append(kubeLinterKinds[:i], kubeLinterKinds[i+1:]...)
			break
		}
	}
	return kubeLinterKinds
}

func gvkFromMetav1APIResource(rsc metav1.APIResource) schema.GroupVersionKind {
	return schema.GroupVersionKind{
		Group:   rsc.Group,
		Version: rsc.Version,
		Kind:    rsc.Kind,
	}
}
