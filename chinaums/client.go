package icbc

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/misu99/gopay"
	"github.com/misu99/gopay/pkg/xhttp"
)

type Client struct {
	isProd bool // 是否正式环境

	appid      string
	appKey     string
	merchantNo string // 商户编号
	terminalNo string // 终端号
	msgSrc     string // 消息来源
	msgSrcId   string // 来源编号（作为账单号的前4位，且在商户系统中此账单号保证唯一。总长度需大于6位，小于28位）
	secretKey  string // 通讯密钥
}

// NewClient 初始化银联支付客户端
// merchantNo：商户编号
// appId：appid
// appKey：appKey
// terminalNo：终端号
// msgSrc：消息来源
// msgSrcId：来源编号（作为账单号的前4位，且在商户系统中此账单号保证唯一。总长度需大于6位，小于28位）
// secretKey：通讯密钥
// isProd：是否正式环境
func NewClient(merchantNo, appid, appKey, terminalNo, msgSrc, msgSrcId, secretKey string, isProd bool) *Client {
	return &Client{
		isProd:     isProd,
		appid:      appid,
		appKey:     appKey,
		merchantNo: merchantNo,
		terminalNo: terminalNo,
		msgSrc:     msgSrc,
		msgSrcId:   msgSrcId,
		secretKey:  secretKey,
	}
}

// getSign 获取签名串
func (c *Client) getSign(body []byte, timestampStr, nonceStr string) (res string, err error) {
	// 1.取报文体，即正文全部内容获得字节数组 ，进行SHA256算法取十六进制小写字符串
	h := sha256.New()
	B := hex.EncodeToString(h.Sum(body))

	// 2.取 、 、 、 进行字符串拼接，以 进行编码获得待签名串 ，取 作为签名密钥 。
	C := c.appid + timestampStr + nonceStr + B

	// 3.以C和D进行HMAC-SHA256算法获得签名字节数组E
	mac := hmac.New(sha256.New, []byte(c.appKey))
	_, _ = mac.Write([]byte(C))
	E := mac.Sum(nil)

	return base64.StdEncoding.EncodeToString(E), nil
}

// OPEN-BODY-SIG签名串
func (c *Client) getSignBodySig(body []byte) (res string, err error) {
	timestampStr := strconv.FormatInt(time.Now().Unix(), 10)
	nonceStr := strconv.FormatInt(time.Now().UnixNano(), 10)

	sign, err := c.getSign(body, timestampStr, nonceStr)
	if err != nil {
		return
	}

	res = fmt.Sprintf(`OPEN-BODY-SIG AppId="%s",Timestamp="%s",Nonce="%s",Signature="%s"`,
		c.appid, timestampStr, nonceStr, sign)

	return
}

// OPEN-FORM-PARAM签名串
func (c *Client) getSignFormSig(body []byte) (res string, err error) {
	timestampStr := strconv.FormatInt(time.Now().Unix(), 10)
	nonceStr := strconv.FormatInt(time.Now().UnixNano(), 10)

	sign, err := c.getSign(body, timestampStr, nonceStr)
	if err != nil {
		return
	}

	params := make(gopay.BodyMap)
	params.
		Set("authorization", "OPEN-FORM-PARAM").
		Set("appId", c.appid).
		Set("timestamp", timestampStr).
		Set("nonce", nonceStr).
		Set("content", string(body)).
		Set("signature", sign)

	res = params.EncodeURLParams()

	return
}

// pubParamsHandle 公共参数处理
func (c *Client) pubParamsHandle(bm gopay.BodyMap) (param []byte, err error) {
	msgID := strconv.FormatInt(time.Now().UnixNano(), 10)

	// 补充业务请求参数
	bm.
		Set("msgId", msgID).                                                            // 消息ID，原样返回
		Set("mid", c.merchantNo).                                                       // 商户号
		Set("tid", c.terminalNo).                                                       // 终端号
		Set("requestTimestamp", time.Now().In(DefaultLocation).Format(timestampLayout)) // 报文请求时间 格式yyyy-MM-dd HH:mm:ss

	// 序列化参数json
	bizContent, err := json.Marshal(bm)
	if err != nil {
		return
	}

	return bizContent, nil
}

// doPost 发起请求
func (c *Client) doPost(ctx context.Context, path string, bm gopay.BodyMap) (bs []byte, err error) {
	param, err := c.pubParamsHandle(bm)
	if err != nil {
		return nil, err
	}

	// 计算参数签名
	authorization, err := c.getSignFormSig(param)
	if err != nil {
		return
	}

	urlBase := baseUrl
	if !c.isProd {
		urlBase = sandboxBaseUrl
	}

	httpClient := xhttp.NewClient()
	httpClient.Header.Set("Authorization", authorization)

	res, bs, err := httpClient.Type(xhttp.TypeForm).Post(urlBase + path).SendString(string(param)).EndBytes(ctx)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP Request Error, StatusCode = %d", res.StatusCode)
	}
	return bs, nil
}

// doGet 发起请求
func (c *Client) doGet(ctx context.Context, path string, bm gopay.BodyMap) (bs []byte, err error) {
	param, err := c.pubParamsHandle(bm)
	if err != nil {
		return nil, err
	}

	// 计算参数签名
	params, err := c.getSignFormSig(param)
	if err != nil {
		return
	}

	urlBase := baseUrl
	if !c.isProd {
		urlBase = sandboxBaseUrl
	}
	urlBase += path + "?" + params

	httpClient := xhttp.NewClient()
	res, bs, err := httpClient.Type(xhttp.TypeUrlencoded).Get(urlBase).EndBytes(ctx)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP Request Error, StatusCode = %d", res.StatusCode)
	}
	return bs, nil
}
