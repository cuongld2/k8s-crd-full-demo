package main

import (
	"context"
	"time"

	"k8s-resource.com/m/api"
	client "k8s-resource.com/m/clientset"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
)

func WatchResources(clientSet client.ExampleInterface, context context.Context) cache.Store {
	projectStore, projectController := cache.NewInformer(
		&cache.ListWatch{
			ListFunc: func(lo metav1.ListOptions) (result runtime.Object, err error) {
				return clientSet.Databases(context).List(lo)
			},
			WatchFunc: func(lo metav1.ListOptions) (watch.Interface, error) {
				return clientSet.Databases(context).Watch(lo)
			},
		},
		&api.Database{},
		1*time.Minute,
		cache.ResourceEventHandlerFuncs{},
	)

	go projectController.Run(wait.NeverStop)
	return projectStore
}
