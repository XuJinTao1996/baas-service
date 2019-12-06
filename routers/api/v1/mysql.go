package v1

import (
	"baas-service/models"
	"baas-service/pkg/app"
	"baas-service/pkg/e"
	"baas-service/pkg/informer"
	"baas-service/pkg/k8s/client"
	"baas-service/pkg/sync"
	"baas-service/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/common/log"
	"net/http"
)

func init() {
	startDeploymentInformer()
	startMysqlClusterInformer()
}

// 启动 deployment informer
func startDeploymentInformer() {
	informer.DeploymentInformer(client.K8sClient)
}

// 启动 mysql cluster informer
func startMysqlClusterInformer() {
	informer.MysqlClusterInformer(client.MysqlClientset)
}

// 获取指定 mysql 实例
func GetMysqlCluster(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		code = e.SUCCESS
	)

	id := utils.Str(c.Param("id")).Int()
	mysqlCluster, state := models.GetMysqlcluster(id)

	if !state {
		log.Errorf("cluster does not exist!")
		appG.Response(http.StatusInternalServerError, e.MYSQL_DOES_NOT_EXIST, nil)
		return
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

// @Summary 新增 mysql cluster 实例
// @Produce json
// @Param 	namespace query string true "Namespace"
// @Param   cluster_name query string true "ClusterName"
// @Param   user query string false "user"
// @Param   password query string true "password"
// @Param   storage_type query string true "storage_type"
// @Param   multi_master query bool true "multi_master"
// @Param   version query string false "version"
// @Param   port query int false "port"
// @Success 200 {object} models.MysqlCluster
// @Failure 500 {string} json "{"code":500,"data":nil,"msg":"mysql 创建失败"}"
// @Router /api/v1/mysql [post]
func CreateMysqlCluster(c *gin.Context) {
	var (
		mysqlCluster models.MysqlCluster
		appG         = app.Gin{C: c}
		code         = e.CREATED
	)

	err := c.ShouldBind(&mysqlCluster)
	if err != nil {
		log.Error(err)
	}

	mysqlCluster.Host = mysqlCluster.RouterDeploymentName()

	if mysqlCluster.Member == 1 {
		mysqlCluster.MultiMaster = false
	}

	isExisted := models.ExistMysqlCluster(mysqlCluster.ClusterName)
	if isExisted {
		appG.Response(http.StatusInternalServerError, e.ERROR_EXIST_MYSQL, nil)
		return
	}

	newMysqlCluster, cStatus := models.AddMysqlCluster(&mysqlCluster)
	if !cStatus {
		log.Error("Failed to create mysql cluster!")
		appG.Response(http.StatusInternalServerError, e.ERROR_CREATE_MYSQL, nil)
		return
	}

	_, err = sync.K8sCreateMysqlCluster(&mysqlCluster)
	if err != nil {
		log.Error("Failed to create mysqlcluster %v", mysqlCluster.ClusterName)
		appG.Response(http.StatusInternalServerError, e.INVALID_PARAMS, nil)
		return
	}

	log.Infof("cluster %v created", mysqlCluster.ClusterName)

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": newMysqlCluster,
	})

}

// Todo
func UpdateMysqlCluster(c *gin.Context) {

}

func DeleteMysqlCluster(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		code = e.MYSQL_DELETED
	)
	id := utils.Str(c.Param("id")).Int()
	mysqlCluster, state := models.GetMysqlcluster(id)

	if !state {
		log.Error("mysql cluster does not existed")
		appG.Response(http.StatusInternalServerError, e.MYSQL_DOES_NOT_EXIST, nil)
	}

	sync.K8sDeleteMysqlCluster(&mysqlCluster)
	models.DeleteMysqlCluster(id)

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": nil,
	})
}
