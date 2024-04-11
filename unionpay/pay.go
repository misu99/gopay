package unionpay

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/misu99/gopay"
)

// 微信小程序支付下单
func (c *Client) MiniWechatPay(ctx context.Context, bm gopay.BodyMap) (rsp *PayRsp, err error) {
	err = bm.CheckEmptyError("merOrderId", "totalAmount", "subOpenId", "notifyUrl")
	if err != nil {
		return nil, err
	}

	// 补充必填业务参数
	bm.
		Set("instMid", "MINIDEFAULT"). // 业务类型
		Set("tradeType", "MINI")       // 交易类型

	var bs []byte
	if bs, err = c.doPost(ctx, miniWechatPayOrderPath, bm); err != nil {
		return nil, err
	}

	rsp = new(PayRsp)
	if err = json.Unmarshal(bs, rsp); err != nil {
		return nil, fmt.Errorf("[%w], bytes: %s", err, string(bs))
	}

	if err = bizErrCheck(rsp.RspBase); err != nil {
		return nil, err
	}

	return rsp, nil
}

// 公众号支付下单
func (c *Client) Webpay(ctx context.Context, bm gopay.BodyMap) (res string, err error) {
	err = bm.CheckEmptyError("merOrderId", "totalAmount", "notifyUrl")
	if err != nil {
		return res, err
	}

	// 补充必填业务参数
	bm.
		Set("instMid", "YUEDANDEFAULT") // 业务类型

	if res, err = c.doUrl(WebPayOrderPath, bm); err != nil {
		return res, err
	}

	return res, nil
}

// h5微信支付下单
func (c *Client) H5WechatPay(ctx context.Context, bm gopay.BodyMap) (rsp *PayRsp, err error) {
	err = bm.CheckEmptyError("merOrderId", "totalAmount", "notifyUrl")
	if err != nil {
		return nil, err
	}

	// 补充必填业务参数
	bm.
		Set("instMid", "H5DEFAULT") // 业务类型

	var bs []byte
	if bs, err = c.doGet(ctx, h5WechatPayOrderPath, bm); err != nil {
		return nil, err
	}

	rsp = new(PayRsp)
	if err = json.Unmarshal(bs, rsp); err != nil {
		return nil, fmt.Errorf("[%w], bytes: %s", err, string(bs))
	}

	if err = bizErrCheck(rsp.RspBase); err != nil {
		return nil, err
	}

	return rsp, nil
}

// h5支付宝支付下单
func (c *Client) H5AliPay(ctx context.Context, bm gopay.BodyMap) (rsp *PayRsp, err error) {
	err = bm.CheckEmptyError("merOrderId", "totalAmount", "notifyUrl")
	if err != nil {
		return nil, err
	}

	// 补充必填业务参数
	bm.
		Set("instMid", "H5DEFAULT") // 业务类型

	var bs []byte
	if bs, err = c.doGet(ctx, h5AliPayOrderPath, bm); err != nil {
		return nil, err
	}

	rsp = new(PayRsp)
	if err = json.Unmarshal(bs, rsp); err != nil {
		return nil, fmt.Errorf("[%w], bytes: %s", err, string(bs))
	}

	if err = bizErrCheck(rsp.RspBase); err != nil {
		return nil, err
	}

	return rsp, nil
}

// Refund 统一退款接口
func (c *Client) Refund(ctx context.Context, bm gopay.BodyMap) (rsp *RefundRsp, err error) {
	err = bm.CheckEmptyError("merOrderId", "refundAmount", "refundOrderId")
	if err != nil {
		return nil, err
	}

	var bs []byte
	if bs, err = c.doPost(ctx, refundOrderPath, bm); err != nil {
		return nil, err
	}

	rsp = new(RefundRsp)
	if err = json.Unmarshal(bs, rsp); err != nil {
		return nil, fmt.Errorf("[%w], bytes: %s", gopay.UnmarshalErr, string(bs))
	}

	if err = bizErrCheck(rsp.RspBase); err != nil {
		return nil, err
	}

	return rsp, nil
}
