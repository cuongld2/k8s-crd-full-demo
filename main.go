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
	flag.StringVar(&kubeconfig, "kubeconfig", "/home/donaldle/Projects/Personal/vultr/k8s-crd/k8s-crd-full-demo/vke-2adb84ff-6b5a-4e5e-a92e-0cafa3b13486.yaml", "path to Kubernetes config file")
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

	projects, err := clientSet.Databases(context).List(metav1.ListOptions{})
	if err != nil {
		panic(err)
	}

	fmt.Printf("projects found: %+v\n", projects)

	store := WatchResources(clientSet, context)

	newDatabase := new(api.Database) // pa == &Student{"", 0}
	newDatabase.Name = "mongodb"
	newDatabase.Kind = "Database" // pa == &Student{"Alice", 0}
	newDatabase.APIVersion = "resource.donald.com/v1"
	newDatabase.Spec.DbName = "mongodb"
	newDatabase.Spec.Description = "Used storing unstructured data"
	newDatabase.Spec.Total = 100
	newDatabase.Spec.Available = 100
	newDatabase.Spec.DbType = "noSQL"
	newDatabase.Spec.Tags = "Web Development, nosql data"

	projectCreated, err := clientSet.Databases(context).Create(newDatabase)
	if err != nil {
		panic(err)
	}

	fmt.Println(projectCreated)

	projectDeleted, err := clientSet.Databases(context).Delete(newDatabase.Name, metav1.GetOptions{})
	if err != nil {
		panic(err)
	}

	fmt.Println(projectDeleted)

	for {
		projectsFromStore := store.List()
		fmt.Printf("project in store: %d\n", len(projectsFromStore))

		time.Sleep(2 * time.Second)
	}
}
