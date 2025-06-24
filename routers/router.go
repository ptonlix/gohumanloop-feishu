// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
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
		ctx.Output.JSON(map[string]string{"error": "Unauthorized"}, false, false)
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
		beego.NSNamespace("/humanloop",
			beego.NSBefore(BearerTokenAuth),
			beego.NSCtrlPost("/request", (*controllers.GoHumanLoopController).Request),
			beego.NSCtrlGet("/status", (*controllers.GoHumanLoopController).Status),
			beego.NSCtrlPost("/continue", (*controllers.GoHumanLoopController).Continue),
			beego.NSCtrlPost("/cancel", (*controllers.GoHumanLoopController).Cancel),
			beego.NSCtrlPost("/cancel_converstation", (*controllers.GoHumanLoopController).CancelConversation),
		),
		beego.NSNamespace("/apikey",
			beego.NSBefore(BearerTokenAuth),
			beego.NSCtrlGet("/get", (*controllers.APIKeyController).GetAPIKeyByKey),
			beego.NSCtrlPost("/update", (*controllers.APIKeyController).UpdateKey),
			beego.NSCtrlPost("/delete", (*controllers.APIKeyController).DeleteKey),
			beego.NSCtrlGet("/list", (*controllers.APIKeyController).ListKeys),
		),
	)

	nsNoAuth := beego.NewNamespace("/api/v1",
		beego.NSCtrlPost("/create_apikey", (*controllers.APIKeyController).CreateKey),
	)

	beego.AddNamespace(nsGoHumanLoop, nsNoAuth)
}
