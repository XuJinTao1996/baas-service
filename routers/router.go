package routers

import (
	"github.com/gin-gonic/gin"
	v1 "rds-front/pkg/controllers/v1"
	"rds-front/pkg/setting"
)

func InitRouter() *gin.Engine {
	r := gin.New()

	r.Use(gin.Logger())

	r.Use(gin.Recovery())

	gin.SetMode(setting.RunMode)

	apiv1 := r.Group("/api/v1")
	{

		// 获取 mysql 实例列表
		apiv1.GET("/mysql", v1.GetMysqlInstances)
		// 新建 mysql 实例
		apiv1.POST("/mysql", v1.CreateMysqlInstance)
		// 更新指定 mysql 实例
		apiv1.PUT("/mysql/:id", v1.UpdateMysqlInstance)
		// 删除指定 mysql 实例
		apiv1.DELETE("/mysql/:id", v1.DeleteMysqlInstance)

	}

	return r
}
