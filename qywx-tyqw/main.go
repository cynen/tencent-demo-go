package main

import (
	"encoding/json"
	"encoding/xml"
	"flag"
	"github.com/gin-gonic/gin"
	"github.com/sbzhu/weworkapi_golang/wxbizmsgcrypt"
	"log"
	"math/rand"
	"model"
	"net/http"
	"sync"
	"time"
	"tongyi"
	"tool"
)

var (
	//全局变量,注意协程互相竞争的情况.
	GlobalTyUrl              string     // 通义的地址
	GlobalTyApiKey           string     // 通义的api-key
	GlobalQywxToken          string     // 验证自建应用的token
	GlobalQywxEncodingAseKey string     // 验证自建应用的秘钥
	GlobalQywxReceiverId     string     // 自建应用时,是corpid
	GlobalQywxCorpSecret     string     // 企业的秘钥
	GlobalQywxAcessToken     string     // 访问企业微信接口的token
	lock                     sync.Mutex // 排异锁.
	WhiteIPList              []string   // ip白名单
)

func main() {

	// 第一步,读取配置参数.
	configfile := flag.String("c", "config.yml", "config file")
	port := flag.String("p", "8888", "Server Port")
	flag.Parse()

	// 读取配置文件.
	appconfig := tool.ReadYamlFile(*configfile)
	log.Println("通义千问配置项: ", appconfig.Tongyi)
	log.Println("企业微信配置项: ", appconfig.QywxConfig)

	// 1.设置全局变量
	setConfig(appconfig)

	// 定时器,实现更新token的目的
	go func() {
		// 在协程中实现死循环
		for {
			// 设置全局的AccessToken
			UpdateAccessToken()
			//GlobalQywxAcessToken = tool.CheckAndGetAccessToken(GlobalQywxAcessToken, GlobalQywxReceiverId, GlobalQywxCorpSecret)
			time.Sleep(time.Second * 3600) // 3600秒更新一次token
		}
	}()

	go func() {
		for {
			UpdateWhiteList()
			// 24小时更新一次
			time.Sleep(time.Hour * 24)
		}
	}()

	// 创建企业微信的加解密工具实例
	wxcpt := wxbizmsgcrypt.NewWXBizMsgCrypt(GlobalQywxToken, GlobalQywxEncodingAseKey, GlobalQywxReceiverId, wxbizmsgcrypt.XmlType)

	// ===== 创建web服务 ====
	engine := gin.Default()
	//gin.SetMode(gin.ReleaseMode)
	// 做简单防护
	engine.GET("/", func(context *gin.Context) {
		context.Status(http.StatusOK)
		context.String(200, "Hello~~")
	})

	// 1.服务器的校验.
	engine.GET("/qywxpush", func(context *gin.Context) {
		// 如何做好防护?
		if !tool.IsWhiteIp(context.ClientIP(), WhiteIPList) {
			log.Println("恶意IP:", context.ClientIP())
			return
		}
		// 解析出url上的参数值如下：
		verifyMsgSign := context.Query("msg_signature")
		verifyTimestamp := context.Query("timestamp")
		verifyNonce := context.Query("nonce")
		verifyEchoStr := context.Query("echostr")

		log.Println("接收到微信的验证推送消息======")
		// 校验,不能为空,有一个为空就需要报错.
		if tool.IsEmpty(verifyMsgSign) || tool.IsEmpty(verifyTimestamp) || tool.IsEmpty(verifyNonce) || tool.IsEmpty(verifyEchoStr) {
			log.Println("verifyMsgSign: ", verifyMsgSign)
			log.Println("verifyTimestamp: ", verifyTimestamp)
			log.Println("verifyNonce: ", verifyNonce)
			log.Println("verifyEchoStr: ", verifyEchoStr)
			log.Println("参数获取异常,存在空参数,请检查")
			return
		}
		// 处理数据
		echoStr, crptErr := wxcpt.VerifyURL(verifyMsgSign, verifyTimestamp, verifyNonce, verifyEchoStr)
		if crptErr != nil {
			log.Println("VerifyUrl Err", crptErr)
		}
		// 需要立即返回. 1秒内.
		context.Writer.Write(echoStr)
	})

	// 2.处理用户消息.
	engine.POST("/qywxpush", func(context *gin.Context) {
		// 安全防护
		if !tool.IsWhiteIp(context.ClientIP(), WhiteIPList) {
			log.Println("恶意IP:", context.ClientIP())
			return
		}

		// 解析出url上的参数值如下：
		verifyMsgSign := context.Query("msg_signature")
		verifyTimestamp := context.Query("timestamp")
		verifyNonce := context.Query("nonce")
		// 参数校验,不能为空. 否则异常.
		if tool.IsEmpty(verifyMsgSign) || tool.IsEmpty(verifyTimestamp) || tool.IsEmpty(verifyNonce) {
			log.Println("verifyMsgSign: ", verifyMsgSign)
			log.Println("verifyTimestamp: ", verifyTimestamp)
			log.Println("verifyNonce: ", verifyNonce)
			log.Println("参数获取异常,存在空参数")
			return
		}
		// 获取post的请求数据 body .
		reqData, err := context.GetRawData()
		if err != nil {
			log.Println("获取reqData失败 \n", err)
			return
		}
		// 1.获取用户发送的消息.
		receviedUserMsg, crptErr := wxcpt.DecryptMsg(verifyMsgSign, verifyTimestamp, verifyNonce, reqData)
		if crptErr != nil {
			log.Println("解密数据出错 \n", crptErr)
			return
		}
		log.Println("服务器接收到的解密后的数据: \n", string(receviedUserMsg))
		// ==== 解密完成 ====================

		// 可能超过5秒. 所以需要单独处理. 如果超过5秒,服务器会再发一次请求.
		// 新开线程去处理查询通义接口.
		go DealWithReceiveDMsg(receviedUserMsg)

		// 立即给微信服务器返回消息.文档要求.
		// 这里是需要返回 200状态码,以及 空串.
		context.Status(200)
		context.Writer.Write([]byte(""))
	})
	// 启动web服务.
	log.Fatal(engine.Run(":" + *port))
}

// 设置全局变量
func setConfig(appconfig tool.AppConfig) {
	GlobalTyUrl = appconfig.Tongyi.Url
	GlobalTyApiKey = appconfig.Tongyi.ApiKey
	GlobalQywxToken = appconfig.QywxConfig.Token
	GlobalQywxEncodingAseKey = appconfig.QywxConfig.EncodingAseKey
	GlobalQywxReceiverId = appconfig.QywxConfig.ReceiverId
	GlobalQywxCorpSecret = appconfig.QywxConfig.CorpSecret
}
func UpdateAccessToken() {
	lock.Lock()
	defer lock.Unlock()
	// 设置全局的AccessToken
	log.Println("准备更新 AccessToken")
	GlobalQywxAcessToken = tool.CheckAndGetAccessToken(GlobalQywxAcessToken, GlobalQywxReceiverId, GlobalQywxCorpSecret)
	log.Println("AccessToken 更新成功")

}

func UpdateWhiteList() {
	log.Println("准备更新企业微信回调IP的白名单列表")
	log.Println("更新前的企业微信回调IP的列表:", WhiteIPList)
	url := "https://qyapi.weixin.qq.com/cgi-bin/getcallbackip?access_token="
	UpdateAccessToken()
	url = url + GlobalQywxAcessToken
	result, err := tool.HttpGetRequest(url)
	if err != nil {
		log.Println("获取企业微信的回调IP异常,", err)
		return
	}
	var qywxBackIPResp model.QywxCallBackIPResp
	jsonData, _ := json.Marshal(result)
	json.Unmarshal(jsonData, &qywxBackIPResp)
	if qywxBackIPResp.ErrCode != 0 {
		log.Println("获取企业微信的回调IP异常,", qywxBackIPResp.ErrMsg)
		return
	}
	WhiteIPList = qywxBackIPResp.WhiteIPS
	log.Println("更新企业微信回调IP的列表成功")
	log.Println("更新后的企业微信回调IP的列表:", WhiteIPList)
}

// DealWithReceiveDMsg 通过获取用户的消息,处理逻辑. 注意,这个是新开的线程在处理.
// 这些都是被动回复的. 返回的消息,需要是按照要求拼接好的xml结构体.
func DealWithReceiveDMsg(msg []byte) {
	// 1.判断msgType
	var msgContent model.MsgContent
	err := xml.Unmarshal(msg, &msgContent)
	if err != nil {
		log.Println("Unmarshal Failed..", err)
		return
	}
	log.Println("解析后的消息结构体struct: ", msgContent)
	//log.Println("MsgType: ", msgContent.MsgType, ",用户: ", msgContent.FromUsername)

	// 2.根据msgType类型分开处理.
	switch msgContent.MsgType {
	case "text":
		// 处理文本消息.调用通义千问.
		var msgTextContent model.MsgTextContent
		err := xml.Unmarshal(msg, &msgTextContent)
		if err != nil {
			log.Println("Unmarshal Failed..", err)
			return
		}
		// 3.调用通义接口,获取返回结果.
		tyResult := GetTongyiResp(msgTextContent.Content)
		//tyResult := GetTextResp(msgContent, GetTongyiResp(msgTextContent.Content))
		// 4.主动推送消息.
		PushMsgToApp(tyResult, msgTextContent.FromUsername, msgTextContent.AgentId)
	case "image":
		// TODO
		PushMsgToApp("暂不支持的类型", msgContent.FromUsername, msgContent.AgentId)
	case "voice":
		// TODO
		PushMsgToApp("暂不支持的类型", msgContent.FromUsername, msgContent.AgentId)
	case "video":
		// TODO

	case "location":
		// TODO

	case "link":
		// TODO

	default:
		// TODO
		PushMsgToApp("暂不支持的类型", msgContent.FromUsername, msgContent.AgentId)
	}
}

// 主动推送消息给企业微信的应用
func PushMsgToApp(msg string, touser string, agentid int32) {
	log.Println("准备推送消息给企业微信应用,touser: ", touser, ", agentid:", agentid)
	url := "https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token="
	// 1.获取accessToken,定时器处理.
	//GlobalQywxAcessToken = tool.CheckAndGetAccessToken(GlobalQywxAcessToken, GlobalQywxReceiverId, GlobalQywxCorpSecret)

	// 2.创建推送消息体
	data := model.SendAppMsgText{
		ToUser:  touser,
		MsgType: "text",
		AgentId: agentid,
		Text: model.AppTextMsgContent{
			Content: msg,
		},
		Safe: 0,
	}
	header := make(map[string]string)
	//3.推送消息
	url = url + GlobalQywxAcessToken
	// 发完不管结果.
	tool.HttpPostJson(url, data, header)
}

// GetTongyiResp 调用通义千问,访问接口
// 返回的是 通义大模型的返回结果.
func GetTongyiResp(msg string) string {
	log.Println("准备调用通义接口")
	// 构造请求体.
	body := tongyi.ReqBody{
		Model: "qwen-turbo",
		Input: tongyi.Input{
			Prompt: msg,
		},
		Parameters: tongyi.Parameters{
			Temperature:  0.1,
			TopP:         0.5,
			TopK:         10,
			Seed:         rand.Int31n(65536), //给随机种子
			ResultFormat: "message",
		},
	}

	header := make(map[string]string)
	header["Authorization"] = "Bearer " + GlobalTyApiKey
	ty_result, err := tool.HttpPostJson(GlobalTyUrl, body, header)
	if err != nil {
		log.Println("调用通义接口报错:", err)
		return "调用通义千问接口异常,请联系管理员"
	}
	log.Println("通义千问返回结果: ", ty_result)

	// 对通义千问的结果进行处理
	jsonData, _ := json.Marshal(ty_result)
	var TyResp tongyi.Resp
	err = json.Unmarshal(jsonData, &TyResp)
	if err != nil {
		log.Println("Error Unmarshal data", err)
		return ""
	}
	// 需要判断数据的暂停原因:
	// TODO

	return TyResp.Output.Choices[0].Message.Content
	// 返回通义千问的响应结果文本.
	//return TyResp.Output.Text
}
