package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"k8s-resource.com/m/api"
	client "k8s-resource.com/m/clientset"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"

	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	//
	// Uncomment to load all auth plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth"
	//
	// Or uncomment to load specific auth plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth/azure"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
)

var kubeconfig string

func init() {
	flag.StringVar(&kubeconfig, "kubeconfig", "/home/donald/.kube/vke-2adb84ff-6b5a-4e5e-a92e-0cafa3b13486.yaml", "path to Kubernetes config file")
	flag.Parse()
}

func main() {
	var config *rest.Config
	var err error

	if kubeconfig == "" {
		log.Printf("using in-cluster configuration")
		config, err = rest.InClusterConfig()
	} else {
		log.Printf("using configuration from '%s'", kubeconfig)
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	}

	if err != nil {
		panic(err)
	}

	api.AddToScheme(scheme.Scheme)

	clientSet, err := client.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	context := context.TODO()

	projects, err := clientSet.Databases("default", context).List(metav1.ListOptions{})
	if err != nil {
		panic(err)
	}

	fmt.Printf("projects found: %+v\n", projects)

	store := WatchResources(clientSet, context)

	for {
		projectsFromStore := store.List()
		fmt.Printf("project in store: %d\n", len(projectsFromStore))

		time.Sleep(2 * time.Second)
	}
}

// // +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// type Database struct {
// 	metav1.TypeMeta   `json:",inline"`
// 	metav1.ObjectMeta `json:"metadata,omitempty"`

// 	Spec interface{} `json:"spec"`
// }

// func main() {
// 	ctx := context.Background()
// 	var kubeconfig *string
// 	if home := homedir.HomeDir(); home != "" {
// 		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "vke-2adb84ff-6b5a-4e5e-a92e-0cafa3b13486.yaml"), "(optional) absolute path to the kubeconfig file")
// 	} else {
// 		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
// 	}
// 	flag.Parse()

// 	// use the current context in kubeconfig
// 	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
// 	if err != nil {
// 		panic(err.Error())
// 	}

// 	// create the clientset
// 	clientSet, err := kubernetes.NewForConfig(config)
// 	if err != nil {
// 		panic(err.Error())
// 	}

// 	projects, err := clientSet.Projects("default").List(metav1.ListOptions{})
// 	if err != nil {
// 		panic(err)
// 	}

// 	fmt.Printf("projects found: %+v\n", projects)

// 	store := WatchResources(clientSet)

// 	for {
// 		projectsFromStore := store.List()
// 		fmt.Printf("project in store: %d\n", len(projectsFromStore))

// 		time.Sleep(2 * time.Second)
// 	}

// 	// defaultNamespace, err := clientset.CoreV1().Namespaces().Get(ctx, "default", metav1.GetOptions{})
// 	// if err != nil {
// 	// 	panic(err.Error())
// 	// }
// 	// // deployments := clientset.AppsV1().Deployments("default")
// 	// // fmt.Println(deployments)
// 	// fmt.Println(defaultNamespace.Name)

// 	// watchObject, err := clientset.AppsV1().RESTClient().Get().Resource("/apis/apiextensions.k8s.io/v1/customresourcedefinitions/dbs").Watch(ctx)
// 	// if err != nil {
// 	// 	panic(err.Error())
// 	// }
// 	// fmt.Println(watchObject)
// 	// clientset.Resource(gvr).
// 	// 	Namespace(namespace).Get("foo", metav1.GetOptions{})
// 	// b, err := clientset.RESTClient().Get().AbsPath("/apis/databases.resource.donald.com/v1/db").DoRaw(ctx)
// 	// if err != nil {
// 	// 	panic(err.Error())
// 	// }
// 	// fmt.Println(b)
// }
