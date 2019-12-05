package client

import (
	mysqlv1alpha1 "github.com/oracle/mysql-operator/pkg/apis/mysql/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

type BaseClientInterface interface {
	Clusters(namespace string) ClusterInterface
}

type ClusterClient struct {
	restClient rest.Interface
}

func NewForConfig(c *rest.Config, groupName string, version string) (*ClusterClient, error) {
	_ = mysqlv1alpha1.AddToScheme(scheme.Scheme)
	config := *c
	config.ContentConfig.GroupVersion = &schema.GroupVersion{Group: groupName, Version: version}
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

func LoadK8sConfig() *rest.Config {
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err)
	}
	return config
}
