package server

import (
	"fmt"
	"github.com/Anniext/Arkitektur/server/middlewares"
	"github.com/Anniext/Arkitektur/system/log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func InitDefaultGin(defaultRegister func(*gin.RouterGroup)) error {
	defaultGin = gin.Default()

	defaultGin.Use(middlewares.CorsHandler)

	defaultRegister(defaultGin.Group("/api")) // 注入路由

	cnf := GetDefaultGinConfig()
	addr := fmt.Sprintf("%s:%d", cnf.Addr, cnf.Port)

	defaultServer = &http.Server{
		Addr:    addr,
		Handler: defaultGin,
	}

	go SafeGoRecoverWarpFunc(func() {
		if err := defaultServer.ListenAndServe(); err != nil {
			log.Error("Gin server start err: ", err)
		}
	})()
	return nil
}

var (
	defaultGin    *gin.Engine
	defaultServer *http.Server
)

func GetDefaultGin() *gin.Engine {
	return defaultGin
}

func GetDefaultServer() *http.Server {
	return defaultServer
}
