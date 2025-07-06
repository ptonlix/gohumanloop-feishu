package controllers

import (
	"encoding/json"
	"strconv"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/ptonlix/gohumanloop-feishu/init/feishu"
	"github.com/ptonlix/gohumanloop-feishu/models"
	"github.com/ptonlix/gohumanloop-feishu/services"
)

// APIKeyController operations for API Keys
type APIKeyController struct {
	beego.Controller
}

// @Title Create API Key
// @Description Create a new API Key
// @Tags APIKey
// @Param	body		body 	models.APIKeyRequestData	true		"API Key creation data"
// @Success 200 {object} models.APIKeyResponseData
// @Failure 400 {object} models.APIResponse
// @router /apikey/create [post]
func (c *APIKeyController) CreateKey() {
	var ak models.APIKeyRequestData
	json.Unmarshal(c.Ctx.Input.RequestBody, &ak)
	logs.Info("incoming message: %s\n", ak)

	// 校验参数是否正确
	if ak.AppId != feishu.FeishuConf.AppId || ak.AppSecret != feishu.FeishuConf.AppSecret {
		c.Data["json"] = models.APIResponse{
			Success: false,
			Error:   "参数校验失败，生成API Key失败",
		}
		c.ServeJSON()
		return
	}

	// 生成API Key
	apiKey, err := services.APIKeyDataService.CreateAPIKey(ak.Name)
	if err != nil {
		c.Data["json"] = models.APIResponse{
			Success: false,
			Error:   "生成API Key失败:" + err.Error(),
		}
		c.ServeJSON()
		return
	}

	// 返回API Key
	c.Data["json"] = models.APIKeyResponseData{
		APIResponse: models.APIResponse{
			Success: true,
		},
		APIKey: *apiKey,
	}
	c.ServeJSON()
}

func (c *APIKeyController) GetAPIKeyByID() {
	id := c.GetString("id")
	// 将string类型的id转换为int64
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.Data["json"] = models.APIResponse{
			Success: false,
			Error:   "无效的ID参数",
		}
		c.ServeJSON()
		return
	}
	apiKey, err := services.APIKeyDataService.GetAPIKeyByID(idInt)
	if err != nil {
		c.Data["json"] = models.APIResponse{
			Success: false,
			Error:   "获取API Key失败:" + err.Error(),
		}
		c.ServeJSON()
		return
	}
	c.Data["json"] = models.APIKeyResponseData{
		APIResponse: models.APIResponse{
			Success: true,
		},
		APIKey: *apiKey,
	}
	c.ServeJSON()
}

// @Title Get API Key By Key
// @Description Get API Key details by actual key value
// @Tags APIKey
// @Param	key		query 	string	true		"The API key value"
// @Security ApiKeyAuth
// @Success 200 {object} models.APIKeyResponseData
// @Failure 400 {object} models.APIResponse
// @router /apikey/get [get]
func (c *APIKeyController) GetAPIKeyByKey() {
	key := c.GetString("key")
	apiKey, err := services.APIKeyDataService.GetAPIKeyByKey(key)
	if err != nil {
		c.Data["json"] = models.APIResponse{
			Success: false,
			Error:   "获取API Key失败:" + err.Error(),
		}
		c.ServeJSON()
		return
	}

	c.Data["json"] = models.APIKeyResponseData{
		APIResponse: models.APIResponse{
			Success: true,
		},
		APIKey: *apiKey,
	}
	c.ServeJSON()
}

// @Title Update API Key
// @Description Update an existing API Key
// @Tags APIKey
// @Param	body		body 	models.APIKey	true		"API Key update data"
// @Security ApiKeyAuth
// @Success 200 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @router /apikey/update [post]
func (c *APIKeyController) UpdateKey() {
	var ak models.APIKey
	json.Unmarshal(c.Ctx.Input.RequestBody, &ak)

	err := services.APIKeyDataService.UpdateAPIKey(&ak)
	if err != nil {
		c.Data["json"] = models.APIResponse{
			Success: false,
			Error:   "更新API Key失败",
		}
		c.ServeJSON()
		return
	}

	c.Data["json"] = models.APIResponse{
		Success: true,
	}
	c.ServeJSON()
}

// @Title Delete API Key
// @Description Delete an API Key
// @Tags APIKey
// @Param	body		body 	models.APIKey	true		"API Key to delete"
// @Security ApiKeyAuth
// @Success 200 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @router /apikey/delete [post]
func (c *APIKeyController) DeleteKey() {
	var ak models.APIKey
	json.Unmarshal(c.Ctx.Input.RequestBody, &ak)

	logs.Info("delete key: %v", ak)

	err := services.APIKeyDataService.DeleteAPIKey(ak.ID)
	if err != nil {
		c.Data["json"] = models.APIResponse{
			Success: false,
			Error:   "删除API Key失败",
		}
		c.ServeJSON()
		return
	}

	c.Data["json"] = models.APIResponse{
		Success: true,
	}
	c.ServeJSON()
}

// @Title List API Keys
// @Description Get all API Keys
// @Tags APIKey
// @Security ApiKeyAuth
// @Success 200 {object} models.APIKeyListResponse
// @Failure 400 {object} models.APIResponse
// @router /apikey/list  [get]
func (c *APIKeyController) ListKeys() {
	apiKeys, err := services.APIKeyDataService.ListAPIKeys()
	if err != nil {
		c.Data["json"] = models.APIResponse{
			Success: false,
			Error:   "获取API Key列表失败",
		}
		c.ServeJSON()
		return
	}

	c.Data["json"] = models.APIKeyListResponse{
		APIResponse: models.APIResponse{
			Success: true,
		},
		APIKeys: apiKeys,
	}
	c.ServeJSON()
}

// @Title Enable API Key
// @Description Enable a disabled API Key
// @Tags APIKey
// @Param	body		body 	models.APIKey	true		"API Key to enable"
// @Security ApiKeyAuth
// @Success 200 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @router /apikey/enable [post]
func (c *APIKeyController) EnableKey() {
	var ak models.APIKey
	json.Unmarshal(c.Ctx.Input.RequestBody, &ak)

	err := services.APIKeyDataService.EnableAPIKey(ak.ID)
	if err != nil {
		c.Data["json"] = models.APIResponse{
			Success: false,
			Error:   "启用API Key失败",
		}
		c.ServeJSON()
		return
	}

	c.Data["json"] = models.APIResponse{
		Success: true,
	}
	c.ServeJSON()
}

// @Title Disable API Key
// @Description Disable an enabled API Key
// @Tags APIKey
// @Param	body		body 	models.APIKey	true		"API Key to disable"
// @Security ApiKeyAuth
// @Success 200 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @router /apikey/disable [post]
func (c *APIKeyController) DisableKey() {
	var ak models.APIKey
	json.Unmarshal(c.Ctx.Input.RequestBody, &ak)

	err := services.APIKeyDataService.DisableAPIKey(ak.ID)
	if err != nil {
		c.Data["json"] = models.APIResponse{
			Success: false,
			Error:   "禁用API Key失败",
		}
		c.ServeJSON()
		return
	}

	c.Data["json"] = models.APIResponse{
		Success: true,
	}
	c.ServeJSON()
}
