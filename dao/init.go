package dao

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type DB struct {
	*gorm.DB
}

func InitDB() *DB {
	var err error
	dsn := "bt_go:07ZjxQG87NTO9hUL@tcp(175.178.40.245:3309)/llfile?charset=utf8mb4&parseTime=True"
	db := new(DB)
	db.DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		//Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Println(err)
		return nil
	}
	fmt.Println(db.CreateBatchSize)
	return db
}
