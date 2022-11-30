package service

import (
	"errors"
	"log"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/sjlleo/traceSysBackend/models"
)

func GetTargetList(c *gin.Context) {
	u := GetRole(c)

	p := models.PaginationQ{}
	if err := c.ShouldBind(&p); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
	}
	if err := u.ListTargets(&p); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	} else {
		c.JSON(200, p)
	}
}

func GetTargetUser(c *gin.Context) {
	targetID := c.Param("targetID")
	tid, err := strconv.Atoi(targetID)
	if err != nil {
		c.JSON(200, gin.H{"error": err.Error()})
		return
	}
	u := GetRole(c)
	res, err := u.ListTargetUser(uint(tid))
	if err != nil {
		c.JSON(200, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, res)
}

func TakeTargetNodeInfo(c *gin.Context) {
	var nodesArr models.NodeInfo
	u := GetRole(c)
	nodesArr.TargetIP = c.Param("ip")
	if nodesArr.TargetIP != "" {
		if err := u.FindTargetIPNodeInfo(&nodesArr); err != nil {
			c.JSON(200, gin.H{"code": 500, "error": err.Error()})
			return
		} else {
			nodesArr.Code = 200
			c.JSON(200, nodesArr)
		}
	} else {
		c.JSON(200, gin.H{"code": 500, "error": errors.New("IP param not found")})
	}
}

func ModifyTarget(c *gin.Context) {
	var t models.Target
	u := GetRole(c)
	if err := c.ShouldBindJSON(&t); err == nil {
		log.Println(t)
		if err := u.ModifyTarget(&t); err != nil {
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
	u := GetRole(c)
	id, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if err := u.DelTarget(int(id)); err != nil {
		c.JSON(200, gin.H{"code": 500, "error": err.Error()})
		return
	} else {
		c.JSON(200, gin.H{"code": 200, "success": "success"})
	}
}

func AddTarget(c *gin.Context) {
	var t models.Target
	session := sessions.Default(c)
	id := session.Get("user_id")

	if err := c.ShouldBindJSON(&t); err != nil {
		c.JSON(200, gin.H{"code": 500, "error": err.Error()})
	}
	t.CreatedUserID = id.(uint)
	if err := models.AddTarget(t); err != nil {
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
