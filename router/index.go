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

	store := cookie.NewStore([]byte("leo-tracesys"))
	r.Use(sessions.Sessions("SESSION", store))

	r.POST("/api/login", service.Login)
	r.POST("/api/register", service.Register)
	r.GET("/api/logout", middleware.Auth(), service.Logout)

	r.GET("/api/node/list", middleware.Auth(), service.ListNodes)
	r.DELETE("/api/node/:id", middleware.Auth(), service.DelNode)
	r.POST("/api/node/add", middleware.Auth(), service.AddNode)
	r.PUT("/api/node/edit", middleware.Auth(), service.ModifyNode)

	r.POST("/api/target/add", middleware.Auth(), service.AddTarget)
	r.GET("/api/target/list", middleware.Auth(), service.GetTargetList)
	r.DELETE("/api/target/:id", middleware.Auth(), service.DelTarget)
	r.PUT("/api/target/edit", middleware.Auth(), service.ModifyTarget)

	r.GET("/api/user/nodes", middleware.Auth(), service.GetNodesForUser)
	r.GET("/api/user/list", middleware.Auth(), service.ListUsers)

	r.GET("/api/tracelist/token/:token", service.GetTraceList)
	r.POST("/api/result/add", service.RecieveDataFromClient)

	r.POST("/api/result", service.SearchResult)

	return r
}
