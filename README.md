# 📦 GoHumanLoop WeWork

<div align="center">
<svg t="1750737000520" class="icon" viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg" p-id="1550" width="200" height="200"><path d="M683.17696 765.12a14.144 14.144 0 0 0 1.792 21.568c30.656 28.736 50.432 67.2 55.872 108.928a59.072 59.072 0 1 0 63.168-74.24 181.376 181.376 0 0 1-100.8-56.192 14.144 14.144 0 0 0-20.032-0.064z" fill="#FB6500" p-id="1551"></path><path d="M922.92096 671.616a58.624 58.624 0 0 0-16.96 35.52 181.12 181.12 0 0 1-56.064 100.928 14.144 14.144 0 1 0 21.504 18.048 181.12 181.12 0 0 1 108.928-55.872 59.072 59.072 0 1 0-57.28-98.688h-0.128z" fill="#0082EF" p-id="1552"></path><path d="M756.58496 505.088a59.072 59.072 0 0 0 35.584 100.416 181.12 181.12 0 0 1 100.928 56.064 14.144 14.144 0 1 0 18.048-21.504 181.44 181.44 0 0 1-55.872-108.928 59.072 59.072 0 0 0-98.688-26.048z" fill="#2DBC00" p-id="1553"></path><path d="M727.59296 597.376l-1.088 1.024a180.928 180.928 0 0 1-110.464 57.792 58.688 58.688 0 0 0-26.176 98.624 59.072 59.072 0 0 0 100.416-35.584 181.76 181.76 0 0 1 56.192-100.928 14.08 14.08 0 1 0-18.88-20.928z" fill="#FFCC00" p-id="1554"></path><path d="M370.79296 87.68c-105.856 11.648-201.856 56.896-270.848 127.616a359.488 359.488 0 0 0-66.112 92.992 323.84 323.84 0 0 0 22.784 327.04c18.752 28.288 49.472 63.616 77.632 88.768l-12.736 100.032-1.408 4.288c-0.384 1.216-0.384 2.624-0.512 3.904l-0.384 3.2 0.384 3.2a32.128 32.128 0 0 0 48.448 24.96h0.512l1.92-1.408 30.4-15.232 90.688-45.632a469.12 469.12 0 0 0 132.608 18.24 476.096 476.096 0 0 0 162.624-28.288 58.88 58.88 0 0 1-40.128-61.696 400.192 400.192 0 0 1-167.232 16.64l-8.96-1.344a399.296 399.296 0 0 1-60.096-12.544 41.28 41.28 0 0 0-32.192 3.328l-2.496 1.216-74.624 43.84-3.2 1.92c-1.792 1.088-2.624 1.408-3.52 1.408a5.184 5.184 0 0 1-4.8-5.312l2.816-11.52 3.328-12.544 5.312-20.672 6.208-22.976a31.296 31.296 0 0 0-11.328-34.816 327.36 327.36 0 0 1-75.328-78.464 255.296 255.296 0 0 1-18.304-257.664c13.44-26.88 31.104-51.776 53.056-74.24 56.576-58.368 136.128-95.488 224.192-105.024a421.312 421.312 0 0 1 91.584 0c87.488 10.048 166.72 47.744 222.912 105.728 21.76 22.464 39.424 47.744 52.48 74.624 17.472 35.712 26.368 73.536 26.368 112.256 0 4.096-0.384 8.128-0.512 12.032a58.944 58.944 0 0 1 72.512 8.512l2.624 3.2a322.176 322.176 0 0 0-32.192-167.616 361.728 361.728 0 0 0-65.408-92.992 443.648 443.648 0 0 0-269.76-128.512 491.52 491.52 0 0 0-109.312-0.448z" fill="#0082EF" p-id="1555"></path></svg>
</div>

GoHumanLoop WeWork 是针对`GoHumanLoop`在企业微信场景下进行审批、获取信息操作的示例服务。方便用户在使用`GohumanLoop`时，对接到自己的企业微信环境中。

---

GoHumanLoop 项目地址：https://github.com/ptonlix/gohumanloop

> **GoHumanLoop**: A Python library empowering AI agents to dynamically request human input (approval/feedback/conversation) at critical stages. Core features:
>
> - `Human-in-the-loop control`: Lets AI agent systems pause and escalate >decisions, enhancing safety and trust.
> - `Multi-channel integration`: Supports Terminal, Email, API, and frameworks like LangGraph/CrewAI (soon).
> - `Flexible workflows`: Combines automated reasoning with human oversight for reliable AI operations.
>
> Ensures responsible AI deployment by bridging autonomous agents and human judgment.
>
> Ensures responsible AI deployment by bridging autonomous agents and human judgment.

---

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
corpsecret = DPs0LAgVQzaQM9OhNQRmvaN2eeYqlfJlScMRWYyRrxY # 应用Secret
corpid = wwfb9fc76f89b0070b # 企业ID
ptoken = YZXxGRAii4o2hdU9G5c6v7 # Token
pkey = yEffDMYVX3L3NVLooTBXeJJNcfm7pwvQjaa8sfKC8oL # EncodingAESKey

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
