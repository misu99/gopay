package icbc

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/misu99/gopay"
	"github.com/misu99/gopay/pkg/util"
)

// 微信小程序支付下单
func (c *Client) MiniWechatPay(ctx context.Context, bm gopay.BodyMap) (rsp *PayRsp, err error) {
	err = bm.CheckEmptyError("merOrderId", "totalAmount", "tradeType", "subOpenId", "notifyUrl")
	if err != nil {
		return nil, err
	}

	var bs []byte
	if bs, err = c.doPost(ctx, miniWechatPayOrderPath, bm); err != nil {
		return nil, err
	}

	rsp = new(PayRsp)
	if err = json.Unmarshal(bs, rsp); err != nil {
		return nil, fmt.Errorf("[%w], bytes: %s", err, string(bs))
	}

	if err := bizErrCheck(rsp.RspBase); err != nil {
		return nil, err
	}

	return rsp, nil
}

// h5微信支付下单
func (c *Client) H5WechatPay(ctx context.Context, bm gopay.BodyMap) (rsp *PayRsp, err error) {
	err = bm.CheckEmptyError("merOrderId", "instMid", "totalAmount", "notifyUrl")
	if err != nil {
		return nil, err
	}

	c.doGet(ctx, h5WechatPayOrderPath, bm)

	//var bs []byte
	//if bs, err = c.doGet(ctx, H5WechatPayOrderPath, bm); err != nil {
	//	return nil, err
	//}

	return rsp, nil
}

// h5支付宝支付下单
func (c *Client) H5AliPay(ctx context.Context, bm gopay.BodyMap) (rsp *PayRsp, err error) {
	err = bm.CheckEmptyError("merOrderId", "instMid", "totalAmount", "notifyUrl")
	if err != nil {
		return nil, err
	}

	c.doGet(ctx, h5AliPayOrderPath, bm)

	//var bs []byte
	//if bs, err = c.doGet(ctx, H5WechatPayOrderPath, bm); err != nil {
	//	return nil, err
	//}

	return rsp, nil
}

// Refund 统一退款接口
func (c *Client) Refund(ctx context.Context, bm gopay.BodyMap) (rsp *RefundRsp, err error) {
	err = bm.CheckEmptyError("merOrderId", "refundAmount", "refundOrderId")
	if err != nil {
		return nil, err
	}

	if bm.GetString("merOrderId") == util.NULL && bm.GetString("targetOrderId") == util.NULL {
		return nil, fmt.Errorf("[%w], %v", gopay.MissParamErr, "merOrderId和targetOrderId必填其一")
	}

	var bs []byte
	if bs, err = c.doPost(ctx, refundOrderPath, bm); err != nil {
		return nil, err
	}

	rsp = new(RefundRsp)
	if err = json.Unmarshal(bs, rsp); err != nil {
		return nil, fmt.Errorf("[%w], bytes: %s", gopay.UnmarshalErr, string(bs))
	}

	if err := bizErrCheck(rsp.RspBase); err != nil {
		return nil, err
	}

	return rsp, nil
}
