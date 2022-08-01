package handler

import "llfile/model"

var DB *model.DB

func init() {
	DB = model.NewModelDB()
}
