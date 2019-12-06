package models

import (
	"baas-service/pkg/utils"
)

type MysqlObjectNameInterface interface {
	SecretName() string
	ConfigmapName() string
	RouterDeploymentName() string
	MysqlHostName() string
	PVCNames() string
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

func (obj *MysqlCluster) PVCNames() []string {
	var pvcNames []string
	for i := 0; i < obj.Member; i++ {
		pvcNames = append(pvcNames, "data-"+obj.ClusterName+"-"+utils.Int(i).Str())
	}
	return pvcNames
}
