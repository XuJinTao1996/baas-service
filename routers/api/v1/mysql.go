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

// @Summary Get a single mysql cluster
// @Produce  json
// @Param id path int true "ID"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/mysql/{id} [get]
func GetMysqlCluster(c *gin.Context) {
	appG := app.Gin{C: c}

	id := utils.Str(c.Param("id")).Int()
	mysqlCluster, state := models.GetMysqlcluster(id)

	if !state {
		log.Errorf("cluster does not exist!")
		appG.Response(http.StatusInternalServerError, e.MYSQL_DOES_NOT_EXIST, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, mysqlCluster)
}

// @Summary Get all mysql cluster
// @Produce json
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/mysqls [get]
func ListMysqlCluster(c *gin.Context) {
	appG := app.Gin{C: c}

	data := make(map[string]interface{})

	data["list"], data["total"] = models.GetAllMysqlClusters()

	appG.Response(http.StatusOK, e.SUCCESS, data)
}

// @Summary add mysql cluster
// @Produce json
// @Param 	namespace query string true "Namespace"
// @Param   cluster_name query string true "ClusterName"
// @Param   user query string false "user"
// @Param   password query string true "password"
// @Param   storage_type query string true "storage_type"
// @Param   multi_master query bool true "multi_master"
// @Param   version query string false "version"
// @Param   port query int false "port"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/mysql [post]
func CreateMysqlCluster(c *gin.Context) {
	var (
		mysqlCluster models.MysqlCluster
		appG         = app.Gin{C: c}
	)

	err := c.ShouldBind(&mysqlCluster)
	if err != nil {
		log.Error(err)
		appG.Response(http.StatusInternalServerError, e.INVALID_PARAMS, nil)
		return
	}

	mysqlCluster.SetHost()
	mysqlCluster.SetPassword()

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

	appG.Response(http.StatusOK, e.CREATED, newMysqlCluster)

}

// Todo
func UpdateMysqlCluster(c *gin.Context) {

}

// @Summary Delete mysql cluster
// @Produce  json
// @Param id path int true "ID"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/mysql/{id} [delete]
func DeleteMysqlCluster(c *gin.Context) {
	var appG = app.Gin{C: c}
	id := utils.Str(c.Param("id")).Int()
	mysqlCluster, state := models.GetMysqlcluster(id)

	if !state {
		log.Error("mysql cluster does not existed")
		appG.Response(http.StatusInternalServerError, e.MYSQL_DOES_NOT_EXIST, nil)
	}

	sync.K8sDeleteMysqlCluster(&mysqlCluster)
	models.DeleteMysqlCluster(id)

	appG.Response(http.StatusOK, e.MYSQL_DELETED, nil)

}
