package icbc

import (
	"encoding/json"
	"time"
)

const (
	// URL
	baseUrl = "https://gw.open.icbc.com.cn"
	//sandboxBaseUrl = "https://APIPCS3.dccnet.com.cn"
	sandboxBaseUrl = "https://apipcs3.dccnet.com.cn:447" // 443端口沙箱地址会出现拒绝连接，应使用447端口

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
	ReturnCode           string `json:"return_code"`
	ReturnMsg            string `json:"return_msg"`
	MsgId                string `json:"msg_id"`
	ThirdPartyReturnCode string `json:"third_party_return_code"` // 第三方报错时返回的报错码
	ThirdPartyReturnMsg  string `json:"third_party_return_msg"`  // 第三方报错时返回的报错信息
}

// PayRsp 通用支付响应参数
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

// PayQueryRsp 通用支付查询响应参数
type PayQueryRsp struct {
	RspBase
	PayStatus             string `json:"pay_status"`
	CardNo                string `json:"card_no"`
	MerId                 string `json:"mer_id"`
	TotalAmt              string `json:"total_amt"`
	PointAmt              string `json:"point_amt"`
	EcouponAmt            string `json:"ecoupon_amt"`
	MerDiscAmt            string `json:"mer_disc_amt"`
	CouponAmt             string `json:"coupon_amt"`
	BankDiscAmt           string `json:"bank_disc_amt"`
	PaymentAmt            string `json:"payment_amt"`
	OutTradeNo            string `json:"out_trade_no"`
	OrderId               string `json:"order_id"`
	PayTime               string `json:"pay_time"`
	TotalDiscAmt          string `json:"total_disc_amt"`
	Attach                string `json:"attach"`
	ThirdTradeNo          string `json:"third_trade_no"`
	CardFlag              string `json:"card_flag"`
	DecrFlag              string `json:"decr_flag"`
	OpenId                string `json:"open_id"`
	PayType               string `json:"pay_type"`
	AccessType            string `json:"access_type"`
	CardKind              string `json:"card_kind"`
	ThirdPartyCouponAmt   string `json:"third_party_coupon_amt"`
	ThirdPartyDiscountAmt string `json:"third_party_discount_amt"`
	OutTradeNoWx          string `json:"out_trade_no_wx"`
	ReturnCodeWx          string `json:"return_code_wx"`
	Appid                 string `json:"appid"`
	BankType              string `json:"bank_type"`
	CashFee               string `json:"cash_fee"`
	FeeType               string `json:"fee_type"`
	IsSubscribe           string `json:"is_subscribe"`
	MchId                 string `json:"mch_id"`
	NonceStr              string `json:"nonce_str"`
	ResultCode            string `json:"result_code"`
	Sign                  string `json:"sign"`
	TimeEnd               string `json:"time_end"`
	TotalFee              string `json:"total_fee"`
	TradeType             string `json:"trade_type"`
	TransactionId         string `json:"transaction_id"`
	SignType              string `json:"sign_type"`
	Body                  string `json:"body"`
	SpbillCreateIp        string `json:"spbill_create_ip"`
	NotifyUrl             string `json:"notify_url"`
	UnionDiscountAmt      string `json:"union_discount_amt"`
	UnionMchtDiscountAmt  string `json:"union_mcht_discount_amt"`
}

// RefundRsp 退款响应
type RefundRsp struct {
	RspBase
	IntrxSerialNo string `json:"intrx_serial_no"` // 退货工行流水号
}
