package model

// VerifyFolderAndUser 验证文件夹是不是用户所属
func VerifyFolderAndUser(folderID string, UserID uint) (f *Folder, ok bool) {
	f = new(Folder)
	err := db.Table("folders").Where("id = ? AND user_id = ?", folderID, UserID).First(&f).Error
	if err != nil {
		return nil, false
	}
	if f == nil {
		return nil, false
	}
	return f, true
}
