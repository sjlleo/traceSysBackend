package service

import (
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/sjlleo/traceSysBackend/models"
)

func GetTaskList(c *gin.Context) {
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

func ModifyTask(c *gin.Context) {
	var t models.Tasks
	u := GetRole(c)
	if err := c.ShouldBindJSON(&t); err == nil {
		if err := u.UpdateTask(&t); err != nil {
			c.JSON(200, gin.H{"code": 500, "error": err.Error()})
			return
		} else {
			c.JSON(200, gin.H{"code": 200, "success": "success"})
		}
	} else {
		c.JSON(200, gin.H{"code": 200, "success": "success"})
	}
}

func DelTask(c *gin.Context) {
	u := GetRole(c)
	id, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if err := u.DeleteTask(uint(id)); err != nil {
		c.JSON(200, gin.H{"code": 500, "error": err.Error()})
		return
	} else {
		c.JSON(200, gin.H{"code": 200, "success": "success"})
	}
}

func AddTask(c *gin.Context) {
	var t models.Tasks
	session := sessions.Default(c)
	id := session.Get("user_id")

	if err := c.ShouldBindJSON(&t); err != nil {
		c.JSON(200, gin.H{"code": 500, "error": err.Error()})
	}
	t.CreatedUserID = id.(uint)
	switch {
	case t.Name == "":
		c.JSON(200, gin.H{"code": 200, "error": "请输入任务名称"})
		return
	case t.Type == 0:
		c.JSON(200, gin.H{"code": 200, "error": "请选择任务类型"})
		return
	case t.CallMethod == 0:
		c.JSON(200, gin.H{"code": 200, "error": "请选择送信方法"})
		return
	case t.TTL == 0:
		c.JSON(200, gin.H{"code": 200, "error": "请选择 TTL"})
		return
	case t.NodeID == 0:
		c.JSON(200, gin.H{"code": 200, "error": "请选择节点"})
		return
	case t.TargetID == 0:
		c.JSON(200, gin.H{"code": 200, "error": "请选择监测目标"})
		return
	case t.ExceedRTT == 0 || t.ExceedPacketLoss == 0:
		c.JSON(200, gin.H{"code": 200, "error": "请选择超时规则"})
		return
	}
	if err := models.CreateTask(&t); err != nil {
		c.JSON(200, gin.H{"code": 500, "error": err.Error()})
		return
	} else {
		c.JSON(200, gin.H{
			"code": 200,
			"msg":  "任务添加成功",
		})
		return
	}
}
