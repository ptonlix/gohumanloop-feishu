package feishu

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

type FeishuConfig struct {
	TemplateConfig
	AppId     string
	AppSecret string
}

var FeishuConf FeishuConfig

func init() {

	var err error

	// 加载并验证配置
	if FeishuConf.AppId, err = beego.AppConfig.String("appid"); err != nil || FeishuConf.AppId == "" {
		logs.Error("飞书配置加载失败: AppID为空或无效")
		panic("飞书配置加载失败: AppID不能为空")
	}

	if FeishuConf.AppSecret, err = beego.AppConfig.String("appsecret"); err != nil || FeishuConf.AppSecret == "" {
		logs.Error("飞书配置加载失败: AppSecret为空或无效")
		panic("飞书配置加载失败: AppSecret不能为空")
	}

	if FeishuConf.ApproveTemplateId, err = beego.AppConfig.String("approve_template_id"); err != nil || FeishuConf.ApproveTemplateId == "" {
		logs.Error("飞书配置加载失败: 审批模板ID为空或无效")
		panic("飞书配置加载失败: 审批模板ID不能为空")
	}

	if FeishuConf.InfoTemplateId, err = beego.AppConfig.String("info_template_id"); err != nil || FeishuConf.InfoTemplateId == "" {
		logs.Error("飞书配置加载失败: 信息模板ID为空或无效")
		panic("飞书配置加载失败: 信息模板ID不能为空")
	}

	if FeishuConf.CreatorUserId, err = beego.AppConfig.String("creator_userid"); err != nil || FeishuConf.CreatorUserId == "" {
		logs.Error("飞书配置加载失败: 创建者用户ID为空或无效")
		panic("飞书配置加载失败: 创建者用户ID不能为空")
	}

	if FeishuConf.ApproverUserId, err = beego.AppConfig.String("approver_userid"); err != nil || FeishuConf.ApproverUserId == "" {
		logs.Error("飞书配置加载失败: 审批者用户ID为空或无效")
		panic("飞书配置加载失败: 审批者用户ID不能为空")
	}
}
