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

func CreateKafkaTopic(c *gin.Context) {
	var kafkaTopic *models.KafkaTopic
	var appG = app.Gin{C: c}

	err := c.ShouldBind(&kafkaTopic)

	if err != nil {
		appG.Response(http.StatusInternalServerError, e.INVALID_PARAMS, nil)
		return
	}

	exist, err := models.ExistKafkaClusterByName(kafkaTopic.ClusterRefName)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_CHECK_KAFKA_EXIST_FAIL, nil)
		return
	}

	if !exist {
		appG.Response(http.StatusInternalServerError, e.KAFKA_TOPIC_REF_CLUSTER_DOES_NOT_EXIST, nil)
		return
	}

	code, err := sync.K8sCreateKafkaTopic(kafkaTopic)

	if err != nil {
		appG.Response(http.StatusInternalServerError, code, nil)
		return
	}

	err = models.AddKafkaTopic(kafkaTopic)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_CREATE_KAFKA_TOPIC, nil)
		return
	}
	appG.Response(http.StatusOK, e.CREATED, nil)
}

func GetKafkaTopic(c *gin.Context) {
	var appG = app.Gin{C: c}
	id := utils.Str(c.Param("id")).Int()
	kafkaTopic, err := models.GetKafkaTopicByID(id)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.KAFKA_TOPIC_DOES_NOT_EXIST, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, kafkaTopic)
}

func ListKafkaTopic(c *gin.Context) {
	var appG = app.Gin{C: c}

	data := make(map[string]interface{})

	total, err := models.GetKafkaTopicTotal()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.KAFKA_TOPIC_DOES_NOT_EXIST, nil)
		return
	}

	kafkaTopics, err := models.GetAllKafkaTopic()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.KAFKA_TOPIC_DOES_NOT_EXIST, nil)
		return
	}

	data["list"] = kafkaTopics
	data["total"] = total

	appG.Response(http.StatusOK, e.SUCCESS, data)
}

func DeleteKafkaTopic(c *gin.Context) {
	var (
		kafkaTopic *models.KafkaTopic
		appG       = app.Gin{C: c}
		err        error
		code       int
	)

	id := utils.Str(c.Param("id")).Int()

	exists, err := models.ExistKafkaTopicByID(id)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_CHECK_KAFKA_TOPIC_EXIST_FAIL, nil)
		return
	}

	if !exists {
		appG.Response(http.StatusInternalServerError, e.KAFKA_TOPIC_DOES_NOT_EXIST, nil)
		return
	}

	code, err = sync.K8sDeleteKafkaTopic(kafkaTopic)
	if err != nil {
		appG.Response(http.StatusInternalServerError, code, nil)
		return
	}

	err = models.DeleteKafkaTopic(id)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.KAFKA_TOPIC_DELETED_FAILED, nil)
		return
	}

	appG.Response(http.StatusOK, e.KAFKA_TOPIC_DELETED, nil)
}
