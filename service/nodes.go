package service

import (
	"fmt"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/sjlleo/traceSysBackend/models"
)

func ListNodes(c *gin.Context) {
	u := GetRole(c)

	p := models.PaginationQ{}
	if err := c.ShouldBind(&p); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
	}
	if err := u.ListNodes(&p); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	} else {
		c.JSON(200, p)
	}
}

func GetNodeInfoForShell(c *gin.Context) {
	token := c.Param("token")
	res, err := models.GetNodeFromToken(token)
	if err != nil {
		c.String(400, "Token Invalid")
		return
	}
	user, _ := models.FindUserByID(res.CreatedUserID)
	str := ""
	str += fmt.Sprintf("%v|", user.Username)
	str += fmt.Sprintf("%v|", res.IP)
	if user.Role == 1 {
		str += fmt.Sprintf("%v", "管理员")
	} else {
		str += fmt.Sprintf("%v", "用户")
	}
	c.String(200, str)
}

func GetNodesForUser(c *gin.Context) {
	session := sessions.Default(c)
	role := session.Get("role")
	res, _ := models.ListNodesUser(role.(uint))
	c.JSON(200, gin.H{"code": 200, "data": res})
}

func DelNode(c *gin.Context) {
	u := GetRole(c)

	id, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if err := u.DelNode(int(id)); err != nil {
		c.JSON(200, gin.H{"code": 500, "error": err.Error()})
		return
	} else {
		c.JSON(200, gin.H{"code": 200, "success": "success"})
	}
}

func AddNode(c *gin.Context) {
	u := GetRole(c)

	if err := u.AddNode(
		c.PostForm("ip"),
		c.PostForm("role"),
		c.PostForm("alias"),
		c.PostForm("secret"),
	); err != nil {
		c.JSON(200, gin.H{"code": 500, "error": err.Error()})
		return
	} else {
		c.JSON(200, gin.H{
			"code": 200,
			"msg":  "节点添加成功",
		})
		return
	}
}

func ModifyNode(c *gin.Context) {
	var node models.Nodes
	u := GetRole(c)

	if err := c.ShouldBind(&node); err == nil {
		if err := u.ModifyNode(&node); err != nil {
			c.JSON(200, gin.H{"code": 500, "error": err.Error()})
			return
		} else {
			c.JSON(200, gin.H{"code": 200, "success": "success"})
		}
	} else {
		c.JSON(200, gin.H{"code": 500, "error": err.Error()})
	}
}
