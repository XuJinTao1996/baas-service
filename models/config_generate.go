package models

import "baas-service/pkg/utils"

type ConfigGenerateInterface interface {
	GenerateConfig(cluster *MysqlCluster) string
}

func (cluster *MysqlCluster) GenerateConfig() map[string]string {
	configs := make(map[string]string)
	configs["default_authentication_plugin"] = cluster.DefaultAuthenticationPlugin
	configs["max_connections"] = utils.Int(cluster.MaxConnections).Str()
	return configs
}
