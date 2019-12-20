package client

import (
	kafkatopicv1alpha1 "github.com/banzaicloud/kafka-operator/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

type kafkaTopicClusterClient struct {
	restClient rest.Interface
	ns         string
}

type KafkaTopicInterface interface {
	List(opts metav1.ListOptions) (*kafkatopicv1alpha1.KafkaTopicList, error)
	Create(cluster *kafkatopicv1alpha1.KafkaTopic) (*kafkatopicv1alpha1.KafkaTopic, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	Delete(name string, opts metav1.ListOptions) error
}

func (kc *kafkaTopicClusterClient) List(opts metav1.ListOptions) (*kafkatopicv1alpha1.KafkaTopicList, error) {
	result := kafkatopicv1alpha1.KafkaTopicList{}
	err := kc.restClient.
		Get().
		Namespace(kc.ns).
		Resource("kafkaclusters").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(&result)

	return &result, err
}

func (kc *kafkaTopicClusterClient) Create(cluster *kafkatopicv1alpha1.KafkaTopic) (*kafkatopicv1alpha1.KafkaTopic, error) {
	result := kafkatopicv1alpha1.KafkaTopic{}
	err := kc.restClient.
		Post().
		Namespace(kc.ns).
		Resource("kafkatopics").
		Body(cluster).
		Do().
		Into(&result)
	return &result, err
}

func (kc *kafkaTopicClusterClient) Watch(opt metav1.ListOptions) (watch.Interface, error) {
	opt.Watch = true
	return kc.restClient.
		Get().
		Namespace(kc.ns).
		Resource("kafkatopics").
		VersionedParams(&opt, scheme.ParameterCodec).
		Watch()
}

func (kc *kafkaTopicClusterClient) Delete(name string, opts metav1.ListOptions) error {
	return kc.restClient.
		Delete().
		Namespace(kc.ns).
		Resource("kafkatopics").
		VersionedParams(&opts, scheme.ParameterCodec).
		Name(name).
		Do().
		Error()
}
