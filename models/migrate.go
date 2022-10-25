package models

import (
	"github.com/sjlleo/traceSysBackend/database"
)

func DBAutoMigration() {
	db := database.GetDB()
	db.AutoMigrate(
		&Users{},
		&Nodes{},
	)
}
