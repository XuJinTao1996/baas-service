package v1

import (
	"github.com/gin-gonic/gin"
	mysqlv1alpha1 "github.com/oracle/mysql-operator/pkg/apis/mysql/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"rds-front/pkg/k8s_client"

	"k8s.io/client-go/rest"
	"net/http"
	//"rds-front/models"
	"rds-front/pkg/e"
	//"rds-front/pkg/util"
)

func load_k8s_config() *k8s_client.ClusterClient {
	var config *rest.Config
	config, _ = rest.InClusterConfig()
	client_set, err := k8s_client.NewForConfig(config)
	if err != nil {
		panic(err)
	}
	return client_set
}

func K8s_get_mysqlcluster() *mysqlv1alpha1.ClusterList {
	client_set := load_k8s_config()
	mysqlcluster, err := client_set.Clusters("mysql-operator").List(metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	return mysqlcluster
}

// 获取所有 mysql 实例
func GetMysqlInstances(c *gin.Context) {

	data := make(map[string]interface{})

	//data["list"], data["total"] = models.GetAllMysqlInstance()
	data["list"] = K8s_get_mysqlcluster()

	c.JSON(http.StatusOK, gin.H{
		"code": e.SUCCESS,
		"msg":  e.GetMsg(e.SUCCESS),
		"data": data,
	})
}

func CreateMysqlInstance(c *gin.Context) {
	name := c.PostForm("db_instance_name")
	data := make(map[string]interface{})
	client_set := load_k8s_config()
	//var pvc_mode []corev1.PersistentVolumeAccessMode{corev1.ReadWriteMany,}
	//sc_name := "alicloud-nas"
	newCluster := mysqlv1alpha1.Cluster{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Cluster",
			APIVersion: "mysql.oracle.com/v1alpha1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: mysqlv1alpha1.ClusterSpec{
			Version: "8.0.12",
			//VolumeClaimTemplate: &corev1.PersistentVolumeClaim{
			//	ObjectMeta: metav1.ObjectMeta{
			//		Name: "test-mysql",
			//	},
			//	Spec: corev1.PersistentVolumeClaimSpec{
			//		AccessModes: pvc_mode,
			//		Resources: corev1.ResourceRequirements{
			//			Requests: corev1.ResourceList{corev1.ResourceStorage}.StorageEphemeral(),
			//		},
			//		StorageClassName: &sc_name,
			//	},
			//},
			RootPasswordSecret: &corev1.LocalObjectReference{
				Name: "wordpress-mysql-root-password",
			},
			Config: &corev1.LocalObjectReference{
				Name: "mysql-config",
			},
		},
	}
	result, err := client_set.Clusters("mysql-operator").Create(&newCluster)
	if err != nil {
		panic(err)
	}

	//db_user := c.PostForm("db_user")
	//db_passwd := c.PostForm("db_passwd")
	//db_port := util.Str(c.PostForm("db_port")).Int()
	//deployment_mode := c.PostForm("deployment_mode")
	//code := e.INVALID_PARAMS
	//if !models.ExistMysqlInstance(name) {
	//	models.AddTag(name, db_user, db_passwd, db_port, deployment_mode)
	code := e.CREATED
	//} else {
	//	code = e.ERROR_ESIST_MYSQL
	//}
	data["data"] = result

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": data,
	})
}

func UpdateMysqlInstance(c *gin.Context) {

}

func DeleteMysqlInstance(c *gin.Context) {

}
