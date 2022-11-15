package clientset

import (
	"context"

	"k8s-resource.com/m/api"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

type ExampleInterface interface {
	Databases(namespace string, ctx context.Context) DatabaseInterface
}

type ExampleClient struct {
	restClient rest.Interface
}

func NewForConfig(c *rest.Config) (*ExampleClient, error) {
	config := *c
	config.ContentConfig.GroupVersion = &schema.GroupVersion{Group: api.GroupName, Version: api.GroupVersion}
	config.APIPath = "/apis"
	config.NegotiatedSerializer = scheme.Codecs.WithoutConversion()
	config.UserAgent = rest.DefaultKubernetesUserAgent()

	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}

	return &ExampleClient{restClient: client}, nil
}

func (c *ExampleClient) Databases(namespace string, ctx context.Context) DatabaseInterface {
	return &databaseClient{
		restClient: c.restClient,
		ns:         namespace,
		ctx:        ctx,
	}
}
