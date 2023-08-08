package orm

import (
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm" //สร้างตาราง
)

var Db *gorm.DB
var err error

func InitDB() {
	dsn := os.Getenv("MYSQL_DNS")	//ต่อ database
  Db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})		// ต่อ database
  if err != nil { 	//กำหนด error
	  panic("failed to connect database")
  }		

  // Migrate the schema
  Db.AutoMigrate(&User{})
}