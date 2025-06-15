package controllers

import (
	"net/http"

	"github.com/beego/beego/logs"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/ptonlix/gohumanloop-wework/init/wework"
	"github.com/xen0n/go-workwx/v2"
)

type WeWorkController struct {
	beego.Controller
}

// OnIncomingMessage 一条消息到来时的回调。
func (dr *WeWorkController) OnIncomingMessage(msg *workwx.RxMessage) error {
	// You can do much more!
	logs.Debug("incoming message: %s\n", msg)
	switch msg.Event {
	case workwx.EventTypeSysApprovalChange:
		// 更新审批状态
		if data, flag := msg.EventSysApprovalChange(); flag {
			logs.Debug("审批状态更新: %s\n", data.GetApprovalInfo())
			// TODO: 审批状态更新
		}
	}

	return nil
}

func (dr *WeWorkController) NewHttpHandler() (http.Handler, error) {
	hh, err := workwx.NewHTTPHandler(wework.WeworkConf.PToken, wework.WeworkConf.PKey, dr)
	if err != nil {
		panic(err)
	}
	return hh, err
}
