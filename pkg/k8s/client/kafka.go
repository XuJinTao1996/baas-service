package client

import (
	kafkav1beta1 "github.com/banzaicloud/kafka-operator/api/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

type KafkaClusterInterface interface {
	List(opts metav1.ListOptions) (*kafkav1beta1.KafkaClusterList, error)
	Create(cluster *kafkav1beta1.KafkaCluster) (*kafkav1beta1.KafkaCluster, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	Delete(name string, opts metav1.ListOptions) error
}

type kafkaClusterClient struct {
	restClient rest.Interface
	ns         string
}

func (kc *kafkaClusterClient) List(opts metav1.ListOptions) (*kafkav1beta1.KafkaClusterList, error) {
	result := kafkav1beta1.KafkaClusterList{}
	err := kc.restClient.
		Get().
		Namespace(kc.ns).
		Resource("kafkaclusters").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(&result)

	return &result, err
}

func (kc *kafkaClusterClient) Create(cluster *kafkav1beta1.KafkaCluster) (*kafkav1beta1.KafkaCluster, error) {
	result := kafkav1beta1.KafkaCluster{}
	err := kc.restClient.
		Post().
		Namespace(kc.ns).
		Resource("kafkaclusters").
		Body(cluster).
		Do().
		Into(&result)
	return &result, err
}

func (kc *kafkaClusterClient) Watch(opt metav1.ListOptions) (watch.Interface, error) {
	opt.Watch = true
	return kc.restClient.
		Get().
		Namespace(kc.ns).
		Resource("kafkaclusters").
		VersionedParams(&opt, scheme.ParameterCodec).
		Watch()
}

func (kc *kafkaClusterClient) Delete(name string, opts metav1.ListOptions) error {
	return kc.restClient.
		Delete().
		Namespace(kc.ns).
		Resource("kafkaclusters").
		VersionedParams(&opts, scheme.ParameterCodec).
		Name(name).
		Do().
		Error()
}
