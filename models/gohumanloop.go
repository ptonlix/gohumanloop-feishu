package models

import (
	"time"
)

const (
	HumanLoopTypeApprove      = "approve"
	HumanLoopTypeInformation  = "information"
	HumanLoopTypeConversation = "conversation"
)

const (
	// 等待处理状态
	HumanLoopStatusPending = "pending"
	// 已批准状态
	HumanLoopStatusApproved = "approved"
	// 已拒绝状态
	HumanLoopStatusRejected = "rejected"
	// 已过期状态
	HumanLoopStatusExpired = "expired"
	// 错误状态
	HumanLoopStatusError = "error"
	// 已完成状态（用于非审批类交互）
	HumanLoopStatusCompleted = "completed"
	// 进行中状态（用于多轮对话的中间状态）
	HumanLoopStatusInProgress = "inprogress"
	// 已取消状态
	HumanLoopStatusCancelled = "cancelled"
)

// APIResponse 基础 API 响应模型
type APIResponse struct {
	Success bool   `json:"success"`         // 请求是否成功
	Error   string `json:"error,omitempty"` // 错误信息（可选）
}

// HumanLoopRequestData 人机协作请求数据
type HumanLoopRequestData struct {
	TaskId         string         `json:"task_id"`            // 任务标识符
	ConversationId string         `json:"conversation_id"`    // 会话标识符
	RequestId      string         `json:"request_id"`         // 请求标识符
	LoopType       string         `json:"loop_type"`          // 循环类型
	Context        map[string]any `json:"context"`            // 提供给人类的上下文信息
	Platform       string         `json:"platform"`           // 使用平台（如微信、飞书）
	Metadata       map[string]any `json:"metadata,omitempty"` // 附加元数据（可选）
}

// HumanLoopStatusParams 人机协作状态查询参数
type HumanLoopStatusParams struct {
	ConversationId string `json:"conversation_id"` // 会话标识符
	RequestId      string `json:"request_id"`      // 请求标识符
	Platform       string `json:"platform"`        // 使用平台
}

// HumanLoopStatusResponse 人机协作状态响应
type HumanLoopStatusResponse struct {
	APIResponse           // 内嵌基础响应
	Status      string    `json:"status"`                 // 请求状态（默认："pending"）
	Response    any       `json:"response,omitempty"`     // 人类响应数据（可选）
	Feedback    any       `json:"feedback,omitempty"`     // 反馈数据（可选）
	RespondedBy string    `json:"responded_by,omitempty"` // 响应者信息（可选）
	RespondedAt time.Time `json:"responded_at,omitempty"` // 响应时间戳（可选）
}

// HumanLoopCancelData 取消单个人机请求
type HumanLoopCancelData struct {
	ConversationId string `json:"conversation_id"` // 会话标识符
	RequestId      string `json:"request_id"`      // 请求标识符
	Platform       string `json:"platform"`        // 使用平台
}

// HumanLoopCancelConversationData 取消整个会话
type HumanLoopCancelConversationData struct {
	ConversationId string `json:"conversation_id"` // 会话标识符
	Platform       string `json:"platform"`        // 使用平台
}

// HumanLoopContinueData 继续人机协作
type HumanLoopContinueData struct {
	ConversationId string         `json:"conversation_id"`    // 会话标识符
	RequestId      string         `json:"request_id"`         // 请求标识符
	TaskId         string         `json:"task_id"`            // 任务标识符
	Context        map[string]any `json:"context"`            // 上下文信息
	Platform       string         `json:"platform"`           // 使用平台
	Metadata       map[string]any `json:"metadata,omitempty"` // 附加元数据（可选）
}
