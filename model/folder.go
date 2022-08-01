package model

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"log"
)

// 程序内部必须动态的维持一张Folder树，用于增删改查
// 用户文件夹
type Folder struct {
	gorm.Model

	UserID uint

	ManagerID *uint
	//是否是用户根目录
	RootNode bool
	ChildF   []*Folder `gorm:"foreignkey:ManagerID"` //子文件夹
	ParentF  *Folder   `gorm:"foreignkey:ManagerID"` //父文件夹

	Name string //文件夹名

	Files []File
}

// InitFolder 初始化用户文件网
func (d *DB) InitFolder(userid uint) {
	folder := new(Folder)
	folder.Name = "root"
	folder.RootNode = true
	folder.UserID = userid
	db.Create(folder)
}

// AddChild 给指定id的文件夹添加子文件夹
func (d *DB) AddChild(folderId, userId uint, newName string) (uint, error) {
	folder := new(Folder)

	fmt.Println(d.db)
	d.db.Table("folders").Where("id = ?", folderId).Preload("ChildF").Find(folder)
	if userId != folder.UserID || folder == nil {
		return 0, errors.New("error user or folders")
	}
	if folder.IfNameExist(newName) {
		return 0, errors.New("这个文件夹下文件夹名已存在")
	}
	fmt.Println("child", folder.ChildF)
	child := folder.addChild(newName)
	folder.ChildF = append(folder.ChildF, child)
	fmt.Println(*folder)
	//开启事务
	tx := d.db.Begin()
	//tx.Create(&child)
	err := db.Session(&gorm.Session{FullSaveAssociations: true}).Updates(&folder).Error
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	tx.Commit()
	return 0, nil
}

func (d *DB) AddFolderFile(folderID uint, file File) error {
	folder := new(Folder)
	db.Model(&Folder{}).Where("id = ?", folderID).First(folder)

	folder.Files = append(folder.Files, file)
	err := db.Session(&gorm.Session{FullSaveAssociations: true}).Updates(folder).Error
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// FolderPath 返回目标文件夹id的路径,与GetUserTree形成链式,勉强能用，太寄吧慢了，后面重构
func FolderPath(targetID uint) string {
	f := new(Folder)
	f.ID = targetID
	path := ""
	//一级拿到相关targetID的对象，当然，其ParentF只有一级
	err := db.Preload("ParentF").Preload(clause.Associations).First(f).Error
	if err != nil {
		log.Println(err)
		return ""
	}
	fmt.Println(f)
	//	fmt.Println(f.ParentF.ID)
	head := f
	for {
		fmt.Println("path", path)
		//遍历每一级拿到相关对象，当然，其在一次循环中ParentF只有一级
		err = db.Preload("ParentF").Preload(clause.Associations).First(head).Error
		if err != nil {
			log.Println(err)
			return ""
		}
		if head.ParentF == nil {
			path = "/root/" + path
			return path
		}
		path = head.Name + "/" + path
		head = head.ParentF
	}
	return path
}

func FolderTree(userid int) *Folder {
	folders := GetUserAllFolder(userid)
	for i, folder := range folders {
		//找到root
		if folder.RootNode == true {
			folders = append(folders[0:i], folders[i+1:]...)
			fmt.Println(len(folders))
			fmt.Println(folders)
			folder.build(folders)
			return folder
		}
	}
	return nil
}

func GetUserAllFolder(userid int) []*Folder {
	folders := make([]*Folder, 0)
	err := db.Model(&Folder{}).Where("user_id = ?", userid).Preload("Files").Find(&folders).Error
	if err != nil {
		err.Error()
	}
	fmt.Println("ss", folders)
	return folders
}

// build 对外调用时f是user的根节点，
// 再来理解下递归，递归分为被影响量和影响量，
// 这里来说，f就是被影响量， 而folders就是影响量，f和folders是链式的也是线性的
// 最开始的f（这里是外调用放root），f是跟着folders的改变而变，也就是两个都得变，不然就没有影响因子
func (f *Folder) build(folders []*Folder) {
	//if len(folders) == 0 {
	//	return
	//}
	fmt.Println("f.", f.ID)
	for i, folder := range folders {
		if f.ID == *folder.ManagerID {
			//fmt.Println("how much")
			folder.ParentF = f
			f.ChildF = append(f.ChildF, folder)
			folders = append(folders[:i], folders[i+1:]...)
			folder.build(folders)
		}
	}
}

// IfNameExist 在f中是否name存在（f是预加载ChildF之后的）, true是存在
func (f *Folder) IfNameExist(name string) bool {
	for _, v := range f.ChildF {
		if v.Name == name {
			return true
		}
	}
	return false
}
