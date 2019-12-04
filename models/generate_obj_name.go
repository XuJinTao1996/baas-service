package models

type MysqlObjectNameInterface interface {
	SecretName() string
	ConfigmapName() string
	RouterDeploymentName() string
	MysqlHostName() string
}

func (obj *MysqlCluster) SecretName() string {
	return obj.ClusterName + "-" + obj.User + "-password"
}

func (obj *MysqlCluster) ConfigmapName() string {
	return obj.ClusterName + "-config"
}

func (obj *MysqlCluster) RouterDeploymentName() string {
	return obj.ClusterName + "-router"
}

func (obj *MysqlCluster) MysqlHostName() string {
	return obj.ClusterName + "-0" + "." + obj.ClusterName
}
