// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"github.com/ptonlix/gohumanloop-wework/controllers"

	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	drc := &controllers.WeWorkController{}
	hl, _ := drc.NewHttpHandler()
	beego.Handler("/gohumanloop", hl)

	nsGoHumanLoop := beego.NewNamespace("/api/v1",
		beego.NSNamespace("/humanloop",
			beego.NSCtrlPost("/request", (*controllers.GoHumanLoopController).Request),
			beego.NSCtrlGet("/status", (*controllers.GoHumanLoopController).Status),
			beego.NSCtrlPost("/continue", (*controllers.GoHumanLoopController).Continue),
			beego.NSCtrlPost("/cancel", (*controllers.GoHumanLoopController).Cancel),
			beego.NSCtrlPost("/cancel_converstation", (*controllers.GoHumanLoopController).CancelConversation),
		),
	)

	beego.AddNamespace(nsGoHumanLoop)
}
