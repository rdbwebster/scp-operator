package v1alpha1

import (
	"context"
    "fmt"
	v1 "github.com/rdbwebster/scp-operator/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

type ManagedOperatorInterface interface {
	List(opts metav1.ListOptions) (*v1.ManagedOperatorList, error)
	Get(name string, options metav1.GetOptions) (*v1.ManagedOperator, error)
	Create(*v1.ManagedOperator) (*v1.ManagedOperator, error)
	Delete(name string) (*v1.ManagedOperator, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	// ...
}

type ManagedOperatorClient struct {
	restClient rest.Interface
	ns         string
}

func (c *ManagedOperatorClient) List(opts metav1.ListOptions) (*v1.ManagedOperatorList, error) {
	result := v1.ManagedOperatorList{}

	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource("managedoperators").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(context.TODO()).
		Into(&result)
	fmt.Printf("Here %+v", result)
	return &result, err
}

func (c *ManagedOperatorClient) Get(name string, opts metav1.GetOptions) (*v1.ManagedOperator, error) {
	result := v1.ManagedOperator{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource("managedoperator").
		Name(name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(context.TODO()).
		Into(&result)

	return &result, err
}

func (c *ManagedOperatorClient) Create(oper *v1.ManagedOperator) (*v1.ManagedOperator, error) {
	result := v1.ManagedOperator{}
	err := c.restClient.
		Post().
		Namespace(c.ns).
		Resource("managedoperator").
		Body(oper).
		Do(context.TODO()).
		Into(&result)

	return &result, err
}

func (c *ManagedOperatorClient) Delete(name string) (*v1.ManagedOperator, error) {
	result := v1.ManagedOperator{}
	err := c.restClient.
		Delete().
		Namespace(c.ns).
		Resource("managedoperator").
		Name(name).
		Do(context.TODO()).
		Into(&result)

	return &result, err
}

func (c *ManagedOperatorClient) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.restClient.
		Get().
		Namespace(c.ns).
		Resource("managedoperator").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch(context.TODO())
}
