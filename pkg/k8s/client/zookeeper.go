package client

import (
	"github.com/pravega/zookeeper-operator/pkg/apis/zookeeper/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

type ZookeeperClusterInterface interface {
	List(opts metav1.ListOptions) (*v1beta1.ZookeeperClusterList, error)
	Create(*v1beta1.ZookeeperCluster) (*v1beta1.ZookeeperCluster, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	Delete(name string, opts metav1.ListOptions) error
}

type zookeeperClusterClient struct {
	restClient rest.Interface
	ns         string
}

func (c *zookeeperClusterClient) List(opts metav1.ListOptions) (*v1beta1.ZookeeperClusterList, error) {
	result := v1beta1.ZookeeperClusterList{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource("zookeeperclusters").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(&result)

	return &result, err
}

func (c *zookeeperClusterClient) Create(cluster *v1beta1.ZookeeperCluster) (*v1beta1.ZookeeperCluster, error) {
	result := v1beta1.ZookeeperCluster{}
	err := c.restClient.
		Post().
		Namespace(c.ns).
		Resource("zookeeperclusters").
		Body(cluster).
		Do().
		Into(&result)
	return &result, err
}

func (c *zookeeperClusterClient) Watch(opt metav1.ListOptions) (watch.Interface, error) {
	opt.Watch = true
	return c.restClient.
		Get().
		Namespace(c.ns).
		Resource("zookeeperclusters").
		VersionedParams(&opt, scheme.ParameterCodec).
		Watch()
}

func (c *zookeeperClusterClient) Delete(name string, opts metav1.ListOptions) error {
	return c.restClient.
		Delete().
		Namespace(c.ns).
		Resource("zookeeperclusters").
		VersionedParams(&opts, scheme.ParameterCodec).
		Name(name).
		Do().
		Error()
}
