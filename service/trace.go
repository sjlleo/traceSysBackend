package service

import (
	"github.com/gin-gonic/gin"
	"github.com/sjlleo/traceSysBackend/models"
)

func GetTraceList(c *gin.Context) {
	token := c.Param("token")
	list := models.GetTraceList(token)
	// result, _ := json.Marshal(list)
	c.JSON(200, list)
}
