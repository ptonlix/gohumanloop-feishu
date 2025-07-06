package controllers

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	"github.com/beego/beego/v2/core/logs"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	larkevent "github.com/larksuite/oapi-sdk-go/v3/event"
	"github.com/larksuite/oapi-sdk-go/v3/event/dispatcher"
	larkapproval "github.com/larksuite/oapi-sdk-go/v3/service/approval/v4"
	larkws "github.com/larksuite/oapi-sdk-go/v3/ws"
	"github.com/ptonlix/gohumanloop-feishu/init/feishu"
	"github.com/ptonlix/gohumanloop-feishu/models"
	"github.com/ptonlix/gohumanloop-feishu/services"
)

type FeishuController struct {
	client *larkws.Client
}

func NewFeishuController() *FeishuController {
	// 注册事件 Register event
	eventHandler := dispatcher.NewEventDispatcher("", "").
		OnCustomizedEvent("approval_instance", func(ctx context.Context, event *larkevent.EventReq) error {
			// logs.Info("[ OnCustomizedEvent access ], type: message, data: %s\n", string(event.Body))
			// 将 event.Body 转换为 字典对象
			var eventData map[string]interface{}
			err := json.Unmarshal(event.Body, &eventData)
			if err != nil {
				logs.Error("解析飞书事件失败: %v", err)
				return err
			}
			logs.Info("飞书事件解析结果: %v", eventData)
			// 从 eventData 中提取审批实例ID
			// 先断言event字段
			eventMap, ok := eventData["event"].(map[string]interface{})
			if !ok {
				logs.Error("飞书事件中缺少event字段")
				return errors.New("飞书事件中缺少event字段")
			}

			// 再从event中获取instance_code
			approvalInstanceId, ok := eventMap["instance_code"].(string)
			if !ok {
				logs.Error("飞书事件中缺少审批实例ID或类型不正确")
				return errors.New("飞书事件中缺少审批实例ID或类型不正确")
			}
			approvalStatus, ok := eventMap["status"].(string)
			if !ok {
				logs.Error("飞书事件中缺少审批实例状态或类型不正确")
				return errors.New("飞书事件中缺少审批实例状态或类型不正确")
			}

			// 确保审批实例ID不为空
			if approvalInstanceId == "" {
				logs.Error("飞书事件中审批实例ID为空")
				return errors.New("飞书事件中审批实例ID为空")
			}

			// 获取审批实例详情
			approvalInstanceDetail, err := services.GHLFeishuService.GetApprovalInstanceDetail(approvalInstanceId)
			if err != nil {
				logs.Error("获取审批实例详情失败: %v", err)
				return err
			}
			// 检查Timeline切片是否为空
			if len(approvalInstanceDetail.Timeline) == 0 {
				logs.Error("Timeline为空")
				return errors.New("Timeline为空")
			}
			// 获取最后一条Timeline的Comment，如果为nil则设置为空字符串
			var responseStr *string
			lastComment := approvalInstanceDetail.Timeline[len(approvalInstanceDetail.Timeline)-1].Comment
			if lastComment == nil {
				emptyStr := ""
				responseStr = &emptyStr
			} else {
				responseStr = lastComment
			}
			responseBy := approvalInstanceDetail.Timeline[len(approvalInstanceDetail.Timeline)-1].OpenId

			feedback := ""
			// 检查CommentList切片是否为空
			if len(approvalInstanceDetail.CommentList) != 0 {
				feedback = *approvalInstanceDetail.CommentList[len(approvalInstanceDetail.CommentList)-1].Comment
			}
			approvalStatus = strings.ToLower(approvalStatus)
			if approvalStatus == models.HumanLoopStatusPending {
				// 审批中，不更新数据库
				return nil
			} else if approvalStatus == models.HumanLoopStatusApproved || approvalStatus == models.HumanLoopStatusRejected {
				// 获取数据库审批记录
				var hl *models.HumanLoop
				hl, err = services.GHLDataService.GetHumanLoopByBySpNo(approvalInstanceId)
				if err != nil {
					logs.Error("获取数据库审批记录失败: %v", err)
					return err
				}

				if hl.LoopType == models.HumanLoopTypeInformation {
					approvalStatus = models.HumanLoopStatusCompleted
				}
			} else {
				approvalStatus = models.HumanLoopStatusCancelled
			}
			// 更新数据库审批实例状态
			_, err = services.GHLDataService.UpdateHumanLoopStatusBySpNo(approvalInstanceId, approvalStatus, *responseStr, *responseBy, feedback)
			if err != nil {
				logs.Error("更新审批实例状态失败: %v", err)
				return err
			}
			return nil
		}).
		OnCustomizedEvent("approval_task", func(ctx context.Context, event *larkevent.EventReq) error {
			return nil
		}).
		OnCustomizedEvent("approval", func(ctx context.Context, event *larkevent.EventReq) error {
			return nil
		}).
		OnCustomizedEvent("approval_cc", func(ctx context.Context, event *larkevent.EventReq) error {
			return nil
		}).
		OnP2ApprovalUpdatedV4(func(ctx context.Context, event *larkapproval.P2ApprovalUpdatedV4) error {
			return nil
		})

	// 构建 client Build client
	cli := larkws.NewClient(feishu.FeishuConf.AppId, feishu.FeishuConf.AppSecret,
		larkws.WithEventHandler(eventHandler),
		larkws.WithLogLevel(larkcore.LogLevelDebug),
	)

	return &FeishuController{
		client: cli,
	}
}

// StartWebSocket 启动WebSocket长连接，在独立的goroutine中运行
func (f *FeishuController) StartWebSocket() {
	go func() {
		// 建立长连接
		err := f.client.Start(context.Background())
		if err != nil {
			panic(err)
		}
	}()
}
