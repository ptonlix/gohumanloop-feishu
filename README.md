# 📦 GoHumanLoop WeWork

<div align="center">
	<img height=160 src="http://cdn.oyster-iot.cloud/企业微信-copy.png"><br>
    <b face="雅黑">审批模板</b>
</div>

**GoHumanLoop WeWork** 是针对`GoHumanLoop`在企业微信场景下进行审批、获取信息操作的示例服务。方便用户在使用`GohumanLoop`时，对接到自己的企业微信环境中。

`GoHumanLoop` 项目地址：https://github.com/ptonlix/gohumanloop ✈️

> [!NOTE] > **GoHumanLoop**: A Python library empowering AI agents to dynamically request human input (approval/feedback/conversation) at critical stages. Core features:
>
> - `Human-in-the-loop control`: Lets AI agent systems pause and escalate >decisions, enhancing safety and trust.
> - `Multi-channel integration`: Supports Terminal, Email, API, and frameworks like LangGraph/CrewAI (soon).
> - `Flexible workflows`: Combines automated reasoning with human oversight for reliable AI operations.
>
> Ensures responsible AI deployment by bridging autonomous agents and human judgment.
>
> Ensures responsible AI deployment by bridging autonomous agents and human judgment.

## 项目部署

GoHumanLoop Wework 支持两种部署方式

> [!IMPORTANT]
> 需要用户提前准备好企业微信和企业微信应用
> 详情见：https://work.weixin.qq.com/
>
> 1. 用户需要获取企业微信`企业ID`
> 2. 用户需要在企业微信中创建应用`GoHumanLoop`，获取`应用ID`和`应用Secret`
> 3. 用户需要在应用中,开启 API 接收消息。获取`Token`和`EncodingAESKey` 用于接收审批消息事件。审批应用中设置可调用接口的应用，关联新创建的应用。
> 4. 用户需要将`企业ID`、`应用ID`、`应用Secret`、`Token`、`EncodingAESKey`配置到 GoHumanLoop 中
> 5. 用户需要在企业微信中创建审批模板，获取`审批模板ID`和`信息模板ID`
> 6. 用户需要将`审批模板ID`、`信息模板ID`、`创建人ID`、`审批人ID`配置到 GoHumanLoop 中

### 配置文件

```yaml
appname = gohumanloop-wework
httpport = 8080 # HTTP 端口按需配置

# wework
agentid = 1000003 # 企业微信应用ID
corpsecret = XXXXX # 应用Secret
corpid = XXXXX # 企业ID
ptoken = XXXXX # Token
pkey = XXXXX # EncodingAESKey

# template
approve_template_id = 8TmoaR5xEaZsuzKyRT4Zt82FLYCYXVN5EVk6R # 审批模板ID
info_template_id = 3WN63LowuwFRsDXft1GbiQi4NrYyLApeejYCBs3S # 信息模板ID
creator_userid = ChenFuDong # 创建人ID，详情参考企业微信文档
approver_userid = ChenFuDong # 审批人ID (默认审批人，实际可通过 GoHumanLoop Metadata数据指定)

# database
datapath = ./data/gohumanloop.db # 数据库路径
```

### 企业微信审批模板

目前这个版本中，支持审批和信息获取。分别使用两个模板，模板格式固定，需要参考以下配置：

<div align="center">
	<img height=240 src="http://cdn.oyster-iot.cloud/202506241753570.png"><br>
    <b face="雅黑">审批模板</b>
</div>

- 参考图片内的字段，都是文本控件和多行文本控件

<div align="center">
	<img height=240 src="http://cdn.oyster-iot.cloud/202506241756802.png"><br>
    <b face="雅黑">审批模板</b>
</div>

- 审批流程可以参考上图设置，审批人设置为自选

<div align="center">
	<img height=240 src="http://cdn.oyster-iot.cloud/202506241800810.png"><br>
    <b face="雅黑">信息获取模板</b>
</div>

- 参考图片内的字段，都是文本控件和多行文本控件

<div align="center">
	<img height=240 src="http://cdn.oyster-iot.cloud/202506242226055.png"><br>
    <b face="雅黑">信息获取模板</b>
</div>

- 信息获取流程不需要具体审批，只需要获取具体信息。没有设置审批人，只设置了办理人，专用于获取信息。

### 部署方式

GoHumanLoop Wework 支持两种部署方式

#### 1. 本地部署

#### 2. Docker 部署

## 项目介绍
