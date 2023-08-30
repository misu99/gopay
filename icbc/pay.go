package icbc

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/misu99/gopay"
	"github.com/misu99/gopay/pkg/util"
)

// Pay 聚合支付接口
func (c *Client) Pay(ctx context.Context, bm gopay.BodyMap) (rsp *PayRsp, err error) {
	err = bm.CheckEmptyError("mer_id", "out_trade_no", "pay_mode", "access_type")
	if err != nil {
		return nil, err
	}

	var bs []byte
	if bs, err = c.doPost(ctx, payPath, bm); err != nil {
		return nil, err
	}

	rspCommon := new(RspCommon)
	if err = json.Unmarshal(bs, rspCommon); err != nil {
		return nil, fmt.Errorf("[%w], bytes: %s", gopay.UnmarshalErr, string(bs))
	}

	rsp = new(PayRsp)
	if err = json.Unmarshal(rspCommon.ResponseBizContent, rsp); err != nil {
		return nil, fmt.Errorf("[%w], bytes: %s", gopay.UnmarshalErr, string(bs))
	}

	if err := bizErrCheck(rsp.RspBase); err != nil {
		return nil, err
	}

	return rsp, c.verifySign(rspCommon)
}

// Refund 统一退款接口
func (c *Client) Refund(ctx context.Context, bm gopay.BodyMap) (rsp *RefundRsp, err error) {
	err = bm.CheckEmptyError("mer_id", "outtrx_serial_no", "ret_total_amt")
	if err != nil {
		return nil, err
	}

	if bm.GetString("order_id") == util.NULL && bm.GetString("out_trade_no") == util.NULL {
		return nil, fmt.Errorf("[%w], %v", gopay.MissParamErr, "order_id和out_trade_no必填其一")
	}

	var bs []byte
	if bs, err = c.doPost(ctx, refundPath, bm); err != nil {
		return nil, err
	}

	rspCommon := new(RspCommon)
	if err = json.Unmarshal(bs, rspCommon); err != nil {
		return nil, fmt.Errorf("[%w], bytes: %s", gopay.UnmarshalErr, string(bs))
	}

	rsp = new(RefundRsp)
	if err = json.Unmarshal(rspCommon.ResponseBizContent, rsp); err != nil {
		return nil, fmt.Errorf("[%w], bytes: %s", gopay.UnmarshalErr, string(bs))
	}

	if err := bizErrCheck(rsp.RspBase); err != nil {
		return nil, err
	}

	return rsp, c.verifySign(rspCommon)
}
