package util

import (
	"os"
)

const (
	B  = 1
	KB = 1024 * B
	MB = 1024 * KB
	GB = 1024 * MB
)

// IfExist 返回是否在硬盘中存在想要的文件
func IfExist(hash string) bool {
	_, err := os.Stat(ToFileRoute(hash))
	if err != nil {
		return false //文件不存在
	}
	return true
}

// ToFileRoute 返回hash值的文件路径
func ToFileRoute(hash string) string {
	return "./file/" + hash + "./"
}
