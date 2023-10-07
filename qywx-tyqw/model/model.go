package model

// ========================= 服务端接收到消息的格式====================

// MsgContent 抽取通用字段,进行消息匹配处理.
type MsgContent struct {
	ToUsername   string `xml:"ToUserName"`
	FromUsername string `xml:"FromUserName"`
	CreateTime   int32  `xml:"CreateTime"`
	MsgType      string `xml:"MsgType"`
	MsgId        int64  `xml:"MsgId"`
	AgentId      int32  `xml:"AgentID"`
}

// 想代码优化...
type MsgTextContent struct {
	MsgContent
	Content string `xml:"Content"`
}

type MsgImageContent struct {
	ToUsername   string `xml:"ToUserName"`
	FromUsername string `xml:"FromUserName"`
	CreateTime   int32  `xml:"CreateTime"`
	MsgType      string `xml:"MsgType"` //image
	PicUrl       string `xml:"PicUrl"`
	MediaId      string `xml:"MediaId"`
	MsgId        int64  `xml:"MsgId"`
	AgentId      int32  `xml:"AgentID"`
}

type MsgVoiceContent struct {
	ToUsername   string `xml:"ToUserName"`
	FromUsername string `xml:"FromUserName"`
	CreateTime   int32  `xml:"CreateTime"`
	MsgType      string `xml:"MsgType"` //voice
	MediaId      string `xml:"MediaId"`
	Format       string `xml:"Format"`
	MsgId        int64  `xml:"MsgId"`
	AgentId      int32  `xml:"AgentID"`
}

type MsgVideoContent struct {
	ToUsername   string `xml:"ToUserName"`
	FromUsername string `xml:"FromUserName"`
	CreateTime   int32  `xml:"CreateTime"`
	MsgType      string `xml:"MsgType"` //video
	MediaId      string `xml:"MediaId"`
	ThumbMediaId string `xml:"ThumbMediaId"`
	MsgId        int64  `xml:"MsgId"`
	AgentId      int32  `xml:"AgentID"`
}

type MsgMapContent struct {
	ToUsername   string  `xml:"ToUserName"`
	FromUsername string  `xml:"FromUserName"`
	CreateTime   int32   `xml:"CreateTime"`
	MsgType      string  `xml:"MsgType"` //location
	Location_X   float64 `xml:"Location_X"`
	Location_Y   float64 `xml:"Location_Y"`
	Scale        int32   `xml:"Scale"`
	Label        string  `xml:"Label"`
	AppType      string  `xml:"AppType"`
	MsgId        int64   `xml:"MsgId"`
	AgentId      int32   `xml:"AgentID"`
}

type MsgLinkContent struct {
	ToUsername   string `xml:"ToUserName"`
	FromUsername string `xml:"FromUserName"`
	CreateTime   int32  `xml:"CreateTime"`
	MsgType      string `xml:"MsgType"` //link
	Title        string `xml:"Title"`
	Description  string `xml:"Description"`
	Url          string `xml:"Url"`
	PicUrl       string `xml:"PicUrl"`
	MsgId        int64  `xml:"MsgId"`
	AgentId      int32  `xml:"AgentID"`
}

// ====================== 主动推送消息的格式 ====================

type SendAppMsgText struct {
	ToUser                 string            `json:"touser"`
	ToParty                string            `json:"toparty"`
	ToTag                  string            `json:"totag"`
	MsgType                string            `json:"msgtype"`
	AgentId                int32             `json:"agentid"`
	Text                   AppTextMsgContent `json:"text"`
	Safe                   int32             `json:"safe"`
	EnableIdTrans          int32             `json:"enable_id_trans"`
	EnableDuplicateCheck   int32             `json:"enable_duplicate_check"`
	DuplicateCheckInterval int32             `json:"duplicate_check_interval"`
}

type AppTextMsgContent struct {
	Content string `json:"content"`
}

type QywxCallBackIPResp struct {
	ErrCode  int32    `json:"errcode"`
	ErrMsg   string   `json:"errmsg"`
	WhiteIPS []string `json:"ip_list"`
}
