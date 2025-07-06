package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/beego/beego/v2/core/logs"
	larkapproval "github.com/larksuite/oapi-sdk-go/v3/service/approval/v4"
	"github.com/ptonlix/gohumanloop-feishu/init/feishu"
	"github.com/ptonlix/gohumanloop-feishu/models"

	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
)

type FeishuService struct {
	app *lark.Client
}

func NewFeishuService() *FeishuService {
	client := lark.NewClient(feishu.FeishuConf.AppId, feishu.FeishuConf.AppSecret, // 默认配置为自建应用
		// lark.WithMarketplaceApp(), // 可设置为商店应用
		lark.WithLogLevel(larkcore.LogLevelDebug),
		lark.WithReqTimeout(60*time.Second),
		lark.WithEnableTokenCache(true),
		lark.WithHttpClient(http.DefaultClient))

	return &FeishuService{
		app: client,
	}
}

var GHLFeishuService *FeishuService = NewFeishuService()

func (f *FeishuService) GetApprovalDetail(approvalCode string) (string, []map[string]interface{}, error) {

	// 创建请求对象
	req := larkapproval.NewGetApprovalReqBuilder().
		ApprovalCode(feishu.FeishuConf.ApproveTemplateId).
		Locale(`zh-CN`).
		Build()
	// 发起请求
	resp, err := f.app.Approval.V4.Approval.Get(context.Background(), req)

	// 处理错误
	if err != nil {
		logs.Error("GetApproval is failed!:", err)
		return "", nil, err
	}

	// 服务端错误处理
	if !resp.Success() {
		logs.Error("logId: %s, error response: \n%s", resp.RequestId(), larkcore.Prettify(resp.CodeError))
		return "", nil, err
	}

	// 检查状态是否可用
	if *resp.Data.Status != "ACTIVE" {
		return "", nil, errors.New("approval is not available")
	}
	// 将返回的FORM 解析成字典数组
	var approvalForms []map[string]interface{}
	err = json.Unmarshal([]byte(*resp.Data.Form), &approvalForms)
	if err != nil {
		return "", nil, fmt.Errorf("解析表单数据失败: %v", err)
	}

	approvalNodeId := ""
	for _, node := range resp.Data.NodeList {
		if *node.Name == "审批" {
			approvalNodeId = *node.NodeId
		}
	}

	return approvalNodeId, approvalForms, nil
}

// 创建审批申请
func (f *FeishuService) CreateApproval(approvalCode, createUser, approverUserId string, humanLoopRequest *models.HumanLoopRequestData) (string, error) {
	approvalNodeId, approvalForms, err := f.GetApprovalDetail(approvalCode)
	if err != nil {
		return "", err
	}
	// 创建表单字段数组
	var formFields []map[string]interface{}
	for _, form := range approvalForms {
		field := make(map[string]interface{})
		field["id"] = form["id"].(string)
		field["type"] = form["type"].(string)

		switch form["name"] {
		case "任务ID":
			field["value"] = humanLoopRequest.TaskId
		case "对话ID":
			field["value"] = humanLoopRequest.ConversationId
		case "请求ID":
			field["value"] = humanLoopRequest.RequestId
		case "HumanLoop类型":
			field["value"] = humanLoopRequest.LoopType
		case "申请内容":
			if msg, ok := humanLoopRequest.Context["message"]; ok {
				if msgStr, ok := msg.(string); ok {
					field["value"] = msgStr
				} else {
					field["value"] = ""
				}
			}
		case "申请问题":
			if question, ok := humanLoopRequest.Context["question"]; ok {
				if questionStr, ok := question.(string); ok {
					field["value"] = questionStr
				} else {
					field["value"] = ""
				}
			}
		case "申请说明":
			if additional, ok := humanLoopRequest.Context["additional"]; ok {
				if additionalStr, ok := additional.(string); ok {
					field["value"] = additionalStr
				} else {
					field["value"] = ""
				}
			}
		}
		formFields = append(formFields, field)
	}

	formStr, err := json.Marshal(formFields)
	if err != nil {
		logs.Error("Marshal formData is failed!:", err)
		return "", err
	}

	// 创建请求对象
	logs.Info("formStr:", string(formStr))
	req := larkapproval.NewCreateInstanceReqBuilder().
		InstanceCreate(larkapproval.NewInstanceCreateBuilder().
			ApprovalCode(approvalCode).
			OpenId(createUser).
			Form(string(formStr)).
			NodeApproverOpenIdList([]*larkapproval.NodeApprover{
				larkapproval.NewNodeApproverBuilder().
					Key(approvalNodeId).
					Value([]string{approverUserId}).
					Build(),
			}).Build()).
		Build()

	// 发起请求
	resp, err := f.app.Approval.V4.Instance.Create(context.Background(), req)

	// 处理错误
	if err != nil {
		logs.Error("CreateApproval is failed!:", err)
		return "", err
	}

	// 服务端错误处理
	if !resp.Success() {
		logs.Error("logId: %s, error response: \n%s", resp.RequestId(), larkcore.Prettify(resp.CodeError))
		return "", errors.New(larkcore.Prettify(resp.CodeError))
	}

	approvalInstanceCode := resp.Data.InstanceCode

	return *approvalInstanceCode, nil
}

// 发起审批流程
func (f *FeishuService) HumanLoopRequestOA(humanLoopRequest models.HumanLoopRequestData) (string, error) {
	// 获取审批人ID
	approverUserId := ""
	if approver, ok := humanLoopRequest.Metadata["approverid"]; ok {
		if approverStr, ok := approver.(string); ok {
			approverUserId = approverStr
		} else {
			approverUserId = ""
		}
	} else {
		approverUserId = feishu.FeishuConf.ApproverUserId
	}

	var approvalInstanceCode string
	var err error
	if humanLoopRequest.LoopType == models.HumanLoopTypeApprove {
		approvalInstanceCode, err = f.CreateApproval(feishu.FeishuConf.ApproveTemplateId, feishu.FeishuConf.CreatorUserId, approverUserId, &humanLoopRequest)
		if err != nil {
			return "", err
		}
	} else if humanLoopRequest.LoopType == models.HumanLoopTypeInformation {
		approvalInstanceCode, err = f.CreateApproval(feishu.FeishuConf.InfoTemplateId, feishu.FeishuConf.CreatorUserId, approverUserId, &humanLoopRequest)
		if err != nil {
			return "", err
		}
	} else {
		logs.Error("HumanLoopRequestOA Type is not approve or information: %v", humanLoopRequest.LoopType)
		return "", errors.New("HumanLoopRequestOA Type is not approve or information")
	}

	return approvalInstanceCode, nil
}

// 开启事件订阅
func (f *FeishuService) SubscribeApprovalEvent(approvalInstanceCode string) error {
	// 创建请求对象
	req := larkapproval.NewSubscribeApprovalReqBuilder().
		ApprovalCode(approvalInstanceCode).
		Build()

	// 发起请求
	resp, err := f.app.Approval.V4.Approval.Subscribe(context.Background(), req)

	// 处理错误
	if err != nil {
		logs.Error("SubscribeApprovalEvent is failed!:", err)
		return err
	}

	// 服务端错误处理
	if !resp.Success() {
		logs.Error("logId: %s, error response: \n%s", resp.RequestId(), larkcore.Prettify(resp.CodeError))
		return err
	}

	// 订阅成功
	logs.Info("SubscribeApprovalEvent is success!")
	return nil
}

// 获取审批实例详情
func (f *FeishuService) GetApprovalInstanceDetail(approvalInstanceCode string) (*larkapproval.GetInstanceRespData, error) {
	// 创建请求对象
	req := larkapproval.NewGetInstanceReqBuilder().
		InstanceId(approvalInstanceCode).
		Build()

	// 发起请求
	resp, err := f.app.Approval.V4.Instance.Get(context.Background(), req)

	// 处理错误
	if err != nil {
		logs.Error("GetApprovalInstanceDetail is failed!:", err)
		return nil, err
	}

	// 服务端错误处理
	if !resp.Success() {
		logs.Error("logId: %s, error response: \n%s", resp.RequestId(), larkcore.Prettify(resp.CodeError))
		return nil, errors.New(larkcore.Prettify(resp.CodeError))
	}

	// 业务处理
	logs.Debug("GetApprovalInstanceDetail response: %s", larkcore.Prettify(resp))
	return resp.Data, nil
}
