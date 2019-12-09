package sync

import (
	"baas-service/models"
	"baas-service/pkg/e"
	"baas-service/pkg/k8s/client"
	"baas-service/pkg/utils"
	mysqlv1alpha1 "github.com/oracle/mysql-operator/pkg/apis/mysql/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// 在 k8s 集群里创建完整的 mysql 集群实例
func K8sCreateMysqlCluster(mysql *models.MysqlCluster) (int, error) {
	var (
		err     error
		code    int
		configs string
	)

	for key, value := range mysql.GenerateConfig() {
		configs += key + "=" + value + "\n"
	}

	data := make(map[string]string)
	data["my.cnf"] = "[mysqld]" + "\n" + configs

	err = createMysqlPasswordSecret(mysql.Namespace, mysql.SecretName(), mysql.Password)
	if err != nil {
		code = e.K8S_MYSQL_PASSWORD_SECRET_CREATE_FAILED
		return code, err
	}

	err = createMysqlConfig(mysql.Namespace, mysql.ConfigmapName(), data)
	if err != nil {
		code = e.K8S_MYSQL_CONFIGMAP_CREATE_FAILED
		return code, err
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
			Resources: &mysqlv1alpha1.Resources{
				Server: &corev1.ResourceRequirements{
					Requests: corev1.ResourceList{
						corev1.ResourceCPU:    resource.MustParse(mysql.CPU),
						corev1.ResourceMemory: resource.MustParse(mysql.Memory),
					},
					Limits: corev1.ResourceList{
						corev1.ResourceCPU:    resource.MustParse(mysql.CPU),
						corev1.ResourceMemory: resource.MustParse(mysql.Memory),
					},
				},
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

	_, err = client.MysqlClientset.Clusters(mysql.Namespace).Create(&newCluster)
	if err != nil {
		code = e.K8S_MYSQL_CLUSTER_CREATE_FAILED
		return code, err
	}

	err = createMysqlRouter(mysql.Namespace, mysql.RouterDeploymentName(), mysql.SecretName(), mysql.MysqlHostName(), mysql.Member, mysql.Port)
	if err != nil {
		code = e.K8S_MYSQL_ROUTER_CREATE_FAILED
		return code, err
	}

	err = createMysqlRouterService(mysql.Namespace, mysql.RouterDeploymentName(), mysql.RouterDeploymentName(), mysql.Port)
	if err != nil {
		code = e.K8S_MYSQL_ROUTER_SERVICE_CREATE_FAILED
		return code, err
	}

	code = e.CREATED
	return code, nil
}

// 创建 mysql 密码
func createMysqlPasswordSecret(ns, name, passwd string) error {
	stringData := make(map[string]string)
	stringData["password"] = passwd
	passwordSecret := corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		StringData: stringData,
	}
	_, err := client.K8sClient.CoreV1().Secrets(ns).Create(&passwordSecret)
	return err
}

// 创建 mysql 配置
func createMysqlConfig(ns, name string, data map[string]string) error {
	mysqlConfig := corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Data: data,
	}
	_, err := client.K8sClient.CoreV1().ConfigMaps(ns).Create(&mysqlConfig)
	return err
}

// 部署 mysql router 实例，用于自动识别 读/写 请求
func createMysqlRouter(ns, name, passwdSecretName, host string, num int, port int) error {

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

	_, err := client.K8sClient.AppsV1().Deployments(ns).Create(&routerDeployment)
	return err
}

// 创建 mysql router 的 service 由于暴露 Mysql service 提供用户连接的接口
func createMysqlRouterService(ns, name, appName string, port int) error {

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

	_, err := client.K8sClient.CoreV1().Services(ns).Create(&mysqlRouterService)
	return err
}

// 在 k8s 中删除指定的 mysql 实例
func K8sDeleteMysqlCluster(mysql *models.MysqlCluster) (int, error) {
	var err error
	var code int
	err = client.MysqlClientset.Clusters(mysql.Namespace).Delete(mysql.ClusterName, metav1.ListOptions{})
	if err != nil {
		code = e.K8S_MYSQL_CLUSTER_DELETE_FAILED
		return code, err
	}
	err = deleteMysqlPasswordSecret(mysql)
	if err != nil {
		code = e.K8S_MYSQL_PASSWORD_SECRET_DELETE_FAILED
		return code, err
	}
	err = deleteMysqlConfig(mysql)
	if err != nil {
		code = e.K8S_MYSQL_CONFIGMAP_DELETE_FAILED
		return code, err
	}
	err = deleteMysqlRouter(mysql)
	if err != nil {
		code = e.K8S_MYSQL_ROUTER_DELETE_FAILED
		return code, err
	}
	err = deleteMysqlRouterService(mysql)
	if err != nil {
		code = e.K8S_MYSQL_ROUTER_SERVICE_DELETE_FAILED
		return code, err
	}
	err = deleteMysqlPVCs(mysql)
	if err != nil {
		code = e.K8S_MYSQL_PVC_DELETE_FAILED
		return code, err
	}
	code = e.MYSQL_DELETED
	return code, nil
}

func deleteMysqlPasswordSecret(mysql *models.MysqlCluster) error {
	err := client.K8sClient.CoreV1().Secrets(mysql.Namespace).Delete(mysql.SecretName(), &metav1.DeleteOptions{})
	if err != nil {
		return err
	}
	return nil
}

func deleteMysqlConfig(mysql *models.MysqlCluster) error {
	err := client.K8sClient.CoreV1().ConfigMaps(mysql.Namespace).Delete(mysql.ConfigmapName(), &metav1.DeleteOptions{})
	if err != nil {
		return err
	}
	return nil
}

func deleteMysqlRouter(mysql *models.MysqlCluster) error {
	err := client.K8sClient.AppsV1().Deployments(mysql.Namespace).Delete(mysql.RouterDeploymentName(), &metav1.DeleteOptions{})
	if err != nil {
		return err
	}
	return nil
}

func deleteMysqlRouterService(mysql *models.MysqlCluster) error {
	err := client.K8sClient.CoreV1().Services(mysql.Namespace).Delete(mysql.RouterDeploymentName(), &metav1.DeleteOptions{})
	if err != nil {
		return err
	}
	return nil
}

func deleteMysqlPVCs(mysql *models.MysqlCluster) error {
	for _, name := range mysql.PVCNames() {
		err := client.K8sClient.CoreV1().PersistentVolumeClaims(mysql.Namespace).Delete(name, &metav1.DeleteOptions{})
		if err != nil {
			return err
		}
	}
	return nil
}
