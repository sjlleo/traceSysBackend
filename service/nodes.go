package service

import (
	"github.com/gin-gonic/gin"
	"github.com/sjlleo/traceSysBackend/models"
)

func ListNodes(c *gin.Context) {
	if res, err := models.ListNodes(); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	} else {
		c.JSON(200, res)
	}
}

func AddNode(c *gin.Context) {
	if err := models.AddNode(
		c.PostForm("ip"),
		c.PostForm("role"),
		c.PostForm("scret"),
	); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	} else {
		c.JSON(200, gin.H{
			"code": 200,
			"msg":  "节点添加成功",
		})
		return
	}
}
