package routers

import (
	v1 "baas-service/pkg/controllers/v1"
	"baas-service/pkg/setting"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.New()

	r.Use(gin.Logger())

	r.Use(gin.Recovery())

	gin.SetMode(setting.RunMode)

	apiv1 := r.Group("/api/v1")
	{

		// 获取 mysql 实例列表
		apiv1.GET("/mysql", v1.ListMysqlCluster)
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
