package taskgroup

import (
	"testing"
)

func initDB() {
	// database.InitDBConn()
	// models.DBAutoMigration()
}

func TestTask(t *testing.T) {
	initDB()
	StartTaskCycle()
}
