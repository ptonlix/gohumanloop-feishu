package controllers

import (
	"encoding/json"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/ptonlix/gohumanloop-wework/models"
	"github.com/ptonlix/gohumanloop-wework/services"

	"github.com/beego/beego/v2/core/logs"
)

type GoHumanLoopController struct {
	beego.Controller
}

func (c *GoHumanLoopController) Request() {
	var ob models.HumanLoopRequestData
	json.Unmarshal(c.Ctx.Input.RequestBody, &ob)
	logs.Info("incoming message: %s\n", ob)

	hl := &models.HumanLoop{
		TaskId:         ob.TaskId,
		RequestId:      ob.RequestId,
		ConversationId: ob.ConversationId,
		LoopType:       ob.LoopType,
		ContextMap:     ob.Context,
		Platform:       ob.Platform,
		MetadataMap:    ob.Metadata,
		Status:         models.HumanLoopStatusPending,
	}
	hl.BeforeSave()

	// 写入数据库
	id, err := services.GHLDataService.CreateHumanLoop(hl)
	if err != nil {
		logs.Error("CreateHumanLoop failed: %v", err)
		c.Data["json"] = models.APIResponse{
			Success: false,
			Error:   err.Error(),
		}
		c.ServeJSON()
		return
	}

	// 发送给服务处理
	logs.Info("CreateHumanLoop success: %v", id)

	// 响应
	c.Data["json"] = models.APIResponse{
		Success: true,
	}
	c.ServeJSON()
}

func (c *GoHumanLoopController) Status() {
	conversationId := c.GetString("conversation_id")
	requestId := c.GetString("request_id")
	platform := c.GetString("platform")
	ob := models.HumanLoopStatusParams{ConversationId: conversationId, RequestId: requestId, Platform: platform}
	logs.Info("incoming message: %v\n", ob)

	// 从数据库获取
	hl, err := services.GHLDataService.GetHumanLoopByRequestId(ob.ConversationId, ob.RequestId, ob.Platform)
	if err != nil {
		logs.Error("GetHumanLoopByID failed: %v", err)
		c.Data["json"] = models.APIResponse{
			Success: false,
			Error:   err.Error(),
		}
		c.ServeJSON()
		return
	}

	// 构造响应
	resp := models.HumanLoopStatusResponse{
		APIResponse: models.APIResponse{
			Success: true,
		},
		Status:      hl.Status,
		Response:    hl.Response,
		Feedback:    hl.Feedback,
		RespondedBy: hl.RespondedBy,
		RespondedAt: hl.RespondedAt,
	}

	// 响应
	c.Data["json"] = resp
	c.ServeJSON()
}

func (c *GoHumanLoopController) Continue() {
	var ob models.HumanLoopContinueData
	json.Unmarshal(c.Ctx.Input.RequestBody, &ob)

	hl := &models.HumanLoop{
		TaskId:         ob.TaskId,
		RequestId:      ob.RequestId,
		ConversationId: ob.ConversationId,
		LoopType:       "conversation",
		ContextMap:     ob.Context,
		Platform:       ob.Platform,
		MetadataMap:    ob.Metadata,
	}
	hl.BeforeSave()

	// 写入数据库
	id, err := services.GHLDataService.CreateHumanLoop(hl)
	if err != nil {
		logs.Error("ContinueHumanLoop failed: %v", err)
		c.Data["json"] = models.APIResponse{
			Success: false,
			Error:   err.Error(),
		}
		c.ServeJSON()
		return
	}

	// 发送给服务处理
	logs.Info("ContinueHumanLoop success: %v", id)
	// 响应
	c.Data["json"] = models.APIResponse{
		Success: true,
	}
	c.ServeJSON()

}

func (c *GoHumanLoopController) Cancel() {
	var ob models.HumanLoopCancelData
	json.Unmarshal(c.Ctx.Input.RequestBody, &ob)

	// 从数据库获取
	hl, err := services.GHLDataService.GetHumanLoopByRequestId(ob.ConversationId, ob.RequestId, ob.Platform)
	if err != nil {
		logs.Error("GetHumanLoopByID failed: %v", err)
		c.Data["json"] = models.APIResponse{
			Success: false,
			Error:   err.Error(),
		}
		c.ServeJSON()
		return
	}

	hl.Status = models.HumanLoopStatusCancelled

	// 更新数据库
	rows, err := services.GHLDataService.UpdateHumanLoop(hl)
	if err != nil {
		logs.Error("UpdateHumanLoop failed: %v", err)
		c.Data["json"] = models.APIResponse{
			Success: false,
			Error:   err.Error(),
		}
		c.ServeJSON()
		return
	}

	logs.Info("CancelHumanLoop success: %v", rows)

	// 响应
	c.Data["json"] = models.APIResponse{
		Success: true,
	}
	c.ServeJSON()

}

func (c *GoHumanLoopController) CancelConversation() {
	var ob models.HumanLoopCancelConversationData
	json.Unmarshal(c.Ctx.Input.RequestBody, &ob)

	// 从数据库获取
	hls, err := services.GHLDataService.GetHumanLoopsByConversationId(ob.ConversationId, ob.Platform)
	if err != nil {
		logs.Error("GetHumanLoopByConversationId failed: %v", err)
		c.Data["json"] = models.APIResponse{
			Success: false,
			Error:   err.Error(),
		}
		c.ServeJSON()
		return
	}

	for _, hl := range hls {
		hl.Status = models.HumanLoopStatusCancelled
	}
	// 批量更新

	err = services.GHLDataService.BatchUpdateHumanLoops(hls)
	if err != nil {
		logs.Error("GetHumanLoopByConversationId failed: %v", err)
		c.Data["json"] = models.APIResponse{
			Success: false,
			Error:   err.Error(),
		}
		c.ServeJSON()
		return
	}

	logs.Info("CancelConversation success: %v", len(hls))

	// 响应
	c.Data["json"] = models.APIResponse{
		Success: true,
	}
	c.ServeJSON()
}
