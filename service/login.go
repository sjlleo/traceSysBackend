package service

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/sjlleo/traceSysBackend/models"
)

// 登录接口
func Login(c *gin.Context) {
	session := sessions.Default(c)

	if roleID, err := models.ValidUser(
		c.PostForm("username"),
		c.PostForm("password"),
	); err != nil {
		c.JSON(401, gin.H{
			"code": 401,
			"msg":  "用户名或密码不正确",
		})
	} else {

		user := session.Get("user")
		if user == nil {
			session.Set("user", c.PostForm("username"))
			session.Set("role", roleID)
			session.Save()
		}

		c.JSON(200, gin.H{
			"code": 200,
			"msg":  "登录成功",
		})
	}
}

func Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "登出成功",
	})
}

// 注册接口
func Register(c *gin.Context) {
	if err := models.CreateUser(
		c.PostForm("username"),
		c.PostForm("password"),
		2,
	); err != nil {
		// 无权限
		c.JSON(403, gin.H{
			"code": 403,
			"msg":  err.Error(),
		})
	} else {
		c.JSON(200, gin.H{
			"code": 200,
			"msg":  "注册成功",
		})
	}
}
