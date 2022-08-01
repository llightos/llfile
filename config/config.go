package config

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/glog"
)

var Logger *glog.Logger

func init() {
	Logger = g.Log()
	g.Log().Print(context.TODO(), "log start")
}
