package client

import (
	kafkatopicv1alpha1 "github.com/banzaicloud/kafka-operator/api/v1alpha1"
	kafkav1beta1 "github.com/banzaicloud/kafka-operator/api/v1beta1"
	mysqlv1alpha1 "github.com/oracle/mysql-operator/pkg/apis/mysql/v1alpha1"
	zookeeperv1beta1 "github.com/pravega/zookeeper-operator/pkg/apis/zookeeper/v1beta1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

var (
	ClusterConfig      *rest.Config
	MysqlClientset     *ClusterClient
	K8sClient          *kubernetes.Clientset
	ZookeeperClientset *ClusterClient
	KafkaClient        *ClusterClient
	KafkaTopicClient   *ClusterClient
)

func init() {
	ClusterConfig = LoadK8sConfig()
	MysqlClientset, _ = NewForConfig(ClusterConfig, mysqlv1alpha1.GroupName, "v1alpha1")
	ZookeeperClientset, _ = NewForConfig(ClusterConfig, "zookeeper.pravega.io", "v1beta1")
	K8sClient, _ = kubernetes.NewForConfig(ClusterConfig)
	KafkaClient, _ = NewForConfig(ClusterConfig, "kafka.banzaicloud.io", "v1beta1")
	KafkaTopicClient, _ = NewForConfig(ClusterConfig, "kafka.banzaicloud.io", "v1alpha1")
}

type BaseClientInterface interface {
	MysqlClusters(namespace string) MysqlClusterInterface
	ZookeeperCluster(namespace string) ZookeeperClusterInterface
	KafkaCluster(namespace string) KafkaClusterInterface
	KafkaTopic(namespace string) KafkaTopicInterface
}

type ClusterClient struct {
	restClient rest.Interface
}

func NewForConfig(c *rest.Config, groupName string, version string) (*ClusterClient, error) {
	_ = mysqlv1alpha1.AddToScheme(scheme.Scheme)
	_ = zookeeperv1beta1.SchemeBuilder.AddToScheme(scheme.Scheme)
	_ = kafkav1beta1.AddToScheme(scheme.Scheme)
	_ = kafkatopicv1alpha1.AddToScheme(scheme.Scheme)
	config := *c
	config.ContentConfig.GroupVersion = &schema.GroupVersion{Group: groupName, Version: version}
	config.APIPath = "/apis"
	config.NegotiatedSerializer = serializer.DirectCodecFactory{scheme.Codecs}
	config.UserAgent = rest.DefaultKubernetesUserAgent()

	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}

	return &ClusterClient{restClient: client}, nil
}

func (c *ClusterClient) MysqlClusters(namespace string) MysqlClusterInterface {
	return &mysqlClusterClient{
		restClient: c.restClient,
		ns:         namespace,
	}
}

func (c *ClusterClient) ZookeeperCluster(namespace string) ZookeeperClusterInterface {
	return &zookeeperClusterClient{
		restClient: c.restClient,
		ns:         namespace,
	}
}

func (c *ClusterClient) KafkaCluster(namespace string) KafkaClusterInterface {
	return &kafkaClusterClient{
		restClient: c.restClient,
		ns:         namespace,
	}
}

func (c *ClusterClient) KafkaTopic(namespace string) KafkaTopicInterface {
	return &kafkaTopicClusterClient{
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
