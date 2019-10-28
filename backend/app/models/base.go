package models

import (
	"fmt"
	"time"

	"iris-ticket/backend/app/models/db"
	"iris-ticket/backend/app/models/sys"
	"iris-ticket/backend/config"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/pelletier/go-toml"
)

func InitDB() {
	var gdb *gorm.DB
	var err error
	GetEnv := "test"
	// if GetEnv() == "test" {
	// } else {}

	dbconfig := config.Conf.Get(GetEnv + ".database").(*toml.Tree)
	driver := dbconfig.Get("driver").(string)
	connect := dbconfig.Get("connect").(string)
	gdb, err = gorm.Open(driver, connect)
	if err != nil {
		panic(err)
	}
	gdb.SingularTable(true)
	// gormconfig := config.Conf.Get("gorm").(*toml.Tree)
	gdb.DB().SetMaxIdleConns(150)
	gdb.DB().SetMaxOpenConns(50)
	gdb.DB().SetConnMaxLifetime(time.Duration(7200) * time.Second)
	db.DB = gdb
}

func Migration() {
	fmt.Println(db.DB.AutoMigrate(new(sys.Menu)).Error)
	fmt.Println(db.DB.AutoMigrate(new(sys.User)).Error)
	fmt.Println(db.DB.AutoMigrate(new(sys.RoleMenu)).Error)
	fmt.Println(db.DB.AutoMigrate(new(sys.Role)).Error)
	fmt.Println(db.DB.AutoMigrate(new(sys.UserRole)).Error)
}
