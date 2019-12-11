package v1

import (
	"baas-service/models"
	"baas-service/pkg/app"
	"baas-service/pkg/e"
	"baas-service/pkg/sync"
	"baas-service/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/common/log"
	"net/http"
)

// @Summary Get a single mysql cluster
// @Produce  json
// @Param id path int true "ID"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/mysql/{id} [get]
func GetMysqlCluster(c *gin.Context) {
	appG := app.Gin{C: c}

	id := utils.Str(c.Param("id")).Int()
	mysqlCluster, err := models.GetMysqlcluster(id)

	if err != nil {
		appG.Response(http.StatusInternalServerError, e.MYSQL_DOES_NOT_EXIST, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, &mysqlCluster)
}

// @Summary Get all mysql cluster
// @Produce json
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/mysqls [get]
func ListMysqlCluster(c *gin.Context) {
	var appG = app.Gin{C: c}

	data := make(map[string]interface{})

	total, err := models.GetArticleTotal()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_MYSQL_CLUSTER_FAIL, nil)
		return
	}

	mysqlClusters, err := models.GetAllMysqlClusters()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_COUNT_MYSQL_CLUSTER_FAIL, nil)
		return
	}

	data["total"] = total
	data["list"] = mysqlClusters

	appG.Response(http.StatusOK, e.SUCCESS, data)
}

// @Summary add mysql cluster
// @Produce json
// @Param 	namespace query string true "Namespace"
// @Param   cluster_name query string true "ClusterName"
// @Param   user query string false "User"
// @Param   password query string true "Password"
// @Param   storage_type query string true "StorageType"
// @Param   multi_master query bool true "multiMaster"
// @Param   version query string false "Version"
// @Param   port query int false "port"
// @Param   volume_size query string false "VolumeSize"
// @Param   default_authentication_plugin query string false "DefaultAuthenticationPlugin"
// @Param   cpu query string false "CPU"
// @Param   memory query string false "Memory"
// @Param   max_connections query int false "MaxConnections"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/mysql [post]
func CreateMysqlCluster(c *gin.Context) {
	var (
		mysqlCluster models.MysqlCluster
		appG         = app.Gin{C: c}
		err          error
		code         int
	)

	err = c.ShouldBind(&mysqlCluster)
	if err != nil {
		log.Error(err)
		appG.Response(http.StatusInternalServerError, e.INVALID_PARAMS, nil)
		return
	}

	mysqlCluster.SetHost()

	if mysqlCluster.Member == 1 {
		mysqlCluster.MultiMaster = false
	}

	exists, existErr := models.ExistMysqlClusterByName(mysqlCluster.ClusterName)
	if existErr != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_CHECK_MYSQL_EXIST_FAIL, nil)
	}

	if exists {
		appG.Response(http.StatusInternalServerError, e.ERROR_EXIST_MYSQL, nil)
		return
	}

	code, err = sync.K8sCreateMysqlCluster(&mysqlCluster)
	if err != nil {
		appG.Response(http.StatusInternalServerError, code, nil)
		return
	}

	mysqlCluster.SetPassword()
	err = models.AddMysqlCluster(&mysqlCluster)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_CREATE_MYSQL, nil)
		return
	}

	appG.Response(http.StatusOK, e.CREATED, nil)
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
	var (
		mysqlCluster *models.MysqlCluster
		appG         = app.Gin{C: c}
		err          error
		code         int
	)

	id := utils.Str(c.Param("id")).Int()

	exists, err := models.ExistMysqlClusterByID(id)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_CHECK_MYSQL_EXIST_FAIL, nil)
		return
	}

	if !exists {
		appG.Response(http.StatusInternalServerError, e.MYSQL_DOES_NOT_EXIST, nil)
		return
	}
	mysqlCluster, _ = models.GetMysqlcluster(id)

	code, err = sync.K8sDeleteMysqlCluster(mysqlCluster)
	if err != nil {
		appG.Response(http.StatusInternalServerError, code, nil)
		return
	}

	err = models.DeleteMysqlCluster(id)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.MYSQL_DELETED_FAILED, nil)
		return
	}

	appG.Response(http.StatusOK, e.MYSQL_DELETED, nil)
}
