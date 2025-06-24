package services

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/ptonlix/gohumanloop-wework/init/wework"
	"github.com/ptonlix/gohumanloop-wework/models"

	"github.com/beego/beego/v2/core/logs"
	"github.com/xen0n/go-workwx/v2"
)

type WeWorkService struct {
	app *workwx.WorkwxApp
}

func NewWeWorkService() *WeWorkService {

	client := workwx.New(wework.WeworkConf.CorpId)

	// work with individual apps
	app := client.WithApp(wework.WeworkConf.CorpSecret, int64(wework.WeworkConf.AgentId))
	app.SpawnAccessTokenRefresher()
	return &WeWorkService{app: app}
}

var GHLWeWorkService *WeWorkService = NewWeWorkService()

func NewwWeWorkServiceWithCtx(ctx context.Context) *WeWorkService {

	client := workwx.New(wework.WeworkConf.CorpId)

	// work with individual apps
	app := client.WithApp(wework.WeworkConf.CorpSecret, int64(wework.WeworkConf.AgentId))
	app.SpawnAccessTokenRefresherWithContext(ctx)
	return &WeWorkService{app: app}
}

func (w *WeWorkService) UpdateUser(userId, externalUserId, unionId string) error {
	if err := w.app.RemarkExternalContact(&workwx.ExternalContactRemark{Userid: userId, ExternalUserid: externalUserId, Description: unionId}); err != nil {
		logs.Error("UpdateUser is failed!:", err)
		return err
	}

	return nil
}

func (w *WeWorkService) SendWelcome(welcomeCode, welcomeText, title, picurl, desc, url string) error {

	att := workwx.Attachments{MsgType: workwx.AttachmentMsgTypeLink, Link: workwx.Link{Title: title, PicURL: picurl, Desc: desc, URL: url}}

	if err := w.app.SendWelcomeMsg(welcomeCode, workwx.Text{Content: welcomeText}, []workwx.Attachments{att}); err != nil {
		logs.Error("SendWelcomeMsg is failed!:", err)
		return err
	}

	return nil
}

func (w *WeWorkService) GetUserInfo(externalUserId string) (*workwx.ExternalContactInfo, error) {
	resp, err := w.app.GetExternalContact(externalUserId)
	if err != nil {
		logs.Error("GetUser is failed!:", err)
		return nil, err
	}

	return resp, nil
}

func (w *WeWorkService) HumanLoopRequestOA(humanLoopRequest models.HumanLoopRequestData) (string, error) {
	// 获取审批人ID
	approverUserId := ""
	if approver, ok := humanLoopRequest.Metadata["approverid"]; ok {
		if approverStr, ok := approver.(string); ok {
			approverUserId = approverStr
		} else {
			approverUserId = ""
		}
	} else {
		approverUserId = wework.WeworkConf.ApproverUserId
	}

	var templateDetail *workwx.OATemplateDetail
	templateId := ""
	if humanLoopRequest.LoopType == models.HumanLoopTypeApprove {
		resp, err := w.app.GetOATemplateDetail(wework.WeworkConf.ApproveTemplateId)
		if err != nil {
			logs.Error("GetOATemplateDetail is failed!:", err)
			return "", err
		}
		logs.Info("GetOATemplateDetail: %v", resp)
		templateDetail = resp
		templateId = wework.WeworkConf.ApproveTemplateId

	} else if humanLoopRequest.LoopType == models.HumanLoopTypeInformation {
		resp, err := w.app.GetOATemplateDetail(wework.WeworkConf.InfoTemplateId)
		if err != nil {
			logs.Error("GetOATemplateDetail is failed!:", err)
			return "", err
		}
		logs.Info("GetOATemplateDetail: %v", resp)
		templateDetail = resp
		templateId = wework.WeworkConf.InfoTemplateId
	} else {
		logs.Error("HumanLoopRequestOA Type is not approve or information: %v", humanLoopRequest.LoopType)
		return "", errors.New("HumanLoopRequestOA Type is not approve or information")
	}

	applyContents := workwx.OAContents{}
	for _, controls := range templateDetail.TemplateContent.Controls {
		value := ""
		for _, title := range controls.Property.Title {
			if title.Text == "任务ID" {
				value = humanLoopRequest.TaskId
			} else if title.Text == "对话ID" {
				value = humanLoopRequest.ConversationId
			} else if title.Text == "请求ID" {
				value = humanLoopRequest.RequestId
			} else if title.Text == "HumanLoop类型" {
				value = humanLoopRequest.LoopType
			} else if title.Text == "申请内容" {
				// 安全地获取并转换Context中的值
				if msg, ok := humanLoopRequest.Context["message"]; ok {
					if msgStr, ok := msg.(string); ok {
						value = msgStr
					} else {
						value = ""
					}
				}
			} else if title.Text == "申请问题" {
				if question, ok := humanLoopRequest.Context["question"]; ok {
					if questionStr, ok := question.(string); ok {
						value = questionStr
					} else {
						value = ""
					}
				}
			} else if title.Text == "申请说明" {
				if additional, ok := humanLoopRequest.Context["additional"]; ok {
					if additionalStr, ok := additional.(string); ok {
						value = additionalStr
					} else {
						value = ""
					}
				}
			}
		}
		applyContents.Contents = append(applyContents.Contents, workwx.OAContent{
			ID:      controls.Property.ID,
			Control: controls.Property.Control,
			Value:   workwx.OAContentValue{Text: value, BankAccount: workwx.OAContentBankAccount{AccountType: 2}},
		})
	}

	humanLoopQuestion := ""
	if question, ok := humanLoopRequest.Context["question"]; ok {
		if questionStr, ok := question.(string); ok {
			humanLoopQuestion = questionStr
		} else {
			humanLoopQuestion = ""
		}
	}

	applyevent := workwx.OAApplyEvent{
		CreatorUserID: wework.WeworkConf.CreatorUserId,
		TemplateID:    templateId,
		Approver: []workwx.OAApprover{
			{Attr: 2, UserID: []string{approverUserId}},
		},
		ApplyData: applyContents,
		SummaryList: []workwx.OASummaryList{
			{
				SummaryInfo: []workwx.OAText{{Text: humanLoopQuestion, Lang: "zh_CN"}},
			},
		},
	}
	json_applyevent, _ := json.Marshal(applyevent)
	logs.Debug("ApplyOAEvent: %v", string(json_applyevent))

	spNo, err := w.app.ApplyOAEvent(applyevent)
	if err != nil {
		logs.Error("ApplyOAEvent is failed!:", err)
		return "", err
	}

	return spNo, nil
}
