package unionpay

import (
	"encoding/json"
	"time"
)

const (
	// URL
	baseUrl        = "https://api-mop.chinaums.com"
	sandboxBaseUrl = "https://test-api-open.chinaums.com"

	RSA             = "RSA"
	RSA2            = "RSA2"
	formatJSON      = "json"
	charsetUTF8     = "UTF-8"
	timestampLayout = "2006-01-02 15:04:05"
)

var (
	DefaultLocation, _ = time.LoadLocation("Asia/Shanghai")
)

// 通用响应参数
type RspCommon struct {
	ResponseBizContent json.RawMessage `json:"response_biz_content"` // 响应参数集合,包含公共和业务参数
	Sign               string          `json:"sign"`                 // 针对返回参数集合的签名
}

// 响应参数
type RspBase struct {
	ErrCode           string `json:"errCode"`           // 平台错误码 0000：正常；
	ErrMsg            string `json:"errMsg"`            // 平台错误信息
	MsgId             string `json:"msgId"`             // 消息ID
	SrcReserve        string `json:"srcReserve"`        // 请求系统预留字段
	ResponseTimestamp string `json:"responseTimestamp"` // 报文响应时间 格式yyyy-MM-dd HH:mm:ss
}

// 小程序支付下单响应
type PayRsp struct {
	RspBase
	ConnectSys     string `json:"connectSys"`
	DelegatedFlag  string `json:"delegatedFlag"`
	MerName        string `json:"merName"`
	Mid            string `json:"mid"`
	SettleRefId    string `json:"settleRefId"`
	Tid            string `json:"tid"`
	TotalAmount    int    `json:"totalAmount"`
	TargetMid      string `json:"targetMid"`
	MiniPayRequest any    `json:"miniPayRequest"`
	TargetStatus   string `json:"targetStatus"`
	SeqId          string `json:"seqId"`
	MerOrderId     string `json:"merOrderId"`
	Status         string `json:"status"`
	TargetSys      string `json:"targetSys"`
}

// 退款响应
type RefundRsp struct {
	RspBase
	PayTime             string `json:"payTime"`
	ConnectSys          string `json:"connectSys"`
	MerName             string `json:"merName"`      // 商户名称
	Mid                 string `json:"mid"`          // 商户号，原样返回
	RefundStatus        string `json:"refundStatus"` // 退款状态详见取值说明
	SettleDate          string `json:"settleDate"`
	SendBackAmount      int    `json:"sendBackAmount"`
	Tid                 string `json:"tid"`                 // 终端号，原样返回
	RefundTargetOrderId string `json:"refundTargetOrderId"` // 目标系统退货订单号
	RefundFundsDesc     string `json:"refundFundsDesc"`     // 退款渠道描述
	RefundFunds         string `json:"refundFunds"`         // 退款渠道列表
	TargetMid           string `json:"targetMid"`           // 支付渠道商户号
	CardAttr            string `json:"cardAttr"`
	TargetStatus        string `json:"targetStatus"` // 目标平台状态
	SeqId               string `json:"seqId"`        // 平台流水号
	MerOrderId          string `json:"merOrderId"`   // 商户订单号
	TargetSys           string `json:"targetSys"`    // 目标平台代码
	BankInfo            string `json:"bankInfo"`
	DelegatedFlag       string `json:"delegatedFlag"`
	SettleRefId         string `json:"settleRefId"`
	RefundOrderId       string `json:"refundOrderId"`       // 退货订单号
	TotalAmount         int    `json:"totalAmount"`         // 支付总金额
	RefundInvoiceAmount int    `json:"refundInvoiceAmount"` // 实付部分退款金额
	ChnlCost            string `json:"chnlCost"`
	Status              string `json:"status"` // 交易状态
}
