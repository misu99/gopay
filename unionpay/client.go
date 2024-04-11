package unionpay

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/url"
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
func NewClient(merchantNo, appid, appKey, terminalNo, secretKey string, isProd bool) *Client {
	return &Client{
		isProd:     isProd,
		appid:      appid,
		appKey:     appKey,
		merchantNo: merchantNo,
		terminalNo: terminalNo,
		secretKey:  secretKey,
	}
}

// 获取签名串
func (c *Client) getSign(appid, appKey, timestamp, nonce string, body []byte) string {
	// 第一步SHA256算法转十六进制
	h := sha256.New()
	h.Write(body)
	sum := h.Sum(nil)
	B := hex.EncodeToString(sum)

	// 第二步拼接参数
	C := appid + timestamp + nonce + B

	// 第三步HMAC-SHA256算法获得签名字节数组
	mac := hmac.New(sha256.New, []byte(appKey))
	_, _ = mac.Write([]byte(C))
	E := mac.Sum(nil)

	return base64.StdEncoding.EncodeToString(E)
}

// OPEN-BODY-SIG签名串
func (c *Client) getSignBodySig(appid, appKey, timestamp, nonce string, body []byte) string {
	sign := c.getSign(appid, appKey, timestamp, nonce, body)

	return fmt.Sprintf(`OPEN-BODY-SIG AppId="%s", Timestamp="%s", Nonce="%s", Signature="%s"`,
		appid, timestamp, nonce, sign)
}

// OPEN-FORM-PARAM签名串
func (c *Client) getSignFormSig(appid, appKey, timestamp, nonce string, body []byte) string {
	sign := c.getSign(appid, appKey, timestamp, nonce, body)

	//params := make(gopay.BodyMap)
	//params.
	//	Set("authorization", "OPEN-FORM-PARAM").
	//	Set("appId", c.appid).
	//	Set("timestamp", timestamp).
	//	Set("nonce", nonce).
	//	Set("content", string(body)).
	//	Set("signature", sign)

	param := "?" + "authorization=" + sign + "&appId=" + appid + "&timestamp=" + "20190812160100" + "&nonce=" + "nonce" + "&content=" + url.QueryEscape(string(body)) + "&signature=" + url.QueryEscape(sign)

	return param
}

// pubParamsHandle 公共参数处理
func (c *Client) pubParamsHandle(bm gopay.BodyMap) (res []byte, err error) {
	msgID := strconv.FormatInt(time.Now().UnixNano(), 10)

	// 补充业务请求参数
	bm.
		Set("msgId", msgID).                                                            // 消息ID，原样返回
		Set("mid", c.merchantNo).                                                       // 商户号
		Set("tid", c.terminalNo).                                                       // 终端号
		Set("requestTimestamp", time.Now().In(DefaultLocation).Format(timestampLayout)) // 报文请求时间 格式yyyy-MM-dd HH:mm:ss

	biz, err := json.Marshal(bm)
	if err != nil {
		return nil, err
	}

	return biz, nil
}

// doPost 发起请求
func (c *Client) doPost(ctx context.Context, path string, bm gopay.BodyMap) (bs []byte, err error) {
	param, err := c.pubParamsHandle(bm)
	if err != nil {
		return nil, err
	}

	// 计算参数签名
	timestamp := time.Now().Format("20060102150405")
	nonce := strconv.FormatInt(time.Now().UnixNano(), 10)
	authorization := c.getSignBodySig(c.appid, c.appKey, timestamp, nonce, param)

	urlBase := baseUrl
	if !c.isProd {
		urlBase = sandboxBaseUrl
	}
	urlBase += path

	httpClient := xhttp.NewClient()
	httpClient.Header.Add("Authorization", authorization)

	res, bs, err := httpClient.Type(xhttp.TypeForm).Post(urlBase).SendString(string(param)).EndBytes(ctx)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("HTTPCode=%d Content=%s", res.StatusCode, string(bs))
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
	timestamp := time.Now().Format("20060102150405")
	nonce := strconv.FormatInt(time.Now().UnixNano(), 10)
	queryParam := c.getSignFormSig(c.appid, c.appKey, timestamp, nonce, param)

	urlBase := baseUrl
	if !c.isProd {
		urlBase = sandboxBaseUrl
	}
	urlBase += path + "?" + queryParam

	httpClient := xhttp.NewClient()
	res, bs, err := httpClient.Get(urlBase).EndBytes(ctx)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("HTTPCode=%d Content=%s", res.StatusCode, string(bs))
	}
	return bs, nil
}
