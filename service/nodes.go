package service

import (
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/sjlleo/traceSysBackend/models"
)

func ListNodes(c *gin.Context) {
	p := models.PaginationQ{}
	if err := c.ShouldBind(&p); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
	}
	if err := models.ListNodes(&p); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	} else {
		c.JSON(200, p)
	}
}

func GetNodesForUser(c *gin.Context) {
	session := sessions.Default(c)
	role := session.Get("role")
	res, _ := models.ListNodesUser(role.(int))
	c.JSON(200, gin.H{"code": 200, "data": res})
}

func DelNode(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if err := models.DelNode(int(id)); err != nil {
		c.JSON(200, gin.H{"code": 500, "error": err.Error()})
		return
	} else {
		c.JSON(200, gin.H{"code": 200, "success": "success"})
	}
}

func AddNode(c *gin.Context) {
	if err := models.AddNode(
		c.PostForm("ip"),
		c.PostForm("role"),
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
	if err := c.ShouldBind(&node); err == nil {
		if err := models.ModifyNode(&node); err != nil {
			c.JSON(200, gin.H{"code": 500, "error": err.Error()})
			return
		} else {
			c.JSON(200, gin.H{"code": 200, "success": "success"})
		}
	} else {
		c.JSON(200, gin.H{"code": 200, "success": "success"})
	}
}
