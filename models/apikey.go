package models

// HumanLoopRequestData 人机协作请求数据
type APIKeyRequestData struct {
	Agentid    string `json:"agentid"`    // 企业微信应用ID
	Corpsecret string `json:"corpsecret"` // 企业微信应用密钥
	Name       string `json:"name"`       //密钥名称
}

type APIKeyResponseData struct {
	APIResponse
	APIKey
}

type APIKeyUpdateData struct {
	ID     int64  `json:"id"`     // 密钥ID
	Name   string `json:"name"`   // 密钥名称
	Status bool   `json:"status"` // 密钥状态
}

type APIKeyDeleteData struct {
	ID int64 `json:"id"` // 密钥ID
}

type APIKeyListResponse struct {
	APIResponse
	APIKeys []*APIKey `json:"api_keys"`
}
