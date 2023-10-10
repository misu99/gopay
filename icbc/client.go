package icbc

import (
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"hash"
	"strconv"
	"time"

	"github.com/misu99/gopay"
	"github.com/misu99/gopay/pkg/util"
	"github.com/misu99/gopay/pkg/xhttp"
	"github.com/misu99/gopay/pkg/xpem"
	"github.com/misu99/gopay/pkg/xrsa"
)

type Client struct {
	isProd bool // 是否正式环境

	merId           string // 商户编号
	appId           string // 商户的工行APPID
	serialNo        string // 收单产品协议编号
	clearingAccount string // 商户清算账号

	privateKey *rsa.PrivateKey // 商户的私钥
	publicKey  *rsa.PublicKey  // 网关的公钥
}

// NewClient 初始化工行客户端
// merId：商户编号
// appId：商户的工行APPID
// serialNo：收单产品协议编号
// clearingAccount：商户清算账号
// privateKey：商户的私钥
// publicKey：工行网关的公钥
// publicKey：是否正式环境
func NewClient(merId, appId, serialNo, clearingAccount, privateKey, publicKey string, isProd bool) (*Client, error) {
	prk, err := xpem.DecodePrivateKey([]byte(xrsa.FormatAlipayPrivateKey(privateKey)))
	if err != nil {
		return nil, err
	}
	puk, err := xpem.DecodePublicKey([]byte(xrsa.FormatAlipayPublicKey(publicKey)))
	if err != nil {
		return nil, err
	}
	return &Client{
		isProd:          isProd,
		merId:           merId,
		appId:           appId,
		serialNo:        serialNo,
		clearingAccount: clearingAccount,
		privateKey:      prk,
		publicKey:       puk,
	}, nil
}

// getRsaSign 获取签名字符串(&拼接参数)
func (c *Client) getRsaSign(path string, bm gopay.BodyMap, signType string, privateKey *rsa.PrivateKey) (sign string, err error) {
	var (
		h              hash.Hash
		hashs          crypto.Hash
		encryptedBytes []byte
	)

	switch signType {
	case RSA:
		h = sha1.New()
		hashs = crypto.SHA1
	case RSA2:
		h = sha256.New()
		hashs = crypto.SHA256
	default:
		h = sha1.New()
		hashs = crypto.SHA1
	}

	signParams := bm.EncodeAliPaySignParams()
	if path != "" {
		signParams = path + "?" + signParams
	}

	if _, err = h.Write([]byte(signParams)); err != nil {
		return
	}

	if encryptedBytes, err = rsa.SignPKCS1v15(rand.Reader, privateKey, hashs, h.Sum(nil)); err != nil {
		return util.NULL, fmt.Errorf("[%w]: %+v", gopay.SignatureErr, err)
	}

	sign = base64.StdEncoding.EncodeToString(encryptedBytes)
	return
}

// getRsaSign 获取签名字符串(json拼接参数)
func (c *Client) getRsaSign2(signParams, signType string, privateKey *rsa.PrivateKey) (sign string, err error) {
	var (
		h              hash.Hash
		hashs          crypto.Hash
		encryptedBytes []byte
	)

	switch signType {
	case RSA:
		h = sha1.New()
		hashs = crypto.SHA1
	case RSA2:
		h = sha256.New()
		hashs = crypto.SHA256
	default:
		h = sha1.New()
		hashs = crypto.SHA1
	}

	if _, err = h.Write([]byte(signParams)); err != nil {
		return
	}

	if encryptedBytes, err = rsa.SignPKCS1v15(rand.Reader, privateKey, hashs, h.Sum(nil)); err != nil {
		return util.NULL, fmt.Errorf("[%w]: %+v", gopay.SignatureErr, err)
	}

	sign = base64.StdEncoding.EncodeToString(encryptedBytes)
	return
}

// pubParamsHandle 公共参数处理
func (c *Client) pubParamsHandle(path string, bm gopay.BodyMap) (param string, err error) {
	msgID := strconv.FormatInt(time.Now().UnixNano(), 10)

	// 通用请求参数
	params := make(gopay.BodyMap)
	params.Set("app_id", c.appId).
		Set("msg_id", msgID).
		Set("format", formatJSON).
		Set("charset", charsetUTF8).
		Set("sign_type", RSA2).
		Set("timestamp", time.Now().In(DefaultLocation).Format(timestampLayout))

	// 业务请求参数
	bm.Set("mer_id", c.merId)
	bm.Set("mer_prtcl_no", c.serialNo)
	bm.Set("icbc_appid", c.appId)
	bm.Set("mer_acct", c.clearingAccount) // 商户清算账号
	bizContent, err := json.Marshal(bm)
	if err != nil {
		return
	}
	params.Set("biz_content", string(bizContent))

	// 计算参数签名
	sign, err := c.getRsaSign(path, params, params.GetString("sign_type"), c.privateKey)
	if err != nil {
		return "", fmt.Errorf("GetRsaSign Error: %w", err)
	}
	params.Set("sign", sign)

	param = params.EncodeURLParams()
	return param, nil
}

// doPost 发起请求
func (c *Client) doPost(ctx context.Context, path string, bm gopay.BodyMap) (bs []byte, err error) {
	param, err := c.pubParamsHandle(path, bm)
	if err != nil {
		return nil, err
	}
	httpClient := xhttp.NewClient()
	url := baseUrl
	if !c.isProd {
		url = sandboxBaseUrl
	}

	res, bs, err := httpClient.Type(xhttp.TypeForm).Post(url + path).SendString(param).EndBytes(ctx)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP Request Error, StatusCode = %d", res.StatusCode)
	}
	return bs, nil
}
