package client

import (
	mysqlv1alpha1 "github.com/oracle/mysql-operator/pkg/apis/mysql/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

type ClusterInterface interface {
	List(opts metav1.ListOptions) (*mysqlv1alpha1.ClusterList, error)
	Create(cluster *mysqlv1alpha1.Cluster) (*mysqlv1alpha1.Cluster, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	Delete(name string, opts metav1.ListOptions) error
}

type mysqlClusterClient struct {
	restClient rest.Interface
	ns         string
}

func (c *mysqlClusterClient) List(opts metav1.ListOptions) (*mysqlv1alpha1.ClusterList, error) {
	result := mysqlv1alpha1.ClusterList{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource("mysqlclusters").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(&result)

	return &result, err
}

func (c *mysqlClusterClient) Create(cluster *mysqlv1alpha1.Cluster) (*mysqlv1alpha1.Cluster, error) {
	result := mysqlv1alpha1.Cluster{}
	err := c.restClient.
		Post().
		Namespace(c.ns).
		Resource("mysqlclusters").
		Body(cluster).
		Do().
		Into(&result)
	return &result, err
}

func (c *mysqlClusterClient) Watch(opt metav1.ListOptions) (watch.Interface, error) {
	opt.Watch = true
	return c.restClient.
		Get().
		Namespace(c.ns).
		Resource("mysqlclusters").
		VersionedParams(&opt, scheme.ParameterCodec).
		Watch()
}

func (c *mysqlClusterClient) Delete(name string, opts metav1.ListOptions) error {
	return c.restClient.
		Delete().
		Namespace(c.ns).
		Resource("mysqlclusters").
		VersionedParams(&opts, scheme.ParameterCodec).
		Name(name).
		Do().
		Error()
}
