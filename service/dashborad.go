package service

import (
	"github.com/gin-gonic/gin"
)

func StatisticsData(c *gin.Context) {
	u := GetRole(c)

	nodeCount, err := u.CountNode()
	if err != nil {
		c.JSON(200, gin.H{
			"code":  500,
			"error": err.Error(),
		})
		return
	}
	targetCount, err := u.CountTarget()
	if err != nil {
		c.JSON(200, gin.H{
			"code":  500,
			"error": err.Error(),
		})
		return
	}
	userCount, err := u.CountUser()
	if err != nil {
		c.JSON(200, gin.H{
			"code":  500,
			"error": err.Error(),
		})
		return
	}
	taskCount, err := u.CountTask()
	if err != nil {
		c.JSON(200, gin.H{
			"code":  500,
			"error": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"code":   200,
		"node":   nodeCount,
		"target": targetCount,
		"user":   userCount,
		"task":   taskCount,
	})

}
