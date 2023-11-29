package db

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

// DB 数据库链接单例
var DB *gorm.DB

func Init(connString string) {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level(这里记得根据需求改一下)
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,        // Disable color
		},
	)
	log.Println("db dsn " + connString)
	db, err := gorm.Open(mysql.Open(connString), &gorm.Config{
		Logger:                                   newLogger,
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		fmt.Println("storage err: ", err)
	} else {
		fmt.Println("数据库连接成功！！")
	}
	sqlDB, err := db.DB()
	if err != nil {
		fmt.Println("storage err: ", err)
	}
	sqlDB.SetMaxIdleConns(3)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)
	DB = db
}
