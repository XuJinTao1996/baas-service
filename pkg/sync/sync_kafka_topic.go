package sync

import (
	"baas-service/models"
	"baas-service/pkg/e"
	"baas-service/pkg/k8s/client"
	kafkatopicv1alpha1 "github.com/banzaicloud/kafka-operator/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func K8sCreateKafkaTopic(kt *models.KafkaTopic) (int, error) {
	var code int
	var err error

	newKafkaTopic := kafkatopicv1alpha1.KafkaTopic{
		TypeMeta: metav1.TypeMeta{
			Kind:       "KafkaTopic",
			APIVersion: "kafka.banzaicloud.io/v1alpha1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: kt.Name,
		},
		Spec: kafkatopicv1alpha1.KafkaTopicSpec{
			ClusterRef: kafkatopicv1alpha1.ClusterReference{
				Name: kt.ClusterRefName,
			},
			Name:              kt.TopicName,
			Partitions:        int32(kt.Partitions),
			ReplicationFactor: int32(kt.ReplicationFactor),
		},
	}

	_, err = client.KafkaTopicClient.KafkaTopic(kt.Namespace).Create(&newKafkaTopic)
	if err != nil {
		code = e.K8S_KAFKA_TOPIC_CREATE_FAILED
		return code, err
	}
	code = e.CREATED
	return code, nil
}

func K8sDeleteKafkaTopic(kt *models.KafkaTopic) (int, error) {
	var code int
	err := client.KafkaTopicClient.KafkaTopic(kt.Namespace).Delete(kt.Name, metav1.ListOptions{})
	if err != nil {
		code = e.K8S_KAFKA_TOPIC_DELETE_FAILED
		return code, err
	}
	code = e.KAFKA_DELETED
	return code, nil
}
