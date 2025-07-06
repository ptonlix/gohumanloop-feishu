# 📦 GoHumanLoop FeiShu

<div align="center">
	<img height=160 src="http://cdn.oyster-iot.cloud/企业微信-copy.png"><br>
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

## 💻 项目部署

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

- 项目配置样例文件在`conf/app.conf.example`中

```yaml
appname = gohumanloop-feishu
httpport = 9800 # HTTP 端口按需配置

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

- 修改配置文件

```
mv conf/app.conf.example conf/app.conf
```

### 企业微信审批模板

目前这个版本中，支持审批和信息获取。分别使用两个模板，模板格式固定，需要参考以下配置：

<div align="center">
	<img height=240 src="http://cdn.oyster-iot.cloud/202506241753570.png"><br>
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
	<img height=240 src="http://cdn.oyster-iot.cloud/202506241756802.png"><br>
    <b face="雅黑">审批模板</b>
</div>

- 审批流程可以参考上图设置，审批人设置为自选

<div align="center">
	<img height=240 src="http://cdn.oyster-iot.cloud/202506241800810.png"><br>
    <b face="雅黑">信息获取模板</b>
</div>

- 参考图片内的字段，都是文本控件和多行文本控件。详情同审批流程模板和说明

<div align="center">
	<img height=240 src="http://cdn.oyster-iot.cloud/202506242226055.png"><br>
    <b face="雅黑">信息获取模板</b>
</div>

- 信息获取流程不需要具体审批，只需要获取具体信息。没有设置审批人，只设置了办理人，专用于获取信息。

### 部署方式

GoHumanLoop Wework 支持两种部署方式手动部署和 Docker 部署。

> [!WARNING]
> 这两种方式均需要有企业微信同一注册主体下的服务器。服务器和域名需要已备案，开启 API 接收消息时也需要域名验证是否是同一注册主体下

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

#### 配置反向代理

以 Nginx 为例，可以参考在 Nginx 配置文件中添加以下路由配置

```shell
location ^~ /humanloop/ {
        proxy_pass http://127.0.0.1:9800/gohumanloop/callback;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Real-Port $remote_port;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
}

location ^~ /api/v1/humanloop/ {
        proxy_pass http://127.0.0.1:9800;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Real-Port $remote_port;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
}

location ^~ /api/v1/apikey/ {
        proxy_pass http://127.0.0.1:9800;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Real-Port $remote_port;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
}
```

## 📖 项目介绍

### 架构设计

<div align="center">
	<img height=240 src="http://cdn.oyster-iot.cloud/202506252306729.png"><br>
    <b face="雅黑">GoHumanLoop与gohumanloop-feishu架构关系</b>
</div>

- `GoHumanLoop`提供了一套统一的 API 接口，通过`API Provider`对外提供。
- `gohumanloop-feishu`实现了`API Consumer`的功能，通过`API Provider`来获取审批相关的信息，并且通过企业微信 WeWork API 实现了与用户的企业微信应用进行交互，发送审批请求和获取审批事件回调等。

### 实现介绍

`gohumanloop-feishu`采用[Beego](https://github.com/beego/beego)作为 Web 框架。`sqlite`作为简单的数据存储。[go-workwx](https://github.com/xen0n/go-workwx)作为企业微信 API 实现。提供一个可拓展的 GoHumanLoop 企业微信审批示例服务。

- 访问 Swagger 文档:

```
go run main.go
```

```
http://127.0.0.1:9800/docs
```

## 🤝 参与贡献

GoHumanLoop Wework 和文档均开源，我们欢迎以问题、文档和 PR 等形式做出贡献。

## 📱 联系方式

<img height=300 src="http://cdn.oyster-iot.cloud/202505231802103.png"/>

🎉 如果你对本项目感兴趣，欢迎扫码联系作者交流

## 🌟 Star History

[![Star History Chart](https://api.star-history.com/svg?repos=gohumanloop-feishu&type=Date)](https://www.star-history.com/#gohumanloop-feishu&Date)
