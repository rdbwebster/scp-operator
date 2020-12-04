package v1alpha1

import (
	v1 "github.com/rdbwebster/scp-operator/api/v1"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

type ExampleV1Alpha1Interface interface {
	SCPclusters(namespace string) SCPclusterInterface
	ManagedOperators(namespace string) ManagedOperatorInterface
}

type ExampleV1Alpha1Client struct {
	restClient rest.Interface
}

func NewForConfig(c *rest.Config) (*ExampleV1Alpha1Client, error) {

	v1.AddToScheme(scheme.Scheme)

	config := *c
	config.ContentConfig.GroupVersion = &v1.GroupVersion
	config.APIPath = "/apis"
	config.NegotiatedSerializer = serializer.NewCodecFactory(scheme.Scheme)
	//config.NegotiatedSerializer = scheme.Codecs.WithoutConversion()
	config.UserAgent = rest.DefaultKubernetesUserAgent()

	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}

	return &ExampleV1Alpha1Client{restClient: client}, nil
}

func (c *ExampleV1Alpha1Client) SCPcluster(namespace string) SCPclusterInterface {
	return &ScpclusterClient{
		restClient: c.restClient,
		ns:         namespace,
	}
}

func (c *ExampleV1Alpha1Client) ManagedOperator(namespace string) ManagedOperatorInterface {
	return &ManagedOperatorClient{
		restClient: c.restClient,
		ns:         namespace,
	}
}
