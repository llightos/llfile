package model

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"strconv"
)

// File 公共文件存储池
type File struct {
	gorm.Model
	FolderID     uint
	Hash         string
	Name         string `gorm:"type:varchar(255)"`
	ExpandedName string `gorm:"type:varchar(255)"`
	Size         uint
	Path         string `gorm:"type:varchar(255)"`
}

// SharedFile 用户拥有的他人分享的文件
type SharedFile struct {
	gorm.Model
	Owner uint //拥有者的id
	// 以下所有的字段都是独立的
	FolderID     uint
	Hash         string
	Name         string
	ExpandedName string
	Size         uint
	Path         string
}

type QueryUserFiles struct {
	UserID     uint
	FileID     uint `gorm:"column:id"`
	FolderID   uint
	Name       string
	ExpandName string `gorm:"column:expanded_name"`
	Size       uint
	Path       string
}

// QueryUserFiles 获取用户的所有文件和路劲
func (d *DB) QueryUserFiles(userId string) []QueryUserFiles {
	res := make([]QueryUserFiles, 0)
	db.Raw("select user_id, files.id, folder_id, files.name, files.expanded_name, files.size, files.path "+
		"from files, folders "+
		"where folder_id = folders.id and user_id = ?", userId).Find(&res)
	return res
}

type QueryUserFolders struct {
	FolderID uint `gorm:"column:id"`
	UserID   uint
	FatherID uint `gorm:"column:manager_id"`
	RootNode bool
	Name     string
}

func (d *DB) QueryUserFolders(userId string) []QueryUserFolders {
	res := make([]QueryUserFolders, 0)
	db.Raw("select id, user_id, manager_id, root_node, name "+
		"from folders "+
		"where user_id = ?", userId).Find(&res)

	return res
}

// 用户个人存储文件

func (d *DB) AddFile(file File) (uint, error) {
	var id uint
	err := db.Transaction(func(tx *gorm.DB) error {
		err := tx.Create(&file).Error
		if err != nil {
			return err
		}
		id = file.ID
		return nil
	})
	return id, err
}

// ChangeFileRoute 改变文件Path
func (d *DB) ChangeFileRoute(fileId, newFolderIdm, userId int, newHeadName, newExpandName string) error {
	var outErr error
	db.Transaction(func(tx *gorm.DB) error {
		fmt.Println("IN????????")
		var id int
		//判断目标文件夹是不是用户的
		tx.Raw("select user_id "+
			"from folders, files "+
			"where files.folder_id = folders.id and files.id = ? and folders.user_id = ?;", fileId, userId).First(&id)
		fmt.Println(id)
		if id == 0 {
			return errors.New("无权限")
		}
		// 改名字
		fmt.Println("fileID", fileId)
		if newHeadName != "" && newExpandName != "" {
			err := tx.Table("files").Where("id = ?", fileId).Updates(&File{
				Path:         FolderPath(uint(newFolderIdm)),
				Name:         newHeadName,
				ExpandedName: newExpandName,
			}).Error
			if err != nil {
				outErr = err
				return err
			}
		}
		if newFolderIdm != 0 {
			err := tx.Table("files").Where("id = ?", fileId).Updates(&File{
				FolderID: uint(newFolderIdm),
				Path:     FolderPath(uint(newFolderIdm)),
			}).Error
			if err != nil {
				outErr = err
			}
			return err
		}
		return nil
	})
	fmt.Println("cw", outErr)
	return outErr
}

// IfHashExist userId是否有目标hash文件, bool1，返回所有位置是否有对应hash文件，
// bool2返回用户是否拥有对应文件
func (d *DB) IfHashExist(userId, hash string) (bool1 bool, bool2 bool) {
	files := make([]File, 0)
	fmt.Println("hash = ", hash)
	d.db.DB.Table("files").Where("hash = ?", hash).Find(&files)
	if len(files) == 0 {
		return false, false
	}
	for _, v := range files {
		fmt.Println("user_id = ", strconv.Itoa(int(v.ID)), userId)
		if strconv.Itoa(int(v.ID)) == userId {
			return true, true
		}
	}
	return true, false
}

func (d *DB) FindFile(userId, headName, expandName string) (file File) {
	file = *new(File)
	d.db.Raw("select files.hash , size from folders, files "+
		"where folders.id = files.folder_id and folders.user_id = ? and files.name = ? and expanded_name = ?", userId, headName, expandName).
		Scan(&file)
	return
}

func (d *DB) FindFileByFileId(userId, fileId string) (file File) {
	file = *new(File)
	fmt.Println("shifou????")
	d.db.Raw("select files.name, files.expanded_name,files.hash, files.size, files.folder_id "+
		"from folders, files "+
		"where folders.id = files.folder_id and folders.user_id = ? and files.id = ?", userId, fileId).
		Scan(&file)
	fmt.Println(file.Hash)
	return
}
