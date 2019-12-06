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

		// 获取 mysql 实例列表
		apiv1.GET("/mysqls", v1.ListMysqlCluster)
		// 获取指定 mysql 实例
		apiv1.GET("/mysql/:id", v1.GetMysqlCluster)
		// 新建 mysql 实例
		apiv1.POST("/mysql", v1.CreateMysqlCluster)
		// 更新指定 mysql 实例
		apiv1.PUT("/mysql/:id", v1.UpdateMysqlCluster)
		// 删除指定 mysql 实例
		apiv1.DELETE("/mysql/:id", v1.DeleteMysqlCluster)

	}

	return r
}
