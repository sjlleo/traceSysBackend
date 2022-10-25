package main

import (
	"github.com/sjlleo/traceSysBackend/database"
	"github.com/sjlleo/traceSysBackend/models"
	"github.com/sjlleo/traceSysBackend/router"
)

func initDB() {
	database.InitDBConn()
	models.DBAutoMigration()
}

func main() {
	initDB()
	r := router.New()
	r.Run(":50888")
}
