package sqlite

import (
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Heartbeat struct {
	gorm.Model
	Ip               string
	ReportedAt       time.Time
	ConfigVer        string
	TargetConfigVer  string
	TargetConfig     string
	RuntimeConfigVer string
	RuntimConfig     string
}

func CreateTable() {
	// Set the password for the SQLite database
	dsn := "file:test.db?_auth&_auth_user=admin&_auth_pass=yourpassword"

	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	db.AutoMigrate(&Heartbeat{})
}

/*
func (db *gorm.DB) CreateHeartbeat() {
	// Create
	db.Create(&Heartbeat{Ip: "", ReportedAt: time.Now()})

}
*/
