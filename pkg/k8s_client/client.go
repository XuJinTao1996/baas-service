package k8s_client

import (
	mysqlv1alpha1 "github.com/oracle/mysql-operator/pkg/apis/mysql/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

type ClientBaseInterface interface {
	Clusters(namespace string) ClusterInterface
}

type ClusterClient struct {
	restClient rest.Interface
}

func NewForConfig(c *rest.Config) (*ClusterClient, error) {
	config := *c
	config.ContentConfig.GroupVersion = &schema.GroupVersion{Group: mysqlv1alpha1.GroupName, Version: "v1alpha1"}
	config.APIPath = "/apis"
	config.NegotiatedSerializer = serializer.DirectCodecFactory{CodecFactory: scheme.Codecs}
	config.UserAgent = rest.DefaultKubernetesUserAgent()

	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}

	return &ClusterClient{restClient: client}, nil
}

func (c *ClusterClient) Clusters(namespace string) ClusterInterface {
	return &mysqlClusterClient{
		restClient: c.restClient,
		ns:         namespace,
	}
}
