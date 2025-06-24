package controllers

import (
	"net/http"

	"github.com/beego/beego/logs"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/ptonlix/gohumanloop-wework/init/wework"
	"github.com/ptonlix/gohumanloop-wework/models"
	"github.com/ptonlix/gohumanloop-wework/services"
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
			GHLStatus := ""
			responseBy := ""
			response := ""
			if data.GetApprovalInfo().SpStatus == "1" { // 审批中
				GHLStatus = models.HumanLoopStatusPending
			} else if data.GetApprovalInfo().SpStatus == "2" { //审批完成
				GHLData, err := services.GHLDataService.GetHumanLoopByBySpNo(data.GetApprovalInfo().SpNo)
				if err != nil {
					logs.Error("GetHumanLoopByBySpNo failed: %v", err)
				}
				if GHLData.LoopType == models.HumanLoopTypeApprove {
					GHLStatus = models.HumanLoopStatusApproved
				} else if GHLData.LoopType == models.HumanLoopTypeInformation {

					GHLStatus = models.HumanLoopStatusCompleted
				}
				responseBy = data.GetApprovalInfo().SpRecord[len(data.GetApprovalInfo().SpRecord)-1].Details[0].Approver.UserID
				response = data.GetApprovalInfo().Comments[len(data.GetApprovalInfo().Comments)-1].CommentContent
			} else if data.GetApprovalInfo().SpStatus == "3" { //拒绝审批
				GHLStatus = models.HumanLoopStatusRejected
				responseBy = data.GetApprovalInfo().SpRecord[len(data.GetApprovalInfo().SpRecord)-1].Details[0].Approver.UserID
				response = data.GetApprovalInfo().Comments[len(data.GetApprovalInfo().Comments)-1].CommentContent
			} else if data.GetApprovalInfo().SpStatus == "4" { //取消申请
				GHLStatus = models.HumanLoopStatusCancelled
			} else {
				GHLStatus = models.HumanLoopStatusError //审批出错
			}
			// 1. 更新审批状态
			id, err := services.GHLDataService.UpdateHumanLoopStatusBySpNo(data.GetApprovalInfo().SpNo, GHLStatus, response, responseBy)
			if err != nil {
				logs.Error("UpdateHumanLoopStatus failed: %v", err)
			}
			logs.Info("UpdateHumanLoopStatus id: %v successfully ", id)
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
