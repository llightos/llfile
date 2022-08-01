package msg

//成功返回信息
const (
	UploadSuccess   = "文件上传成功"
	DownloadSuccess = "文件下载成功"
)

const (
	NoUploadEventId  = "没有找到字节断点信息"
	ErrUploadEventId = "错误的字节断点信息"
	ErrBreakpoint    = "错误的字节断点"
	ErrUploadFile    = "上传文件失败"
)
