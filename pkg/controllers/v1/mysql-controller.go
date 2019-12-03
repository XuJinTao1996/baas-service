package v1

import (
	"baas-service/models"
	"baas-service/pkg/e"
	"baas-service/pkg/sync_status"
	"baas-service/pkg/util"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/common/log"
	"net/http"
)

// 获取所有 mysql 集群
func GetMysqlClusters(c *gin.Context) {

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

	data := make(map[string]interface{})

	namespace := c.PostForm("namespace")
	clusterName := c.PostForm("clusterName")
	member := util.Str(c.PostForm("member")).Int()
	user := c.PostForm("user")
	password := c.PostForm("password")
	port := util.Str(c.PostForm("port")).Int()
	multiMaster := util.Str(c.PostForm("multiMaster")).Bool()
	version := c.PostForm("version")
	storageType := c.PostForm("storageType")
	volumeSize := util.Str(c.PostForm("volumeSize")).Int()

	mysqlCluster := models.MysqlCluster{
		Namespace:   namespace,
		ClusterName: clusterName,
		Member:      member,
		User:        user,
		Passwd:      password,
		Port:        port,
		MultiMaster: multiMaster,
		Version:     version,
		StorageType: storageType,
		VolumeSize:  volumeSize,
		ServiceUrl:  clusterName + "-router",
	}

	code := e.CREATED

	isExisted := models.ExistMysqlCluster(clusterName)
	if isExisted {
		log.Infof("cluster %v existed!", clusterName)
		code = e.ERROR_ESIST_MYSQL
	} else {
		cStatus := models.AddMysqlCluster(&mysqlCluster)
		if !cStatus {
			code = e.ERROR_CREATE_MYSQL
			log.Error("Failed to create mysql cluster!")
		} else {
			result, err := sync_status.K8sCreateMysqlCluster(&mysqlCluster)
			if err != nil {
				log.Error("Failed to create mysqlcluster %v", mysqlCluster.ClusterName)
				code = e.ERROR_CREATE_MYSQL
			} else {
				data["result"] = result
				log.Infof("cluster %v created", clusterName)
				code = e.CREATED
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": data,
	})

}

func UpdateMysqlCluster(c *gin.Context) {

}

func DeleteMysqlCluster(c *gin.Context) {

}
