package middleware

import (
	"log"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		log.Println(session.Get("user"))
		if session.Get("user") != nil {
			c.Next()
			return
		}
		c.JSON(401, gin.H{
			"code": 401,
			"msg":  "无权限",
		})
		c.Abort()
	}
}
