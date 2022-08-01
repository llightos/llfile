package model

import (
	"fmt"
	"testing"
)

//func TestFolder(t *testing.T) {
//	NewModelDB().AddChild(1, 2, "nihao")
//}

func TestInit(t *testing.T) {
	modelDB := NewModelDB()
	fmt.Println(modelDB)
	//NewModelDB().db.AutoMigrate(&UserInfo{}, &File{}, &UserFile{}, &Folder{})
}

func TestAddUser(t *testing.T) {
	NewModelDB().AddUser("username1", "111111")
}

func TestDB_AddFolder(t *testing.T) {
	//NewModelDB().AddChild(1, 1, "222")
	_, err := NewModelDB().AddChild(1, 1, "new")
	if err != nil {
		fmt.Println(err)
	}
}

//
func TestFolder_FolderPath(t *testing.T) {
	tempFolder := new(Folder)
	NewModelDB().db.Model(&Folder{}).Where("id = ?", 4).First(tempFolder)
	fmt.Println(FolderPath(4))
}

func TestGetUserTree(t *testing.T) {

}

func TestGetAllFolder(t *testing.T) {
	fmt.Println(GetUserAllFolder(1))
}

func TestFolderTree(t *testing.T) {
	tree := FolderTree(1)
	fmt.Println(tree)
}
