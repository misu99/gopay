package icbc

import (
	"encoding/json"
	"time"
)

const (
	// URL
	baseUrl        = "https://gw.open.icbc.com.cn"
	sandboxBaseUrl = "https://syb-test.allinpay.com/apiweb"

	RSA             = "RSA"
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
	ReturnCode           int    `json:"return_code"`
	ReturnMsg            string `json:"return_msg"`
	MsgId                string `json:"msg_id"`
	ThirdPartyReturnCode string `json:"third_party_return_code"` // 第三方报错时返回的报错码
	ThirdPartyReturnMsg  string `json:"third_party_return_msg"`  // 第三方报错时返回的报错信息
}

// PayRsp 通用支付响应
type PayRsp struct {
	RspBase
	TotalAmt         string `json:"total_amt"`          // 订单总金额
	OutTradeNo       string `json:"out_trade_no"`       // 商户系统订单号，原样返回
	OrderId          string `json:"order_id"`           // 工行订单号
	PayTime          string `json:"pay_time"`           // 支付完成时间，格式为：yyyyMMdd
	MerId            string `json:"mer_id"`             // 商户编号
	PayMode          string `json:"pay_mode"`           // 支付方式
	AccessType       string `json:"access_type"`        // 收单接入方式
	CardKind         string `json:"card_kind"`          // 卡种
	TradeType        string `json:"trade_type"`         // 支付方式为微信时返回，交易类型
	WxDataPackage    string `json:"wx_data_package"`    // 支付方式为微信时返回，微信数据包，用于之后唤起微信支付
	ZfbDataPackage   string `json:"zfb_data_package"`   // 支付方式为支付宝时返回，支付宝数据包，用于之后唤起支付宝支付
	UnionDataPackage string `json:"union_data_package"` // 支付方式为云闪付时返回，云闪付受理订单号，用于之后进行银联云闪付支付
}

// RefundRsp 退款响应
type RefundRsp struct {
	RspBase
	IntrxSerialNo string `json:"intrx_serial_no"` // 退货工行流水号
}
