package sync

import (
	"baas-service/models"
	"baas-service/pkg/e"
	"baas-service/pkg/k8s/client"
	"github.com/pravega/zookeeper-operator/pkg/apis/zookeeper/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func K8sCreateZookeeperCluster(zc *models.ZookeeperCluster) (int, error) {
	var code int
	newZookeeperCluster := v1beta1.ZookeeperCluster{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ZookeeperCluster",
			APIVersion: "zookeeper.pravega.io/v1beta1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: zc.ClusterName,
		},
		Spec: v1beta1.ZookeeperClusterSpec{
			Size: int32(zc.Size),
		},
	}

	_, err := client.ZookeeperClientset.ZookeeperCluster(zc.Namespace).Create(&newZookeeperCluster)
	if err != nil {
		code = e.K8S_ZOOKEEPER_CLUSTER_CREATE_FAILED
		return code, err
	}
	code = e.CREATED
	return code, nil
}

func K8sDeleteZookeeperCluster(zc *models.ZookeeperCluster) (int, error) {
	var code int
	err := client.ZookeeperClientset.ZookeeperCluster(zc.Namespace).Delete(zc.ClusterName, metav1.ListOptions{})
	if err != nil {
		code = e.K8S_ZOOKEEPER_CLUSTER_DELETE_FAILED
		return code, err
	}

	code = e.ZOOKEEPER_DELETED
	return code, nil
}
