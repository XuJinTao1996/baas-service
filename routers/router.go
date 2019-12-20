package routers

import (
	_ "baas-service/docs"
	"baas-service/pkg/middleware"
	"baas-service/pkg/setting"
	v1 "baas-service/routers/api/v1"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func InitRouter() *gin.Engine {

	gin.SetMode(setting.RunMode)

	r := gin.New()

	r.Use(gin.Recovery())

	r.Use(middleware.Logger())

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	apiv1 := r.Group("/api/v1")
	{
		apiv1.GET("/mysqls", v1.ListMysqlCluster)
		apiv1.GET("/mysql/:id", v1.GetMysqlCluster)
		apiv1.POST("/mysql", v1.CreateMysqlCluster)
		//apiv1.PUT("/mysql/:id", v1.UpdateMysqlCluster)
		apiv1.DELETE("/mysql/:id", v1.DeleteMysqlCluster)

		apiv1.GET("/kafkas", v1.ListKafkaCluster)
		apiv1.POST("/kafka", v1.CreateKafkaCluster)
		apiv1.GET("/kafka/:id", v1.GetKafkaCluster)
		apiv1.DELETE("/kafka/:id", v1.DeleteKafkaCluster)

		apiv1.POST("/zookeeper", v1.CreateZookeeperCluster)
		apiv1.GET("/zookeeper/:id", v1.GetZookeeperCluster)
		apiv1.GET("/zookeepers", v1.ListZookeeperCluster)
		apiv1.DELETE("/zookeeper/:id", v1.DeleteZookeeperCluster)

		apiv1.POST("/kafka_topic", v1.CreateKafkaTopic)
		apiv1.GET("/kafka_topic/:id", v1.GetKafkaTopic)
		apiv1.GET("/kafka_topics", v1.ListKafkaTopic)
		apiv1.DELETE("/kafka_topic/:id", v1.DeleteKafkaTopic)
	}

	return r
}
