package config

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/glog"
)

var Logger *glog.Logger
var ServerIPAddress = "127.0.0.1:8081"
var LimitSpeed = 200.0 //限速为32kb的倍数 这里就是3.2mb
var LimitSpeedInt = 200

func init() {
	Logger = g.Log()
	g.Log().Print(context.TODO(), "log start")
}
