package icbc

import (
	"crypto"
	"crypto/rsa"
	"encoding/base64"
	"fmt"

	"github.com/misu99/gopay"
)

// verifySign 验证响应签名
func (c *Client) verifySign(rsp *RspCommon) (err error) {
	sign := rsp.Sign
	signData := rsp.ResponseBizContent

	signBytes, _ := base64.StdEncoding.DecodeString(sign)
	hashs := crypto.SHA1
	h := hashs.New()
	h.Write(signData)
	if err = rsa.VerifyPKCS1v15(c.publicKey, hashs, h.Sum(nil), signBytes); err != nil {
		return fmt.Errorf("[%w]: %v", gopay.VerifySignatureErr, err)
	}
	return nil
}

// 回调通知验证签名
func (c *Client) VerifyNotifySign(path string, bm gopay.BodyMap) (err error) {
	sign := bm.Get("sign")
	bm.Remove("sign")
	signData := path + "?" + bm.EncodeAliPaySignParams()

	signBytes, _ := base64.StdEncoding.DecodeString(sign)
	hashs := crypto.SHA1
	h := hashs.New()
	h.Write([]byte(signData))
	if err = rsa.VerifyPKCS1v15(c.publicKey, hashs, h.Sum(nil), signBytes); err != nil {
		return fmt.Errorf("[%w]: %v", gopay.VerifySignatureErr, err)
	}
	return nil
}
