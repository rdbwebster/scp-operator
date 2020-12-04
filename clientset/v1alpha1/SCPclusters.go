package v1alpha1

import (
	"context"

	v1 "github.com/rdbwebster/scp-operator/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

type SCPclusterInterface interface {
	List(opts metav1.ListOptions) (*v1.SCPclusterList, error)
	Get(name string, options metav1.GetOptions) (*v1.SCPcluster, error)
	Create(*v1.SCPcluster) (*v1.SCPcluster, error)
	Delete(name string) (*v1.SCPcluster, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	// ...
}

type ScpclusterClient struct {
	restClient rest.Interface
	ns         string
}

func (c *ScpclusterClient) List(opts metav1.ListOptions) (*v1.SCPclusterList, error) {
	result := v1.SCPclusterList{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource("scpclusters").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(context.TODO()).
		Into(&result)

	return &result, err
}

func (c *ScpclusterClient) Get(name string, opts metav1.GetOptions) (*v1.SCPcluster, error) {
	result := v1.SCPcluster{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource("scpcluster").
		Name(name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(context.TODO()).
		Into(&result)

	return &result, err
}

func (c *ScpclusterClient) Create(cluster *v1.SCPcluster) (*v1.SCPcluster, error) {
	result := v1.SCPcluster{}
	err := c.restClient.
		Post().
		Namespace(c.ns).
		Resource("scpclusters").
		Body(cluster).
		Do(context.TODO()).
		Into(&result)

	return &result, err
}

func (c *ScpclusterClient) Delete(name string) (*v1.SCPcluster, error) {
	result := v1.SCPcluster{}
	err := c.restClient.
		Delete().
		Namespace(c.ns).
		Resource("scpclusters").
		Name(name).
		Do(context.TODO()).
		Into(&result)

	return &result, err
}

func (c *ScpclusterClient) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.restClient.
		Get().
		Namespace(c.ns).
		Resource("scpclusters").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch(context.TODO())
}
