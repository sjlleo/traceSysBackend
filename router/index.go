package router

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/sjlleo/traceSysBackend/middleware"
	"github.com/sjlleo/traceSysBackend/service"
)

func New() *gin.Engine {
	r := gin.Default()

	store := cookie.NewStore([]byte("dogdai"))
	r.Use(sessions.Sessions("SESSION", store))

	r.POST("/login", service.Login)
	r.POST("/register", service.Register)
	r.GET("/logout", middleware.Auth(), service.Logout)

	r.GET("/ping", middleware.Auth(), service.Ping)
	r.GET("/node/list", middleware.Auth(), service.ListNodes)
	r.POST("/node/add", middleware.Auth(), service.AddNode)
	return r
}
