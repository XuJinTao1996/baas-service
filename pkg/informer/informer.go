package informer

import (
	"baas-service/pkg/k8s/client"
	mysqlv1alpha1 "github.com/oracle/mysql-operator/pkg/apis/mysql/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
	"time"
)

func WatchResources(clientSet client.BaseClientInterface) cache.Store {
	clusterStore, clusterController := cache.NewInformer(&cache.ListWatch{
		ListFunc: func(lo metav1.ListOptions) (object runtime.Object, e error) {
			return clientSet.Clusters("mysql-operator").List(lo)
		},
		WatchFunc: func(lo metav1.ListOptions) (i watch.Interface, e error) {
			return clientSet.Clusters("mysql-operator").Watch(lo)
		},
	},
		&mysqlv1alpha1.Cluster{},
		1*time.Minute,
		cache.ResourceEventHandlerFuncs{},
	)

	go clusterController.Run(wait.NeverStop)
	return clusterStore
}
