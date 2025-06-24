package models

// APIResponse 基础 API 响应模型
type APIResponse struct {
	Success bool   `json:"success"`         // 请求是否成功
	Error   string `json:"error,omitempty"` // 错误信息（可选）
}
