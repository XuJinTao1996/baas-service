package informer

import (
	"baas-service/models"
	"baas-service/pkg/e"
	"baas-service/pkg/k8s/client"
	"fmt"
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
	"log"
	"time"
)

// 对 mysql cluster 资源的 list & watch
func MysqlClusterInformer(clientSet client.BaseClientInterface) {
	_, clusterController := cache.NewInformer(&cache.ListWatch{
		ListFunc: func(opt metav1.ListOptions) (result runtime.Object, e error) {
			return clientSet.Clusters("").List(opt)
		},
		WatchFunc: func(opt metav1.ListOptions) (i watch.Interface, e error) {
			return clientSet.Clusters("").Watch(opt)
		},
	},
		&mysqlv1alpha1.Cluster{},
		1*time.Minute,
		cache.ResourceEventHandlerFuncs{
			AddFunc:    onAddMysql,
			UpdateFunc: onUpdateMysqlCluster,
		},
	)

	go clusterController.Run(wait.NeverStop)
}

func onAddMysql(obj interface{}) {
	mysqlcluster := obj.(*mysqlv1alpha1.Cluster)
	log.Println("add a mysql cluster", mysqlcluster.Name)
}

func DeploymentInformer(clientSet *kubernetes.Clientset) {
	factory := informers.NewSharedInformerFactory(clientSet, 0)
	deploymentInformer := factory.Apps().V1().Deployments().Informer()
	defer utilruntime.HandleCrash()

	// 启动 informer, list & watch
	go factory.Start(wait.NeverStop)

	if !cache.WaitForCacheSync(wait.NeverStop, deploymentInformer.HasSynced) {
		utilruntime.HandleError(fmt.Errorf("Time out waiting for caches to sync"))
		return
	}

	// 使用自定义事件 handler
	deploymentInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    onAddDeployment,
		UpdateFunc: onUpdateDeployment,
	})
}

func onAddDeployment(obj interface{}) {
	deployment := obj.(*appsv1.Deployment)
	log.Println("add a deployment", deployment.Name)
}

func onUpdateDeployment(oldObj, newObj interface{}) {
	deployment := newObj.(*appsv1.Deployment)
	log.Println("add a deployment", deployment.Name)
}

// 更新 mysql cluster 的状态
func onUpdateMysqlCluster(oldObj, newObj interface{}) {
	k8sMysqlCluster := oldObj.(*mysqlv1alpha1.Cluster)
	state := k8sMysqlCluster.Status.Conditions
	if len(state) > 0 {

		exists, err := models.ExistMysqlClusterByName(k8sMysqlCluster.Name)
		if err != nil {
			log.Println(e.GetMsg(e.ERROR_CHECK_MYSQL_EXIST_FAIL))
		}

		if !exists {
			log.Println(e.MYSQL_DOES_NOT_EXIST)
		}

		mysqlCluster, err := models.GetMysqlclusterByName(k8sMysqlCluster.Name)
		if err != nil {
			log.Println(e.ERROR_GET_MYSQL_CLUSTER_FAIL)
			return
		}

		if string(state[0].Type) == mysqlCluster.Status {
			log.Println("mysql cluster status has no changed!")
			return
		}

		newMysqlCluster, result := models.UpdateMysqlCluster(mysqlCluster, string(state[0].Type))
		if !result {
			log.Println("mysql cluster %v update failed", newMysqlCluster.ClusterName)
			return
		}
	}
}
