// @APIVersion 0.1.0
// @Title GoHumanLoop-Wework
// @Description 是针对GoHumanLoop在企业微信场景下进行审批、获取信息操作的示例服务。方便用户在使用`GohumanLoop`时，对接到自己的企业微信环境中。
// @Contact baird0917@163.com
// @License MIT
package routers

import (
	"github.com/beego/beego/logs"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context"
	"github.com/ptonlix/gohumanloop-wework/controllers"
	"github.com/ptonlix/gohumanloop-wework/models"
	"github.com/ptonlix/gohumanloop-wework/services"
)

func BearerTokenAuth(ctx *context.Context) {
	authHeader := ctx.Input.Header("Authorization")
	logs.Info("authHeader: %v", authHeader)
	if authHeader == "" || len(authHeader) < 7 || authHeader[:7] != "Bearer " {
		ctx.Output.SetStatus(401)
		ctx.Output.JSON(models.APIResponse{Success: false, Error: "Unauthorized"}, false, false)
		return
	}
	token := authHeader[7:]
	// 从数据库获取
	apiKey, err := services.APIKeyDataService.GetAPIKeyByKey(token)
	if err != nil {
		ctx.Output.SetStatus(401)
		ctx.Output.JSON(models.APIResponse{Success: false, Error: "Unauthorized"}, false, false)
		return
	}
	// 检查是否有效
	if apiKey.Status != true {
		ctx.Output.SetStatus(401)
		ctx.Output.JSON(models.APIResponse{Success: false, Error: "Unauthorized"}, false, false)
		return
	}
}

func init() {
	drc := &controllers.WeWorkController{}
	hl, _ := drc.NewHttpHandler()
	beego.Handler("/gohumanloop/callback", hl)

	nsGoHumanLoop := beego.NewNamespace("/api/v1",
		beego.NSBefore(func(ctx *context.Context) {
			if ctx.Input.URL() == "/api/v1/apikey/create" {
				return // 跳过认证
			}
			BearerTokenAuth(ctx) // 其他路由执行认证
		}),
		beego.NSNamespace("/humanloop",
			beego.NSCtrlPost("/request", (*controllers.GoHumanLoopController).Request),
			beego.NSCtrlGet("/status", (*controllers.GoHumanLoopController).Status),
			beego.NSCtrlPost("/continue", (*controllers.GoHumanLoopController).Continue),
			beego.NSCtrlPost("/cancel", (*controllers.GoHumanLoopController).Cancel),
			beego.NSCtrlPost("/cancel_converstation", (*controllers.GoHumanLoopController).CancelConversation),
		),
		beego.NSNamespace("/apikey",

			beego.NSCtrlPost("/create", (*controllers.APIKeyController).CreateKey),
			beego.NSCtrlGet("/get", (*controllers.APIKeyController).GetAPIKeyByKey),
			beego.NSCtrlPost("/update", (*controllers.APIKeyController).UpdateKey),
			beego.NSCtrlPost("/delete", (*controllers.APIKeyController).DeleteKey),
			beego.NSCtrlGet("/list", (*controllers.APIKeyController).ListKeys),
			beego.NSCtrlPut("/enable", (*controllers.APIKeyController).EnableKey),
			beego.NSCtrlPut("/disable", (*controllers.APIKeyController).DisableKey),
		),
	)

	beego.AddNamespace(nsGoHumanLoop)
}
