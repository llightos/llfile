package model

// addChild 返回一个子文件夹，不更新数据库
func (f *Folder) addChild(name string) *Folder {
	newFolder := new(Folder)
	newFolder.Name = name
	newFolder.UserID = f.UserID
	newFolder.ParentF = f
	return newFolder
}
