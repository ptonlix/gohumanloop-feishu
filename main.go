package main

import (
	_ "github.com/ptonlix/gohumanloop-feishu/init/feishu"

	_ "github.com/ptonlix/gohumanloop-feishu/routers"

	"github.com/beego/beego/v2/core/logs"

	beego "github.com/beego/beego/v2/server/web"

	_ "github.com/ptonlix/gohumanloop-feishu/docs"
)

var (
	AppName      string // 应用名称
	AppVersion   string // 应用版本
	BuildVersion string // 编译版本
	BuildTime    string // 编译时间
	GitRevision  string // Git版本
	GitBranch    string // Git分支
	GoVersion    string // Golang信息
)

func init() {
	logs.SetLogger(logs.AdapterFile, `{"filename":"./log/project.log","level":7,"maxlines":0,"maxsize":0,"daily":true,"maxdays":10,"color":true}`)
}

// @title           gohumanloop-feishu
// @version         v0.1.0
// @description     针对GoHumanLoop在企业微信场景下进行审批、获取信息操作的示例服务。方便用户在使用`GohumanLoop`时，对接到自己的企业微信环境中。
// @termsOfService  http://swagger.io/terms/

// @contact.name   Baird
// @contact.url    https://github.com/ptonlix/gohumanloop-feishu
// @contact.email  baird0917@163.com

// @license.name  MIT
// @license.url   https://github.com/ptonlix/gohumanloop-feishu/blob/main/LICENSE

// @host      127.0.0.1:9800
// @BasePath  /api/v1

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/docs"] = "docs"
	}
	Version()
	beego.Run()
}

// Version 版本信息
func Version() {
	logs.Info("App Name:\t\t%s", AppName)
	logs.Info("App Version:\t%s", AppVersion)
	logs.Info("Build version:\t%s", BuildVersion)
	logs.Info("Build time:\t%s", BuildTime)
	logs.Info("Git revision:\t%s", GitRevision)
	logs.Info("Git branch:\t%s", GitBranch)
	logs.Info("Golang Version:\t%s", GoVersion)
}
