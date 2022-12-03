package service

import (
	"errors"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/sjlleo/traceSysBackend/models"
)

// 登录接口
func Login(c *gin.Context) {
	session := sessions.Default(c)

	if id, roleID, err := models.ValidUser(
		c.PostForm("username"),
		c.PostForm("password"),
	); err != nil {
		c.JSON(200, gin.H{
			"code": 401,
			"msg":  "用户名或密码不正确",
		})
	} else {

		user := session.Get("user")
		if user == nil {
			session.Set("user_id", id)
			session.Set("user", c.PostForm("username"))
			session.Set("role", roleID)
			session.Save()
		}

		c.JSON(200, gin.H{
			"code": 200,
			"role": roleID,
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
	var user models.Users
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(200, gin.H{
			"code": 500,
			"msg":  err.Error(),
		})
		return
	}

	user.Role = 2

	if err := models.CreateUser(user); err != nil {
		// 无权限
		c.JSON(200, gin.H{
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

func MyInformation(c *gin.Context) {
	session := sessions.Default(c)
	id := session.Get("user_id").(uint)
	res, _ := models.FindUserByID(id)
	c.JSON(200, res)
}

func CreateUser(c *gin.Context) {
	var user models.Users
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(200, gin.H{
			"code": 500,
			"msg":  err.Error(),
		})
		return
	}

	if err := models.CreateUser(user); err != nil {
		// 无权限
		c.JSON(200, gin.H{
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

func ListUsers(c *gin.Context) {
	p := models.PaginationQ{}
	if err := c.ShouldBind(&p); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
	}
	if err := models.ListUsers(&p); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	} else {
		c.JSON(200, p)
	}
}

func UpdateUser(c *gin.Context) {
	var user models.Users
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(200, gin.H{
			"error": err.Error(),
		})
	} else {
		if err := models.UpdateUser(user); err != nil {
			c.JSON(200, gin.H{
				"error": err.Error(),
			})
		} else {
			c.JSON(200, gin.H{"code": 200, "msg": "success"})
		}
	}
}

func DeleteUser(c *gin.Context) {
	if id_str := c.Param("id"); id_str != "" {
		user := models.Users{}
		id, err := strconv.ParseInt(id_str, 10, 32)
		if err != nil {
			c.JSON(200, gin.H{
				"error": errors.New("ID 不合法").Error(),
			})
		}
		user.ID = uint(id)
		// log.Println(user)
		if err := models.DelUser(user); err != nil {
			c.JSON(200, gin.H{
				"error": err.Error,
			})
		} else {
			c.JSON(200, gin.H{"code": 200, "msg": "success"})
		}
	} else {
		c.JSON(200, gin.H{
			"error": errors.New("ID 不合法").Error(),
		})
	}
}

func UpdatePwd(c *gin.Context) {
	var pwd models.Password
	if err := c.ShouldBindJSON(&pwd); err != nil {
		c.JSON(200, gin.H{
			"err": err.Error(),
		})
		return
	}
	session := sessions.Default(c)
	pwd.UserID = uint(session.Get("user_id").(int))
	if err := models.UpdatePassword(&pwd); err != nil {
		c.JSON(200, gin.H{
			"code": 500,
			"msg":  err.Error(),
		})
		return
	} else {
		c.JSON(200, gin.H{
			"code": 200,
			"msg":  "修改成功",
		})
	}
}
