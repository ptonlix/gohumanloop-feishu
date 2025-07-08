# 📦 GoHumanLoop FeiShu

<div align="center">
	<img height=80 src="https://lf-package-cn.feishucdn.com/obj/feishu-static/lark/open/website/images/899fa60e60151c73aaea2e25871102dc.svg"><br>
</div>

**GoHumanLoop FeiShu** 是针对`GoHumanLoop`在飞书场景下进行审批、获取信息等人机协同操作的示例服务。方便用户在使用`GohumanLoop`时，对接到自己的企业微信环境中。

`GoHumanLoop` 项目地址：https://github.com/ptonlix/gohumanloop ✈️

> [!NOTE] > **GoHumanLoop**: A Python library empowering AI agents to dynamically request human input (approval/feedback/conversation) at critical stages. Core features:
>
> - `Human-in-the-loop control`: Lets AI agent systems pause and escalate >decisions, enhancing safety and trust.
> - `Multi-channel integration`: Supports Terminal, Email, API, and frameworks like LangGraph/CrewAI (soon).
> - `Flexible workflows`: Combines automated reasoning with human oversight for reliable AI operations.
>
> Ensures responsible AI deployment by bridging autonomous agents and human judgment.

## 💻 项目部署

> [!IMPORTANT]
> 需要用户提前准备好飞书企业自建应用
> 详情见：https://open.feishu.cn/
>
> 1. 用户需要获取飞书`企业ID`
> 2. 用户需要在飞书中创建应用`GoHumanLoop`，获取`应用ID`和`应用Secret`
> 3. 用户需要在应用中,开启 API 接收消息事件。（使用使用长连接接收回调）
> 4. 用户需要将`应用ID`和`应用Secret`配置到 GoHumanLoop 中
> 5. 用户需要在飞书中创建审批模板，获取`审批模板ID`和`信息模板ID`
> 6. 用户需要将`审批模板ID`、`信息模板ID`、`创建人OpenId`、`审批人OpenId`配置到 GoHumanLoop 中

### 配置文件

- 项目配置样例文件在`conf/app.conf.example`中

```yaml
appname = gohumanloop-feishu
httpport = 9800 # HTTP 端口按需配置

# wework
appid = XXX # 企业应用ID
appsecret = XXX # 应用Secret

# template
approve_template_id = B102B40F-E288-4A48-8220-FB0E69990222 # 审批模板ID
info_template_id = 4B131588-A5E2-4CE6-91D1-5A208CA62EDA
creator_userid = ou_ed1d616f3023f7781808ce479b5f0fdd # 创建人ID，用户在应用中的OpenId
approver_userid = ou_ed1d616f3023f7781808ce479b5f0fdd  # 审批人ID，用户在应用中的OpenId

# database
datapath = ./data/gohumanloop.db # 数据库路径
```

- 修改配置文件

```
mv conf/app.conf.example conf/app.conf
```

### 飞书审批模板

目前这个版本中，支持审批和信息获取。分别使用两个模板，模板格式固定，需要参考以下配置：

<div align="center">
	<img height=240 src="http://cdn.oyster-iot.cloud/202507062313394.png"><br>
    <b face="雅黑">审批模板</b>
</div>

- 参考图片内的字段，都是文本控件和多行文本控件。包括以下字段
  1. 任务 ID
  2. 对话 ID
  3. 请求 ID
  4. HumanLoop 类型
  5. 申请内容
  6. 申请问题
  7. 申请说明

以上字段由 GoHumanLoop 库来传输并自动填充并自动发起审批流程

<div align="center">
	<img height=240 src="http://cdn.oyster-iot.cloud/202507062313050.png"><br>
    <b face="雅黑">审批模板</b>
</div>

- 审批流程可以参考上图设置，审批人设置为自选

<div align="center">
	<img height=240 src="http://cdn.oyster-iot.cloud/202507062314889.png"><br>
    <b face="雅黑">信息获取模板</b>
</div>

- 参考图片内的字段，都是文本控件和多行文本控件。详情同审批流程模板和说明

<div align="center">
	<img height=240 src="http://cdn.oyster-iot.cloud/202507062315370.png"><br>
    <b face="雅黑">信息获取模板</b>
</div>

- 信息获取流程不需要具体审批，只需要获取具体信息,设置必须填写回执信息，否则无法获取信息。

### 部署方式

GoHumanLoop FeiShu 支持两种部署方式手动部署和 Docker 部署。

> [!WARNING]
> 通过采用飞书 SDK 的 Websocket 长连接来接收回调消息，需要在飞书应用中开启长连接接收消息事件,支持本地部署。区别于企业微信的 Webhook 方案需要注册同企微主体下的服务器。服务器和域名需要已备案等这些繁琐的操作

#### 1. 手动部署

Go 版本要求：1.23.0

- 下载代码

```shell
git clone https://github.com/ptonlix/gohumanloop-feishu.git
```

- 编译

```shell
make build
```

## 运行

```
./gohumanloop-feishu
```

#### 2. Docker 部署

- 提前安装好 Docker 服务

```
docker pull ptonlix/gohumanloop-feishu:latest
```

- 运行容器

```
docker run -d \
  --name gohumanloop-feishu \
  -v /path/to/local/conf:/app/conf \
  -v /path/to/local/data:/app/data \
  -p 9800:9800 \
  ptonlix/gohumanloop-feishu:latest
```

## 📖 项目介绍

### 架构设计

<div align="center">
	<img height=240 src="http://cdn.oyster-iot.cloud/202507090033672.png"><br>
    <b face="雅黑">GoHumanLoop与gohumanloop-feishu架构关系</b>
</div>

- `GoHumanLoop`提供了一套统一的 API 接口，通过`API Provider`对外提供。
- `gohumanloop-feishu`实现了`API Consumer`的功能，通过`API Provider`来获取审批相关的信息，并且通过飞书 SDK 实现了与用户的飞书进行交互，发送审批请求和获取审批事件回调等人机交互协同的操作。

### 实现介绍

`gohumanloop-feishu`采用[Beego](https://github.com/beego/beego)作为 Web 框架。`sqlite`作为简单的数据存储。[飞书 SDK](https://open.feishu.cn/document/server-side-sdk/golang-sdk-guide/preparations)作为飞书 API 实现。提供一个可拓展的 GoHumanLoop 飞书审批示例服务。

- 访问 Swagger 文档:

```
go run main.go
```

```
http://127.0.0.1:9800/docs
```

## 🤝 参与贡献

GoHumanLoop FeiShu 和文档均开源，我们欢迎以问题、文档和 PR 等形式做出贡献。

## 📱 联系方式

<img height=300 src="http://cdn.oyster-iot.cloud/202505231802103.png"/>

🎉 如果你对本项目感兴趣，欢迎扫码联系作者交流

## 🌟 Star History

[![Star History Chart](https://api.star-history.com/svg?repos=gohumanloop-feishu&type=Date)](https://www.star-history.com/#gohumanloop-feishu&Date)
