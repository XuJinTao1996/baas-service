package main

import (
	"fmt"
	"net/http"
	"rds-front/pkg/setting"
	"rds-front/routers"
)

func main() {
	router := routers.InitRouter()

	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", setting.HttpPort),
		Handler:        router,
		ReadTimeout:    setting.ReadTimeOut,
		WriteTimeout:   setting.WriteTimeOut,
		MaxHeaderBytes: 1 << 20,
	}

	s.ListenAndServe()
}
