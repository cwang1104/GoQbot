package db

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"log"
	"qbot/utils"
	"time"
)

var db *gorm.DB

//DbInit 数据库初始化
func DbInit() error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		utils.GlobalConf.Mysql.DBUserName, utils.GlobalConf.Mysql.DBPassword,
		utils.GlobalConf.Mysql.DBAddress, utils.GlobalConf.Mysql.DBName)
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
	sqlDb.SetMaxOpenConns(utils.GlobalConf.Mysql.DBMaxOpenConns)
	sqlDb.SetMaxIdleConns(utils.GlobalConf.Mysql.DBMaxIdleConns)
	sqlDb.SetConnMaxLifetime(time.Duration(utils.GlobalConf.Mysql.DBMaxLifeTime) * time.Second)

	return nil
}
