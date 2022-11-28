package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"k8s-resource.com/m/api"
	client "k8s-resource.com/m/clientset"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"

	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

var kubeconfig string

func init() {
	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	flag.StringVar(&kubeconfig, "kubeconfig", path+"/vke.yaml", "path to Kubernetes config file")
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

	for _, k := range projects.Items {

		fmt.Println(k.Name)

	}

}
