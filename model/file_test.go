package model

import (
	"fmt"
	"github.com/gogf/gf/v2/crypto/gmd5"
	"llfile/dao"
	"testing"
)

func Test_init(t *testing.T) {
	db = dao.InitDB()
	err := db.AutoMigrate(&UserInfo{}, &File{}, &Folder{})
	if err != nil {
		panic(err)
	}
}

func TestFolder_AddUserFile(t *testing.T) {
	NewModelDB().AddFolderFile(1, File{
		Hash:         gmd5.MustEncrypt("这是一个测试文件"),
		Name:         "test",
		ExpandedName: "txt",
		Size:         100,
		Path:         "/root",
	})
}

func TestDB_AddFile(t *testing.T) {
	NewModelDB().AddFile(File{
		FolderID:     2,
		Hash:         "diuadnajcnfanfwa",
		Name:         "111",
		ExpandedName: "txt",
		Size:         771,
		Path:         FolderPath(2),
	})
}

func TestDB_FindFile(t *testing.T) {
	file := NewModelDB().FindFile("8", "test", "zip")
	fmt.Println(file)
}

func TestQueryUserFiles(t *testing.T) {
	files := NewModelDB().QueryUserFiles("8")
	fmt.Println(files)
}

func TestQueryUserFolders(t *testing.T) {
	folders := NewModelDB().QueryUserFolders("7")
	fmt.Println(folders)
}

func TestGorm(t *testing.T) {
	NewModelDB().db.Table("")
}
