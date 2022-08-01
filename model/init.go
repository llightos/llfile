package model

import "llfile/dao"

// 对临时内调用， 主要用于表的结构AutoMigrate
var db *dao.DB

// DB 对外调用的对象，拥有对外的所有方法
type DB struct {
	db *dao.DB
}

func NewModelDB() *DB {
	localDB := new(DB)
	localDB.db = db
	return localDB
}

func init() {
	db = dao.InitDB()
	//err := db.AutoMigrate(&UserInfo{}, &File{}, &UserFile{}, &Folder{})
	//if err != nil {
	//	panic(err)
	//}
}
