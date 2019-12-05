package v1

import (
	"baas-service/models"
	"baas-service/pkg/e"
	"baas-service/pkg/sync"
	"baas-service/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/common/log"
	"net/http"
)

// 获取指定 mysql 实例
func GetMysqlCluster(c *gin.Context) {
	var code int
	id := utils.Str(c.Param("id")).Int()
	mysqlCluster, state := models.GetMysqlcluster(id)
	if !state {
		code = e.MYSQL_DOES_NOT_ESIST
		log.Errorf("cluster does not exist!")
	} else {
		code = e.SUCCESS
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": mysqlCluster,
	})
}

// 获取所有 mysql 集群
func ListMysqlCluster(c *gin.Context) {

	data := make(map[string]interface{})

	data["list"], data["total"] = models.GetAllMysqlClusters()

	c.JSON(http.StatusOK, gin.H{
		"code": e.SUCCESS,
		"msg":  e.GetMsg(e.SUCCESS),
		"data": data,
	})
}

// 创建 mysql 集群
func CreateMysqlCluster(c *gin.Context) {
	var mysqlCluster models.MysqlCluster

	err := c.ShouldBind(&mysqlCluster)
	if err != nil {
		log.Error(err)
	}

	mysqlCluster.Host = mysqlCluster.RouterDeploymentName()

	code := e.CREATED

	isExisted := models.ExistMysqlCluster(mysqlCluster.ClusterName)
	if isExisted {
		log.Infof("cluster %v existed!", mysqlCluster.ClusterName)
		code = e.ERROR_ESIST_MYSQL
	} else {
		cStatus := models.AddMysqlCluster(&mysqlCluster)
		if !cStatus {
			code = e.ERROR_CREATE_MYSQL
			log.Error("Failed to create mysql cluster!")
		} else {
			_, err := sync.K8sCreateMysqlCluster(&mysqlCluster)
			if err != nil {
				log.Error("Failed to create mysqlcluster %v", mysqlCluster.ClusterName)
				code = e.ERROR_CREATE_MYSQL
			} else {
				log.Infof("cluster %v created", mysqlCluster.ClusterName)
				code = e.CREATED
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": mysqlCluster,
	})

}

func UpdateMysqlCluster(c *gin.Context) {

}

func DeleteMysqlCluster(c *gin.Context) {
	var code int
	id := utils.Str(c.Param("id")).Int()
	mysqlCluster, state := models.GetMysqlcluster(id)
	if !state {
		log.Error("mysql cluster does not existed")
		code = e.MYSQL_DOES_NOT_ESIST
	} else {
		sync.K8sDeleteMysqlCluster(&mysqlCluster)
		models.DeleteMysqlcluster(id)
		code = e.MYSQL_DELETED
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": nil,
	})
}
