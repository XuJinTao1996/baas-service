package informer

import (
	"baas-service/models"
	"baas-service/pkg/e"
	"baas-service/pkg/k8s/client"
	"fmt"
	kafkav1beta1 "github.com/banzaicloud/kafka-operator/api/v1beta1"
	mysqlv1alpha1 "github.com/oracle/mysql-operator/pkg/apis/mysql/v1alpha1"
	zookeeperv1beta1 "github.com/pravega/zookeeper-operator/pkg/apis/zookeeper/v1beta1"
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
			return clientSet.MysqlClusters("").List(opt)
		},
		WatchFunc: func(opt metav1.ListOptions) (i watch.Interface, e error) {
			return clientSet.MysqlClusters("").Watch(opt)
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
	log.Println("add a mysql cluster: ", mysqlcluster.Name)
}

func onAddDeployment(obj interface{}) {
	deployment := obj.(*appsv1.Deployment)
	log.Println("add a deployment: ", deployment.Name)
}

func onUpdateDeployment(oldObj, newObj interface{}) {
	deployment := newObj.(*appsv1.Deployment)
	log.Println("add a deployment: ", deployment.Name)
}

// 更新 mysql cluster 的状态
func onUpdateMysqlCluster(oldObj, newObj interface{}) {
	k8sMysqlCluster := oldObj.(*mysqlv1alpha1.Cluster)
	state := k8sMysqlCluster.Status.Conditions
	if len(state) > 0 {

		exists, err := models.ExistMysqlClusterByName(k8sMysqlCluster.Name)
		if err != nil {
			log.Println(e.GetMsg(e.ERROR_CHECK_MYSQL_EXIST_FAIL))
			return
		}

		if !exists {
			log.Println(e.MYSQL_DOES_NOT_EXIST)
			return
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
			log.Printf("mysql cluster %v update failed: ", newMysqlCluster.ClusterName)
			return
		}
	}
}

func ZookeeperClusterInformer(clientSet client.BaseClientInterface) {
	_, clusterController := cache.NewInformer(&cache.ListWatch{
		ListFunc: func(opt metav1.ListOptions) (result runtime.Object, e error) {
			return clientSet.ZookeeperCluster("").List(opt)
		},
		WatchFunc: func(opt metav1.ListOptions) (i watch.Interface, e error) {
			return clientSet.ZookeeperCluster("").Watch(opt)
		},
	},
		&zookeeperv1beta1.ZookeeperCluster{},
		1*time.Minute,
		cache.ResourceEventHandlerFuncs{
			AddFunc:    onAddZookeeperCluster,
			UpdateFunc: onUpdateZookeeperCluster,
		},
	)

	go clusterController.Run(wait.NeverStop)
}

func onAddZookeeperCluster(obj interface{}) {
	zookeeperCluster := obj.(*zookeeperv1beta1.ZookeeperCluster)
	log.Println("add a zookeeper cluster: ", zookeeperCluster.Name)
}

func onUpdateZookeeperCluster(oldobj, newobj interface{}) {

	var newZookeeperCluster *models.ZookeeperCluster
	var result bool

	zookeeperObj := oldobj.(*zookeeperv1beta1.ZookeeperCluster)
	readyMembers := zookeeperObj.Status.ReadyReplicas

	exists, err := models.ExistZookeeperClusterByName(zookeeperObj.Name)
	if err != nil {
		log.Println(e.GetMsg(e.ERROR_CHECK_ZOOKEEPER_EXIST_FAIL))
		return
	}

	if !exists {
		log.Println(e.GetMsg(e.ZOOKEEPER_DOES_NOT_EXIST))
		return
	}

	zookeeperCluster, err := models.GetZookeeperClusterByName(zookeeperObj.Name)
	if err != nil {
		log.Println(e.GetMsg(e.ERROR_GET_ZOOKEEPER_CLUSTER_FAIL))
		return
	}

	if readyMembers == zookeeperObj.Spec.Size {
		newZookeeperCluster, result = models.UpdateZookeeperCluster(zookeeperCluster, "Ready")
	} else {
		newZookeeperCluster, result = models.UpdateZookeeperCluster(zookeeperCluster, "NotReady")
	}

	if !result {
		log.Printf("zookeeper cluster %v update failed ", newZookeeperCluster.ClusterName)
	}
}

func KafkaClusterInformer(clientSet client.BaseClientInterface) {
	_, clusterController := cache.NewInformer(&cache.ListWatch{
		ListFunc: func(opt metav1.ListOptions) (result runtime.Object, e error) {
			return clientSet.KafkaCluster("").List(opt)
		},
		WatchFunc: func(opt metav1.ListOptions) (i watch.Interface, e error) {
			return clientSet.KafkaCluster("").Watch(opt)
		},
	},
		&kafkav1beta1.KafkaCluster{},
		1*time.Minute,
		cache.ResourceEventHandlerFuncs{
			AddFunc:    onAddKafkaCluster,
			UpdateFunc: onUpdateKafkaCluster,
		},
	)

	go clusterController.Run(wait.NeverStop)
}

func onAddKafkaCluster(obj interface{}) {
	kafkaCluster := obj.(*kafkav1beta1.KafkaCluster)
	log.Printf("add a zookeeper cluster: %v", kafkaCluster.Name)
}

func onUpdateKafkaCluster(oldobj, newobj interface{}) {

	var result bool

	kafkaObj := oldobj.(*kafkav1beta1.KafkaCluster)
	state := kafkaObj.Status.State

	exists, err := models.ExistKafkaClusterByName(kafkaObj.Name)
	if err != nil {
		log.Println(e.GetMsg(e.ERROR_CHECK_KAFKA_EXIST_FAIL))
		return
	}

	if !exists {
		log.Println(e.GetMsg(e.KAFKA_DOES_NOT_EXIST))
		return
	}

	kafkaCluster, err := models.GetKafkaClusterByName(kafkaObj.Name)
	if err != nil {
		log.Println(e.GetMsg(e.ERROR_GET_ZOOKEEPER_CLUSTER_FAIL))
		return
	}

	fmt.Println(state, kafkaCluster.Status)
	if string(state) != kafkaCluster.Status {
		_, result = models.UpdateKafkaCluster(kafkaCluster, string(state))
	}

	if !result {
		log.Printf("zookeeper cluster %v update failed ", kafkaCluster.ClusterName)
	}
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
