package v1alpha1

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"

	"github.com/tkircsi/cluster-crd/api/types/v1alpha1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

type ExtensionV1Alpha1Interface interface {
	Clusters(namespace string) ClusterInterface
}

type ExtensionV1Alpha1Client struct {
	restClient rest.Interface
}

func NewForConfig(c *rest.Config) (*ExtensionV1Alpha1Client, error) {
	crdConfig := *c
	crdConfig.ContentConfig.GroupVersion = &schema.GroupVersion{Group: v1alpha1.GroupName, Version: v1alpha1.GroupVersion}
	crdConfig.APIPath = "/apis"
	crdConfig.NegotiatedSerializer = serializer.NewCodecFactory(scheme.Scheme)
	crdConfig.UserAgent = rest.DefaultKubernetesUserAgent()

	client, err := rest.RESTClientFor(&crdConfig)
	if err != nil {
		return nil, err
	}
	return &ExtensionV1Alpha1Client{restClient: client}, nil
}

func (c *ExtensionV1Alpha1Client) Clusters(namespace string) ClusterInterface {
	return &clusterClient{
		restClient: c.restClient,
		ns:         namespace,
	}
}
