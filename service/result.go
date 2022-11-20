package service

import (
	"github.com/gin-gonic/gin"
	"github.com/sjlleo/traceSysBackend/models"
)

func RecieveDataFromClient(c *gin.Context) {
	var clientData models.ClientData
	if err := c.ShouldBindJSON(&clientData); err != nil {
		c.JSON(500, err.Error())
		return
	} else {
		if err := models.AddTraceData(&clientData); err != nil {
			c.JSON(500, err.Error())
		} else {
			c.JSON(200, gin.H{
				"code": 200,
				"msg":  "success",
			})
		}
	}
}
