package db

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"log"
	"time"
)

var db *gorm.DB

const (
	DbUserName     = "root"
	DbPassword     = "4524"
	DbAddress      = "127.0.0.1:3306"
	DbName         = "gqbot"
	DBMaxOpenConns = 100
	DBMaxIdleConns = 15
	DBMaxLifeTime  = 200
)

//DbInit 数据库初始化
func DbInit() error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		DbUserName, DbPassword, DbAddress, DbName)
	var err error
	db, err = gorm.Open(mysql.New(mysql.Config{
		DSN:                     dsn,
		DontSupportRenameColumn: true,
	}), &gorm.Config{
		SkipDefaultTransaction: false,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})

	if err != nil {
		log.Println("gorm open failed", err)
		return err
	}

	sqlDb, err := db.DB()
	if err != nil {
		log.Println("db.DB failed", err)
		return err
	}
	sqlDb.SetMaxOpenConns(DBMaxOpenConns)
	sqlDb.SetMaxIdleConns(DBMaxIdleConns)
	sqlDb.SetConnMaxLifetime(time.Second * DBMaxLifeTime)

	return nil
}
