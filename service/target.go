package service

import (
	"log"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/sjlleo/traceSysBackend/models"
)

func GetTargetList(c *gin.Context) {
	p := models.PaginationQ{}
	if err := c.ShouldBind(&p); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
	}
	if err := models.ListTargets(&p); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	} else {
		c.JSON(200, p)
	}
}

func ModifyTarget(c *gin.Context) {
	var t models.Target
	if err := c.Bind(&t); err == nil {
		log.Println(t)
		if err := models.ModifyTarget(&t); err != nil {
			c.JSON(200, gin.H{"code": 500, "error": err.Error()})
			return
		} else {
			c.JSON(200, gin.H{"code": 200, "success": "success"})
		}
	} else {
		c.JSON(200, gin.H{"code": 200, "success": "success"})
	}
}

func DelTarget(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if err := models.DelTarget(int(id)); err != nil {
		c.JSON(200, gin.H{"code": 500, "error": err.Error()})
		return
	} else {
		c.JSON(200, gin.H{"code": 200, "success": "success"})
	}
}

func AddTarget(c *gin.Context) {
	session := sessions.Default(c)
	id := session.Get("user_id")
	interval, err := strconv.Atoi(c.PostForm("interval"))
	if err!= nil {
		c.JSON(200, gin.H{"code": 500, "error": err.Error()})
	}
	method, err := strconv.Atoi(c.PostForm("method"))
	if err!= nil {
		c.JSON(200, gin.H{"code": 500, "error": err.Error()})
	}
	if err = models.AddTarget(
		c.PostForm("ip"),
		id.(int),
		interval,
		method,
		c.PostForm("nodeid"),
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
