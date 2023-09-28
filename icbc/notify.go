package icbc

import (
	"encoding/json"
	"fmt"
	"github.com/misu99/gopay"
	"strconv"
	"time"
)

// 异步回调通知响应体参数
type NotifyReq struct {
	RspBase
	CardNo       string `json:"card_no"`
	MerId        string `json:"mer_id"`
	TotalAmt     string `json:"total_amt"`
	PointAmt     string `json:"point_amt"`
	EcouponAmt   string `json:"ecoupon_amt"`
	MerDiscAmt   string `json:"mer_disc_amt"`
	CouponAmt    string `json:"coupon_amt"`
	BankDiscAmt  string `json:"bank_disc_amt"`
	PaymentAmt   string `json:"payment_amt"`
	OutTradeNo   string `json:"out_trade_no"`
	OrderId      string `json:"order_id"`
	PayTime      string `json:"pay_time"`
	TotalDiscAmt string `json:"total_disc_amt"`
	Attach       string `json:"attach"`
	ThirdTradeNo string `json:"third_trade_no"`
	CardFlag     string `json:"card_flag"`
	DecrFlag     string `json:"decr_flag"`
	OpenId       string `json:"open_id"`
	PayType      string `json:"pay_type"`
	AccessType   string `json:"access_type"`
	CardKind     string `json:"card_kind"`
}

// 异步回调通用返回参数
type NotifyRspCommon struct {
	ResponseBizContent json.RawMessage `json:"response_biz_content"` // 响应业务参数
	SignType           string          `json:"sign_type"`            // 签名类型
	Sign               string          `json:"sign"`                 // 针对返回参数集合的签名
}

// 异步回调通知返回响应体参数
type NotifyRsp struct {
	ReturnCode int    `json:"return_code"`
	ReturnMsg  string `json:"return_msg"`
	MsgId      string `json:"msg_id"`
}

// 获取回调通知返回
func (c *Client) GetNotifyRsp(code int, msg string) (rsp NotifyRspCommon, err error) {
	msgID := strconv.FormatInt(time.Now().UnixNano(), 10)

	// 业务响应参数
	bm := make(gopay.BodyMap)
	bm.Set("return_code", code)
	bm.Set("return_msg", msg)
	bm.Set("msg_id", msgID)
	bizContent, err := json.Marshal(bm)
	if err != nil {
		return
	}

	// 通用响应参数
	params := make(gopay.BodyMap)
	params.Set("response_biz_content", string(bizContent)).
		Set("sign_type", RSA2)

	// 计算参数签名
	sign, err := c.getRsaSign("", params, params.GetString("sign_type"), c.privateKey)
	if err != nil {
		return rsp, fmt.Errorf("GetRsaSign Error: %w", err)
	}

	return NotifyRspCommon{ResponseBizContent: bizContent, SignType: RSA2, Sign: sign}, nil
}
