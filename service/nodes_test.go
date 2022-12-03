package service

import (
	"fmt"
	"log"
	"testing"

	"github.com/sjlleo/traceSysBackend/database"
	"github.com/sjlleo/traceSysBackend/models"
)

func TestNodeStr(t *testing.T) {
	database.Init()

	token := "cuxfDeNdaB4TZb8dDyKBD"
	res, err := models.GetNodeFromToken(token)
	if err != nil {
		// c.String(400, "Token Invalid")
		return
	}
	user, _ := models.FindUserByID(res.CreatedUserID)
	str := ""
	str += fmt.Sprintf("%v|",user.Username)
	str += fmt.Sprintf("%v|",res.IP)
	if res.Role == 1 {
		str += fmt.Sprintf("%v|", "管理员")
	} else {
		str += fmt.Sprintf("%v|", "用户")
	}
	log.Println(str)
}
