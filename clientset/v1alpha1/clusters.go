package v1alpha1

import (
	"context"

	"github.com/tkircsi/cluster-crd/api/types/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

type ClusterInterface interface {
	List(opts metav1.ListOptions) (*v1alpha1.ClusterList, error)
	Get(name string, options metav1.GetOptions) (*v1alpha1.Cluster, error)
	Create(*v1alpha1.Cluster) (*v1alpha1.Cluster, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
}

type clusterClient struct {
	restClient rest.Interface
	ns         string
}

func (c *clusterClient) List(opts metav1.ListOptions) (*v1alpha1.ClusterList, error) {
	result := v1alpha1.ClusterList{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource("clusters").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(context.Background()).
		Into(&result)

	return &result, err
}

func (c *clusterClient) Get(name string, opts metav1.GetOptions) (*v1alpha1.Cluster, error) {
	result := v1alpha1.Cluster{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource("clusters").
		Name(name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(context.Background()).
		Into(&result)

	return &result, err
}

func (c *clusterClient) Create(project *v1alpha1.Cluster) (*v1alpha1.Cluster, error) {
	result := v1alpha1.Cluster{}
	err := c.restClient.
		Post().
		Namespace(c.ns).
		Resource("clusters").
		Body(project).
		Do(context.Background()).
		Into(&result)

	return &result, err
}

func (c *clusterClient) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.restClient.
		Get().
		Namespace(c.ns).
		Resource("clusters").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch(context.Background())
}
