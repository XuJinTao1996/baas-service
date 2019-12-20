package v1

import (
	"baas-service/models"
	"baas-service/pkg/app"
	"baas-service/pkg/e"
	"baas-service/pkg/sync"
	"baas-service/pkg/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CreateZookeeperCluster(c *gin.Context) {
	var zookeeperCluster *models.ZookeeperCluster
	var appG = app.Gin{C: c}
	err := c.ShouldBind(&zookeeperCluster)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.INVALID_PARAMS, nil)
		return
	}
	code, err := sync.K8sCreateZookeeperCluster(zookeeperCluster)
	if err != nil {
		appG.Response(http.StatusInternalServerError, code, nil)
		return
	}
	err = models.AddZookeeperCluster(zookeeperCluster)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_CREATE_ZOOKEEPER, nil)
		return
	}
	appG.Response(http.StatusOK, e.CREATED, nil)
}

func GetZookeeperCluster(c *gin.Context) {
	var appG = app.Gin{C: c}
	id := utils.Str(c.Param("id")).Int()
	zookeeperCluster, err := models.GetZookeeperClusterByID(id)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ZOOKEEPER_DOES_NOT_EXIST, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, zookeeperCluster)
}

func ListZookeeperCluster(c *gin.Context) {
	var appG = app.Gin{C: c}

	data := make(map[string]interface{})

	total, err := models.GetZookeeperClusterTotal()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_ZOOKEEPER_CLUSTER_FAIL, nil)
		return
	}

	kafkaClusters, err := models.GetAllZookeeperCluster()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_COUNT_ZOOKEEPER_CLUSTER_FAIL, nil)
		return
	}

	data["list"] = kafkaClusters
	data["total"] = total

	appG.Response(http.StatusOK, e.SUCCESS, data)
}

func DeleteZookeeperCluster(c *gin.Context) {
	var (
		zookeeperCluster *models.ZookeeperCluster
		appG             = app.Gin{C: c}
		err              error
		code             int
	)

	id := utils.Str(c.Param("id")).Int()

	exists, err := models.ExistZookeeperClusterByID(id)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_CHECK_ZOOKEEPER_EXIST_FAIL, nil)
		return
	}

	if !exists {
		appG.Response(http.StatusInternalServerError, e.ZOOKEEPER_DOES_NOT_EXIST, nil)
		return
	}

	status, err := models.CheckZookeeperClusterUsageStatus(id)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.CHECK_ZOOKEEPER_USAGE_STATUS_FAILED, nil)
		return
	}

	if status {
		appG.Response(http.StatusInternalServerError, e.KAFLA_CLUSTER_IS_USING_THE_CURRENT_ZOOKEEPER_CLUSTER, nil)
		return
	}

	zookeeperCluster, _ = models.GetZookeeperClusterByID(id)

	code, err = sync.K8sDeleteZookeeperCluster(zookeeperCluster)
	if err != nil {
		appG.Response(http.StatusInternalServerError, code, nil)
		return
	}

	err = models.DeleteZookeeperCluster(id)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ZOOKEEPER_DELETED_FAILED, nil)
		return
	}

	appG.Response(http.StatusOK, e.ZOOKEEPER_DELETED, nil)
}
