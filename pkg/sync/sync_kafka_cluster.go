package sync

import (
	"baas-service/models"
	"baas-service/pkg/e"
	"baas-service/pkg/k8s/client"
	"baas-service/pkg/utils"
	kafkav1beta1 "github.com/banzaicloud/kafka-operator/api/v1beta1"
	"io/ioutil"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	KAFKA_CONFIG_FILE          = "assets/kafka.conf"
	KAFKA_CAPACITY_CONFIG_FILE = "assets/capacity.conf"
	KAFKA_CLUSTER_CONFIG_FILE  = "assets/cluster.conf"
)

func ReadConfigfile(filePath string) string {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	return string(content)
}

func K8sCreateKafkaCluster(kc *models.KafkaCluster) (int, error) {
	var code int
	var err error

	brokerConfigGroups := make(map[string]kafkav1beta1.BrokerConfig)
	brokerConfigGroups["default"] = kafkav1beta1.BrokerConfig{
		StorageConfigs: []kafkav1beta1.StorageConfig{{
			MountPath: "/kafka-logs",
			PVCSpec: &corev1.PersistentVolumeClaimSpec{
				AccessModes: []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce},
				Resources: corev1.ResourceRequirements{
					Requests: corev1.ResourceList{
						corev1.ResourceStorage: resource.MustParse(utils.Int(kc.StorageSize).Str() + "Gi"),
					},
				},
				StorageClassName: &kc.StorageType,
			},
		}},
	}

	config := ReadConfigfile(KAFKA_CONFIG_FILE)
	capacityConfig := ReadConfigfile(KAFKA_CAPACITY_CONFIG_FILE)
	clusterConfig := ReadConfigfile(KAFKA_CLUSTER_CONFIG_FILE)

	zk, err := models.GetZookeeperClusterByID(kc.ZookeeperClusterID)
	if err != nil {
		return e.ERROR_GET_ZOOKEEPER_CLUSTER_FAIL, err
	}

	newKafkaCluster := kafkav1beta1.KafkaCluster{
		TypeMeta: metav1.TypeMeta{
			Kind:       "KafkaCluster",
			APIVersion: "kafka.banzaicloud.io/v1beta1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: kc.ClusterName,
		},
		Spec: kafkav1beta1.KafkaClusterSpec{
			HeadlessServiceEnabled: true,
			ZKAddresses:            []string{zk.ClusterName + "-" + "client:2181"},
			OneBrokerPerNode:       kc.OnBrokerPerNode,
			ClusterImage:           kc.ClusterImage,
			ReadOnlyConfig:         "auto.create.topics.enable=false",
			BrokerConfigGroups:     brokerConfigGroups,
			Brokers: []kafkav1beta1.Broker{
				{
					Id:                0,
					BrokerConfigGroup: "default",
				}, {
					Id:                1,
					BrokerConfigGroup: "default",
				},
				{
					Id:                2,
					BrokerConfigGroup: "default",
				}},
			RollingUpgradeConfig: kafkav1beta1.RollingUpgradeConfig{
				FailureThreshold: 1,
			},
			ListenersConfig: kafkav1beta1.ListenersConfig{
				InternalListeners: []kafkav1beta1.InternalListenerConfig{
					{
						Type:                            "plaintext",
						Name:                            "internal",
						ContainerPort:                   29092,
						UsedForInnerBrokerCommunication: true,
					},
					{
						Type:                            "plaintext",
						Name:                            "controller",
						ContainerPort:                   29093,
						UsedForInnerBrokerCommunication: true,
						UsedForControllerCommunication:  true,
					},
				},
			},
			CruiseControlConfig: kafkav1beta1.CruiseControlConfig{
				TopicConfig: &kafkav1beta1.TopicConfig{
					Partitions:        12,
					ReplicationFactor: 3,
				},
				Config:         config,
				CapacityConfig: capacityConfig,
				ClusterConfig:  clusterConfig,
			},
		},
	}

	_, err = client.KafkaClient.KafkaCluster(kc.Namespace).Create(&newKafkaCluster)
	if err != nil {
		code = e.K8S_KAFKA_CLUSTER_CREATE_FAILED
		return code, err
	}
	code = e.CREATED
	return code, nil
}

func K8sDeleteKafkaCluster(kc *models.KafkaCluster) (int, error) {
	var code int
	err := client.KafkaClient.KafkaCluster(kc.Namespace).Delete(kc.ClusterName, metav1.ListOptions{})
	if err != nil {
		code = e.K8S_KAFKA_CLUSTER_DELETE_FAILED
		return code, err
	}
	code = e.KAFKA_DELETED
	return code, nil
}
