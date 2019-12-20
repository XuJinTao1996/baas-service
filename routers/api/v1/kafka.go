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

func CreateKafkaCluster(c *gin.Context) {
	var kafkaCluster *models.KafkaCluster
	var appG = app.Gin{C: c}

	err := c.ShouldBind(&kafkaCluster)

	if err != nil {
		appG.Response(http.StatusInternalServerError, e.INVALID_PARAMS, nil)
		return
	}

	code, err := sync.K8sCreateKafkaCluster(kafkaCluster)

	if err != nil {
		appG.Response(http.StatusInternalServerError, code, nil)
		return
	}

	err = models.AddKafkaCluster(kafkaCluster)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_CREATE_KAFKA, nil)
		return
	}
	appG.Response(http.StatusOK, e.CREATED, nil)
}

func GetKafkaCluster(c *gin.Context) {
	var appG = app.Gin{C: c}
	id := utils.Str(c.Param("id")).Int()
	kafkaCluster, err := models.GetKafkaClusterByID(id)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.KAFKA_DOES_NOT_EXIST, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, kafkaCluster)
}

func ListKafkaCluster(c *gin.Context) {
	var appG = app.Gin{C: c}

	data := make(map[string]interface{})

	total, err := models.GetKafkaClusterTotal()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_KAFKA_CLUSTER_FAIL, nil)
		return
	}

	kafkaClusters, err := models.GetAllKafkaCluster()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_COUNT_KAFKA_CLUSTER_FAIL, nil)
		return
	}

	data["list"] = kafkaClusters
	data["total"] = total

	appG.Response(http.StatusOK, e.SUCCESS, data)
}

func DeleteKafkaCluster(c *gin.Context) {
	var (
		kafkaCluster *models.KafkaCluster
		appG         = app.Gin{C: c}
		err          error
		code         int
	)

	id := utils.Str(c.Param("id")).Int()

	exists, err := models.ExistKafkaClusterByID(id)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_CHECK_KAFKA_EXIST_FAIL, nil)
		return
	}

	if !exists {
		appG.Response(http.StatusInternalServerError, e.KAFKA_DOES_NOT_EXIST, nil)
		return
	}

	kafkaCluster, _ = models.GetKafkaClusterByID(id)

	code, err = sync.K8sDeleteKafkaCluster(kafkaCluster)
	if err != nil {
		appG.Response(http.StatusInternalServerError, code, nil)
		return
	}

	err = models.DeleteKafkaCluster(id)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.KAFKA_DELETED_FAILED, nil)
		return
	}

	appG.Response(http.StatusOK, e.KAFKA_DELETED, nil)
}
