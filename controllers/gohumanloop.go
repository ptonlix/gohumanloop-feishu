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

// @Title Create HumanLoop Request
// @Description Create a new HumanLoop request
// @Tags HumanLoop
// @Param	body		body 	models.HumanLoopRequestData	true		"body for HumanLoop request"
// @Security ApiKeyAuth
// @Success 200 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @router /api/v1/humanloop/request [post]
func (c *GoHumanLoopController) Request() {
	var ob models.HumanLoopRequestData
	json.Unmarshal(c.Ctx.Input.RequestBody, &ob)
	logs.Info("incoming message: %s\n", ob)

	// 发送给服务处理
	spNo, err := services.GHLWeWorkService.HumanLoopRequestOA(ob)
	if err != nil {
		logs.Error("HumanLoopRequestOA failed: %v", err)
		c.Data["json"] = models.APIResponse{
			Success: false,
			Error:   err.Error(),
		}
		c.ServeJSON()
		return
	}
	logs.Info("CreateHumanLoop spNo: %v", spNo)

	hl := &models.HumanLoop{
		TaskId:         ob.TaskId,
		RequestId:      ob.RequestId,
		ConversationId: ob.ConversationId,
		SpNo:           spNo,
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
	logs.Info("CreateHumanLoop success: %v", id)

	// 响应
	c.Data["json"] = models.APIResponse{
		Success: true,
	}
	c.ServeJSON()
}

// @Title Get HumanLoop Status
// @Description Get the status of a HumanLoop request
// @Tags HumanLoop
// @Param	conversation_id		query	string	true	"Conversation ID"
// @Param	request_id		query	string	true	"Request ID"
// @Param	platform		query	string	true	"Platform"
// @Security ApiKeyAuth
// @Success 200 {object} models.HumanLoopStatusResponse
// @Failure 400 {object} models.APIResponse
// @router /api/v1/humanloop/status [get]
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

// @Title Continue HumanLoop
// @Description Continue a HumanLoop conversation (not supported)
// @Tags HumanLoop
// @Param	body		body 	models.HumanLoopContinueData	true		"body for continuing HumanLoop"
// @Security ApiKeyAuth
// @Success 200 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @router /api/v1/humanloop/continue [post]
func (c *GoHumanLoopController) Continue() {
	var ob models.HumanLoopContinueData
	json.Unmarshal(c.Ctx.Input.RequestBody, &ob)

	//TODO 对话模式暂不支持
	logs.Error("ContinueHumanLoop not supported")
	c.Data["json"] = models.APIResponse{
		Success: false,
		Error:   "ContinueHumanLoop not supported",
	}
	c.ServeJSON()
	return

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

// @Title Cancel HumanLoop
// @Description Cancel a specific HumanLoop request
// @Tags HumanLoop
// @Param	body		body 	models.HumanLoopCancelData	true		"body for canceling HumanLoop"
// @Security ApiKeyAuth
// @Success 200 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @router /humanloop/cancel [post]
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

// @Title Cancel Conversation
// @Description Cancel all HumanLoop requests in a conversation
// @Tags HumanLoop
// @Param	body		body 	models.HumanLoopCancelConversationData	true		"body for canceling conversation"
// @Security ApiKeyAuth
// @Success 200 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @router /humanloop/cancel_conversation [post]
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
