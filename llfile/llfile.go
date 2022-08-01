package llfile

//llfile层对实际的file进行管理，维持在线用户的文件树（当然构建是model层来完成）
const (
	//FilePath 是实体file的存储路径，所有用户文件夹都是抽象出来的，实体文件夹以hash命名，
	//在一个统一放在一个文件夹里，
	//后期更改以时间为索引建文件夹
	FilePath = "./file"
)

func initFileRoute(path string) {

}
