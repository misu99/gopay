package unionpay

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"github.com/misu99/gopay"
)

func (c *Client) VerifyNotifySign(bm gopay.BodyMap) (err error) {
	srcSign := bm.Get("sign")

	// 构建签名串
	bm.Remove("sign")
	signData := bm.EncodeAliPaySignParams()
	signData += c.secretKey

	// 计算签名
	h := md5.New()
	h.Write([]byte(signData))
	sign := hex.EncodeToString(h.Sum(nil))

	if srcSign != sign {
		return errors.New("sign not match")
	}

	return nil
}
