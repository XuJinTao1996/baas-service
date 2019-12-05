package sync

import (
	"baas-service/models"
	"baas-service/pkg/k8s/client"
	"baas-service/pkg/utils"
	"github.com/google/martian/log"
	mysqlv1alpha1 "github.com/oracle/mysql-operator/pkg/apis/mysql/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var (
	ClusterConfig  *rest.Config
	MysqlClientset *client.ClusterClient
	K8sClient      *kubernetes.Clientset
)

func init() {
	ClusterConfig = client.LoadK8sConfig()
	MysqlClientset, _ = client.NewForConfig(ClusterConfig, mysqlv1alpha1.GroupName, "v1alpha1")
	K8sClient, _ = kubernetes.NewForConfig(ClusterConfig)
}

func K8sGetMysqlClsuter() *mysqlv1alpha1.ClusterList {
	mysqlcluster, err := MysqlClientset.Clusters("mysql-operator").List(metav1.ListOptions{})
	if err != nil {
		panic(err)
	}

	return mysqlcluster
}

// 在 k8s 集群里创建完整的 mysql 集群实例
func K8sCreateMysqlCluster(mysql *models.MysqlCluster) (*mysqlv1alpha1.Cluster, error) {

	data := make(map[string]string)
	data["my.cnf"] = "[mysqld]" + "\n" +
		"default_authentication_plugin=mysql_native_password"

	secret, secretErr := createMysqlPasswordSecret(mysql.Namespace, mysql.SecretName(), mysql.Password)
	if secretErr != nil {
		log.Errorf("failed to create secret: %v", secretErr)
	} else {
		log.Infof("%v", secret)
	}

	configmap, configErr := createMysqlConfig(mysql.Namespace, mysql.ConfigmapName(), data)
	if configErr != nil {
		log.Errorf("failed to create configmap: %v", configErr)
	} else {
		log.Infof("%v", configmap)
	}

	newCluster := mysqlv1alpha1.Cluster{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Cluster",
			APIVersion: "mysql.oracle.com/v1alpha1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: mysql.ClusterName,
		},
		Spec: mysqlv1alpha1.ClusterSpec{
			Members:     int32(mysql.Member),
			MultiMaster: mysql.MultiMaster,
			Version:     mysql.Version,
			RootPasswordSecret: &corev1.LocalObjectReference{
				Name: mysql.SecretName(),
			},
			Config: &corev1.LocalObjectReference{
				Name: mysql.ConfigmapName(),
			},
			VolumeClaimTemplate: &corev1.PersistentVolumeClaim{
				ObjectMeta: metav1.ObjectMeta{
					Name: "data",
				},
				Spec: corev1.PersistentVolumeClaimSpec{
					AccessModes: []corev1.PersistentVolumeAccessMode{corev1.ReadWriteMany},
					Resources: corev1.ResourceRequirements{
						Requests: corev1.ResourceList{
							corev1.ResourceStorage: resource.MustParse(mysql.VolumeSize),
						},
					},
					StorageClassName: &mysql.StorageType,
				},
			},
		},
	}

	result, err := MysqlClientset.Clusters(mysql.Namespace).Create(&newCluster)
	if err != nil {
		panic(err)
	}

	mysqlRouter, routerErr := createMysqlRouter(mysql.Namespace, mysql.RouterDeploymentName(), mysql.SecretName(), mysql.MysqlHostName(), mysql.Member, mysql.Port)
	if routerErr != nil {
		log.Errorf("failed create mysql router: %v", routerErr)
	} else {
		log.Infof("%v", mysqlRouter)
	}

	mysqlRouterService, serviceErr := createMysqlRouterService(mysql.Namespace, mysql.RouterDeploymentName(), mysql.RouterDeploymentName(), mysql.Port)
	if serviceErr != nil {
		log.Errorf("failed create mysql router service: %v", serviceErr)
	} else {
		log.Infof("%v", mysqlRouterService)
	}

	return result, err
}

// 创建 mysql 密码
func createMysqlPasswordSecret(ns, name, passwd string) (*corev1.Secret, error) {
	stringData := make(map[string]string)
	stringData["password"] = passwd
	passwordSecret := corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		StringData: stringData,
	}
	result, err := K8sClient.CoreV1().Secrets(ns).Create(&passwordSecret)
	return result, err
}

// 创建 mysql 配置
func createMysqlConfig(ns, name string, data map[string]string) (*corev1.ConfigMap, error) {
	mysqlConfig := corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Data: data,
	}
	result, err := K8sClient.CoreV1().ConfigMaps(ns).Create(&mysqlConfig)
	return result, err
}

// 部署 mysql router 实例，用于自动识别 读/写 请求
func createMysqlRouter(ns, name, passwdSecretName, host string, num int, port int) (*appsv1.Deployment, error) {

	labelSelector := make(map[string]string)
	labelSelector["app"] = name

	routerDeployment := appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: labelSelector,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name:   name,
					Labels: labelSelector,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Name:            name,
						Image:           "mysql/mysql-router:8.0.12",
						ImagePullPolicy: "IfNotPresent",
						Command:         []string{"/bin/bash", "-cx", "exec /run.sh mysqlrouter"},
						Env: []corev1.EnvVar{
							{
								Name: "MYSQL_PASSWORD",
								ValueFrom: &corev1.EnvVarSource{
									SecretKeyRef: &corev1.SecretKeySelector{
										LocalObjectReference: corev1.LocalObjectReference{
											Name: passwdSecretName,
										},
										Key: "password",
									},
								},
							},
							{
								Name:  "MYSQL_USER",
								Value: "root",
							},
							{
								Name:  "MYSQL_PORT",
								Value: utils.Int(port).Str(),
							},
							{
								Name:  "MYSQL_HOST",
								Value: host,
							},
							{
								Name:  "MYSQL_INNODB_NUM_MEMBERS",
								Value: utils.Int(num).Str(),
							},
						},
					}},
				},
			},
		},
	}

	result, err := K8sClient.AppsV1().Deployments(ns).Create(&routerDeployment)
	return result, err
}

// 创建 mysql router 的 service 由于暴露 Mysql service 提供用户连接的接口
func createMysqlRouterService(ns, name, appName string, port int) (*corev1.Service, error) {

	labelSelector := make(map[string]string)
	labelSelector["app"] = appName

	mysqlRouterService := corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name:       name,
					Protocol:   corev1.ProtocolTCP,
					Port:       int32(port),
					TargetPort: intstr.IntOrString{IntVal: 6446},
				},
			},
			Selector: labelSelector,
			Type:     "ClusterIP",
		},
	}

	result, err := K8sClient.CoreV1().Services(ns).Create(&mysqlRouterService)
	return result, err
}

// 在 k8s 中删除指定的 mysql 实例
func K8sDeleteMysqlCluster(mysql *models.MysqlCluster) {
	err := MysqlClientset.Clusters(mysql.Namespace).Delete(mysql.ClusterName, metav1.ListOptions{})
	if err != nil {
		log.Errorf("failed to delete mysqlcluster %v", err)
	}
	deleteMysqlPasswordSecret(mysql)
	deleteMysqMysqlConfig(mysql)
	deleteMysqMysqlRouter(mysql)
	deleteMysqlRouterService(mysql)
}

func deleteMysqlPasswordSecret(mysql *models.MysqlCluster) {
	err := K8sClient.CoreV1().Secrets(mysql.Namespace).Delete(mysql.SecretName(), &metav1.DeleteOptions{})
	if err != nil {
		log.Errorf("failed to delete Secrets%v", err)
	}
}

func deleteMysqMysqlConfig(mysql *models.MysqlCluster) {
	err := K8sClient.CoreV1().ConfigMaps(mysql.Namespace).Delete(mysql.ConfigmapName(), &metav1.DeleteOptions{})
	if err != nil {
		log.Errorf("failed to delete MysqMysqlConfig %v", err)
	}
}

func deleteMysqMysqlRouter(mysql *models.MysqlCluster) {
	err := K8sClient.AppsV1().Deployments(mysql.Namespace).Delete(mysql.RouterDeploymentName(), &metav1.DeleteOptions{})
	if err != nil {
		log.Errorf("failed to delete Deployments %v", err)
	}
}

func deleteMysqlRouterService(mysql *models.MysqlCluster) {
	err := K8sClient.CoreV1().Services(mysql.Namespace).Delete(mysql.RouterDeploymentName(), &metav1.DeleteOptions{})
	if err != nil {
		log.Errorf("failed to delete Services %v", err)
	}
}

func SyncMysqlClusterStatus(name, status string) bool {
	mysqCluster, state := models.GetMysqlclusterByName(name)
	if !state {
		log.Errorf("mysql cluster does not exist!")
		return false
	}
	newMysqlCluster, result := models.UpdateMysqlCluster(mysqCluster, status)
	if !result {
		log.Errorf("mysql cluster %v update failed", newMysqlCluster.ClusterName)
		return false
	}
	return true
}
