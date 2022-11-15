package clientset

import (
	"context"

	"k8s-resource.com/m/api"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

type DatabaseInterface interface {
	List(opts metav1.ListOptions) (*api.DatabaseList, error)
	Get(name string, options metav1.GetOptions) (*api.Database, error)
	Create(*api.Database) (*api.Database, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	// ...
}

type databaseClient struct {
	restClient rest.Interface
	ns         string
	ctx        context.Context
}

func (c *databaseClient) List(opts metav1.ListOptions) (*api.DatabaseList, error) {
	result := api.DatabaseList{}
	err := c.restClient.
		Get().
		AbsPath("/apis/resource.donald.com/v1/databases").
		Do(c.ctx).
		Into(&result)

	return &result, err
}

func (c *databaseClient) Get(name string, opts metav1.GetOptions) (*api.Database, error) {
	result := api.Database{}
	err := c.restClient.
		Get().
		AbsPath("/apis/resource.donald.com/v1/databases").
		Name(name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(c.ctx).
		Into(&result)

	return &result, err
}

func (c *databaseClient) Create(database *api.Database) (*api.Database, error) {
	result := api.Database{}
	err := c.restClient.
		Post().
		AbsPath("/apis/resource.donald.com/v1/databases").
		Body(database).
		Do(c.ctx).
		Into(&result)

	return &result, err
}

func (c *databaseClient) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.restClient.
		Get().
		AbsPath("/apis/resource.donald.com/v1/databases").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch(c.ctx)
}
