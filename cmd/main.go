package main

import (
	"baas-service/pkg/informer"
	"baas-service/pkg/k8s/client"
	"baas-service/pkg/setting"
	"baas-service/routers"
	"fmt"
	"net/http"
)

func main() {
	router := routers.InitRouter()

	informer.DeploymentInformer(client.K8sClient)
	informer.MysqlClusterInformer(client.MysqlClientset)

	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", setting.HttpPort),
		Handler:        router,
		ReadTimeout:    setting.ReadTimeOut,
		WriteTimeout:   setting.WriteTimeOut,
		MaxHeaderBytes: 1 << 20,
	}

	s.ListenAndServe()
}
