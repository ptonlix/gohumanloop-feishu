package wework

import (
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
)

type TemplateConfig struct {
	ApproveTemplateId string
	InfoTemplateId    string
	CreatorUserId     string
	ApproverUserId    string
}

type WeworkConfig struct {
	TemplateConfig
	AgentId    int
	CorpId     string
	CorpSecret string
	PKey       string
	PToken     string
}

var WeworkConf WeworkConfig

func init() {
	// 获取ChatGpt配置
	var err error

	// 加载并验证配置
	if WeworkConf.AgentId, err = beego.AppConfig.Int("agentid"); err != nil {
		logs.Error("Wework配置加载失败: agentid")
		panic("Wework配置加载失败: agentid 必须为有效的整数")
	}

	if WeworkConf.CorpId, err = beego.AppConfig.String("corpid"); err != nil || WeworkConf.CorpId == "" {
		logs.Error("Wework配置加载失败: corpid为空或无效")
		panic("Wework配置加载失败: corpid 不能为空")
	}

	if WeworkConf.CorpSecret, err = beego.AppConfig.String("corpsecret"); err != nil || WeworkConf.CorpSecret == "" {
		logs.Error("Wework配置加载失败: corpsecret为空或无效")
		panic("Wework配置加载失败: corpsecret 不能为空")
	}

	if WeworkConf.PKey, err = beego.AppConfig.String("pkey"); err != nil || WeworkConf.PKey == "" {
		logs.Error("Wework配置加载失败: pkey为空或无效")
		panic("Wework配置加载失败: pkey 不能为空")
	}

	if WeworkConf.PToken, err = beego.AppConfig.String("ptoken"); err != nil || WeworkConf.PToken == "" {
		logs.Error("Wework配置加载失败: ptoken为空或无效")
		panic("Wework配置加载失败: ptoken 不能为空")
	}

	if WeworkConf.ApproveTemplateId, err = beego.AppConfig.String("approve_template_id"); err != nil || WeworkConf.ApproveTemplateId == "" {
		logs.Error("Wework配置加载失败: approve_template_id为空或无效")
		panic("Wework配置加载失败: approve_template_id 不能为空")
	}

	if WeworkConf.InfoTemplateId, err = beego.AppConfig.String("info_template_id"); err != nil || WeworkConf.InfoTemplateId == "" {
		logs.Error("Wework配置加载失败: info_template_id为空或无效")
		panic("Wework配置加载失败: info_template_id 不能为空")
	}

	if WeworkConf.CreatorUserId, err = beego.AppConfig.String("creator_userid"); err != nil || WeworkConf.CreatorUserId == "" {
		logs.Error("Wework配置加载失败: creator_userid为空或无效")
		panic("Wework配置加载失败: creator_userid 不能为空")
	}

	if WeworkConf.ApproverUserId, err = beego.AppConfig.String("approver_userid"); err != nil || WeworkConf.ApproverUserId == "" {
		logs.Error("Wework配置加载失败: approver_userid为空或无效")
		panic("Wework配置加载失败: approver_userid 不能为空")
	}
}
