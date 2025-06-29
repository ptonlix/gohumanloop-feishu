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
	if msg == nil {
		logs.Error("received nil message")
		return nil
	}

	logs.Debug("incoming message: %s\n", msg)

	switch msg.Event {
	case workwx.EventTypeSysApprovalChange:
		data, flag := msg.EventSysApprovalChange()
		if !flag {
			logs.Error("failed to parse approval change event")
			return nil
		}

		approvalInfo := data.GetApprovalInfo()
		logs.Debug("审批状态更新: %s\n", approvalInfo)

		var GHLStatus, responseBy, response, feedback string

		switch approvalInfo.SpStatus {
		case "1": // 审批中
			GHLStatus = models.HumanLoopStatusPending
		case "2": // 审批完成
			GHLData, err := services.GHLDataService.GetHumanLoopByBySpNo(approvalInfo.SpNo)
			if err != nil {
				logs.Error("GetHumanLoopByBySpNo failed: %v", err)
				return nil
			}

			if GHLData.LoopType == models.HumanLoopTypeApprove {
				GHLStatus = models.HumanLoopStatusApproved
			} else if GHLData.LoopType == models.HumanLoopTypeInformation {
				GHLStatus = models.HumanLoopStatusCompleted
			}

			// 安全获取审批人和评论
			if len(approvalInfo.SpRecord) > 0 && len(approvalInfo.SpRecord[len(approvalInfo.SpRecord)-1].Details) > 0 {
				responseBy = approvalInfo.SpRecord[len(approvalInfo.SpRecord)-1].Details[0].Approver.UserID
				response = approvalInfo.SpRecord[len(approvalInfo.SpRecord)-1].Details[0].Speech
			}

			if len(approvalInfo.Comments) > 0 {
				response = approvalInfo.Comments[len(approvalInfo.Comments)-1].CommentContent
			}

		case "3": // 拒绝审批
			GHLStatus = models.HumanLoopStatusRejected
			// 安全获取审批人和评论
			if len(approvalInfo.SpRecord) > 0 && len(approvalInfo.SpRecord[len(approvalInfo.SpRecord)-1].Details) > 0 {
				responseBy = approvalInfo.SpRecord[len(approvalInfo.SpRecord)-1].Details[0].Approver.UserID
				response = approvalInfo.SpRecord[len(approvalInfo.SpRecord)-1].Details[0].Speech
			}

			if len(approvalInfo.Comments) > 0 {
				feedback = approvalInfo.Comments[len(approvalInfo.Comments)-1].CommentContent
			}

		case "4": // 取消申请
			GHLStatus = models.HumanLoopStatusCancelled
		default:
			GHLStatus = models.HumanLoopStatusError // 审批出错
			logs.Warn("unknown approval status: %s", approvalInfo.SpStatus)
		}

		// 更新审批状态
		logs.Debug("更新审批状态: %s, %s, %s, %s, %s", approvalInfo.SpNo, GHLStatus, response, responseBy, feedback)
		id, err := services.GHLDataService.UpdateHumanLoopStatusBySpNo(
			approvalInfo.SpNo,
			GHLStatus,
			response,
			responseBy,
			feedback,
		)
		if err != nil {
			logs.Error("UpdateHumanLoopStatus failed: %v", err)
			return nil
		}
		logs.Info("UpdateHumanLoopStatus id: %v successfully", id)
	}

	return nil
}

func (dr *WeWorkController) NewHttpHandler() (http.Handler, error) {
	hh, err := workwx.NewHTTPHandler(wework.WeworkConf.PToken, wework.WeworkConf.PKey, dr)
	if err != nil {
		logs.Error("Failed to create HTTP handler: %v", err)
		return nil, err
	}
	return hh, nil
}
