package informer

import (
	"baas-service/pkg/k8s/client"
	"baas-service/pkg/sync"
	"fmt"
	"github.com/google/martian/log"
	mysqlv1alpha1 "github.com/oracle/mysql-operator/pkg/apis/mysql/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"time"
)

func MysqlClusterInformer(clientSet client.BaseClientInterface) {
	_, clusterController := cache.NewInformer(&cache.ListWatch{
		ListFunc: func(opt metav1.ListOptions) (result runtime.Object, e error) {
			return clientSet.Clusters("mysql-operator").List(opt)
		},
		WatchFunc: func(opt metav1.ListOptions) (i watch.Interface, e error) {
			return clientSet.Clusters("mysql-operator").Watch(opt)
		},
	},
		&mysqlv1alpha1.Cluster{},
		1*time.Minute,
		cache.ResourceEventHandlerFuncs{
			AddFunc: onAddMysql,
			UpdateFunc: func(oldObj, newObj interface{}) {
				fmt.Println("update mysqlcluster", oldObj.(*mysqlv1alpha1.Cluster).Status)
			},
		},
	)

	go clusterController.Run(wait.NeverStop)
}

func onAddMysql(obj interface{}) {
	mysqlcluster := obj.(*mysqlv1alpha1.Cluster)
	fmt.Println("add a mysqlcluster", mysqlcluster.Name)
}

func DeploymentInformer(clientSet *kubernetes.Clientset) {
	factory := informers.NewSharedInformerFactory(clientSet, 0)
	deploymentInformer := factory.Apps().V1().Deployments().Informer()
	defer utilruntime.HandleCrash()

	// 启动 informer, list & watch
	go factory.Start(wait.NeverStop)

	if !cache.WaitForCacheSync(wait.NeverStop, deploymentInformer.HasSynced) {
		utilruntime.HandleError(fmt.Errorf("Timed out waiting for caches to sync"))
		return
	}

	// 使用自定义事件 handler
	deploymentInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    onAdd,
		UpdateFunc: onUpdate,
		DeleteFunc: func(interface{}) { fmt.Println("delete not implemented") },
	})
}

func onAdd(obj interface{}) {
	deployment := obj.(*appsv1.Deployment)
	fmt.Println("add a deployment", deployment.Name)
}

// update mysql cluster
func onUpdate(oldobj, newobj interface{}) {
	mysqlCluster := newobj.(*mysqlv1alpha1.Cluster)
	state := mysqlCluster.Status.Conditions[0].Type
	result := sync.SyncMysqlClusterStatus(mysqlCluster.Name, string(state))
	if result {
		log.Infof("success update mysql cluster %v", mysqlCluster.Name)
	}
}
