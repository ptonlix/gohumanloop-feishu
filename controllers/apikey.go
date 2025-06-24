package controllers

import (
	"encoding/json"
	"strconv"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/ptonlix/gohumanloop-wework/init/wework"
	"github.com/ptonlix/gohumanloop-wework/models"
	"github.com/ptonlix/gohumanloop-wework/services"
)

type APIKeyController struct {
	beego.Controller
}

func (c *APIKeyController) CreateKey() {
	var ak models.APIKeyRequestData
	json.Unmarshal(c.Ctx.Input.RequestBody, &ak)
	logs.Info("incoming message: %s\n", ak)

	// 校验参数是否正确
	if ak.Agentid != strconv.Itoa(int(wework.WeworkConf.AgentId)) || ak.Corpsecret != wework.WeworkConf.CorpSecret {
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
			Error:   "生成API Key失败",
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

func (c *APIKeyController) EnableKey() {
	id, _ := c.GetInt64("id")
	err := services.APIKeyDataService.EnableAPIKey(id)
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

func (c *APIKeyController) DisableKey() {
	id, _ := c.GetInt64("id")
	err := services.APIKeyDataService.DisableAPIKey(id)
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
