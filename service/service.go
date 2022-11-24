package service

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/sjlleo/traceSysBackend/models"
)

func GetRole(c *gin.Context) models.User {
	var u models.User

	session := sessions.Default(c)
	roleID := session.Get("role").(uint)
	userID := session.Get("user_id").(uint)
	if roleID == 1 {
		u = &models.Admin{
			RoleID: roleID,
			UserID: userID,
		}
	} else {
		u = &models.Normal{
			RoleID: roleID,
			UserID: userID,
		}
	}
	return u
}
