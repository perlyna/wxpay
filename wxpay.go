package wxpay

import (
	"errors"
	"fmt"
	"net/http"
)

const (
	domainAPI  = "api.mch.weixin.qq.com"
	domainAPI2 = "api2.mch.weixin.qq.com"
)

// WechatPay 微信支付SDK
type WechatPay struct {
	appID     string       // 应用 appid
	mchID     string       // 商户号
	apiSecret string       // 商户号 API Secret
	domain    string       // 微信支付域名
	notifyURL string       // 支付通知地址
	signType  SignType     // 签名类型, 默认为MD5签名
	client    *http.Client // http client
	tlsClient *http.Client // 带有证书的 http client
}

// New 创建微信支付模块
// 默认的签名类型为 MD5
func New(appid, mchid, apiSecret string, opts ...Option) *WechatPay {
	config := &WechatPay{
		appID:     appid,
		mchID:     mchid,
		apiSecret: apiSecret,
		domain:    domainAPI,
		signType:  MD5,
		client:    DefaultClient,
	}
	for _, opt := range opts {
		opt(config)
	}
	return config
}

// 向 Map 中添加 appid、mch_id、nonce_str、sign_type、sign
// 该函数适用于商户适用于统一下单等接口，不适用于红包、代金券接口
func (c *WechatPay) fillRequestData(params map[string]string) {
	if _, ok := params["appid"]; !ok { // 优先使用参数的appid
		params["appid"] = c.appID
	}
	if _, ok := params["mch_id"]; !ok { // 优先使用参数的mch_id
		params["mch_id"] = c.mchID
	}
	if _, ok := params["nonce_str"]; !ok { // 优先使用参数的随机数
		params["nonce_str"] = GenerateNonceStr(32)
	}
	if _, ok := params["sign_type"]; !ok { // 优先使用参数的签名类型
		params["sign_type"] = c.signType
	}
	params[FieldSign] = Sign(params, c.apiSecret, params["sign_type"])
}

// processResponseXML 处理 HTTPS API返回数据，转换成Map对象。return_code为SUCCESS时，验证签名。
func (c *WechatPay) processResponseXML(body []byte, signType SignType) (map[string]string, error) {
	ret, err := xmlToMap(body) // 转换成Map对象
	if err != nil {
		return ret, err
	}
	retCode, ok := ret["return_code"]
	if !ok { // 返回书架没有 return_code 参数
		return ret, fmt.Errorf("no `return_code` in XML: \n%s", body)
	}
	if retCode == "FAIL" { // 请求通信标识, 返回通信失败的错误信息
		return ret, fmt.Errorf(ret["return_msg"]) // eg: 签名失败, 参数格式校验错误
	}
	if retCode != "SUCCESS" { // 无效的通信标识
		return ret, fmt.Errorf("return_code value %s is invalid in XML: %s", retCode, body)
	}
	if !IsSignatureValid(ret, c.apiSecret, signType) { //  验签失败
		return ret, fmt.Errorf("invalid sign value in XML: %s", body)
	}
	return ret, nil
}

// MicroPay 扫码支付
// 文档链接: https://pay.weixin.qq.com/wiki/doc/api/micropay.php?chapter=9_10&index=1
func (c *WechatPay) MicroPay(params map[string]string) (ret map[string]string, err error) {
	c.fillRequestData(params)
	body, err := Request(c.client, c.domain, "/pay/micropay", mapToXML(params), params["nonce_str"])
	if err == nil {
		return nil, err
	}
	return c.processResponseXML(body, params["sign_type"])
}

// UnifiedOrder 统一下单
// 文档链接: https://pay.weixin.qq.com/wiki/doc/api/jsapi.php?chapter=9_1
// 文档链接: https://pay.weixin.qq.com/wiki/doc/api/native.php?chapter=9_1
// 文档链接: https://pay.weixin.qq.com/wiki/doc/api/app/app.php?chapter=9_1
// 文档链接: https://pay.weixin.qq.com/wiki/doc/api/H5.php?chapter=9_20&index=1
// 文档链接: https://pay.weixin.qq.com/wiki/doc/api/wxa/wxa_api.php?chapter=9_1
func (c *WechatPay) UnifiedOrder(params map[string]string) (ret map[string]string, err error) {
	if _, ok := params["notify_url"]; !ok {
		params["notify_url"] = c.notifyURL
	}
	c.fillRequestData(params)
	body, err := Request(c.client, c.domain, "/pay/unifiedorder", mapToXML(params), params["nonce_str"])
	if err != nil {
		return nil, err
	}
	return c.processResponseXML(body, params["sign_type"])
}

// OrderQuery 查询订单
// 文档链接: https://pay.weixin.qq.com/wiki/doc/api/jsapi.php?chapter=9_2
func (c *WechatPay) OrderQuery(params map[string]string) (ret map[string]string, err error) {
	c.fillRequestData(params)
	body, err := Request(c.client, c.domain, "/pay/orderquery", mapToXML(params), params["nonce_str"])
	if err != nil {
		return nil, err
	}
	return c.processResponseXML(body, params["sign_type"])
}

// Reverse 撤销订单
// 文档链接: https://pay.weixin.qq.com/wiki/doc/api/micropay.php?chapter=9_11&index=3
func (c *WechatPay) Reverse(params map[string]string) (ret map[string]string, err error) {
	if c.tlsClient == nil {
		return nil, errors.New("the request need cert")
	}
	c.fillRequestData(params)
	body, err := Request(c.tlsClient, c.domain, "/secapi/pay/reverse", mapToXML(params), params["nonce_str"])
	if err != nil {
		return nil, err
	}
	return c.processResponseXML(body, params["sign_type"])
}

// CloseOrder 关闭订单
// 文档链接: https://pay.weixin.qq.com/wiki/doc/api/jsapi.php?chapter=9_3
func (c *WechatPay) CloseOrder(params map[string]string) (ret map[string]string, err error) {
	c.fillRequestData(params)
	body, err := Request(c.client, c.domain, "/pay/closeorder", mapToXML(params), params["nonce_str"])
	if err != nil {
		return nil, err
	}
	return c.processResponseXML(body, params["sign_type"])
}

// Refund 申请退款
// 文档链接: https://pay.weixin.qq.com/wiki/doc/api/jsapi.php?chapter=9_4
func (c *WechatPay) Refund(params map[string]string) (ret map[string]string, err error) {
	if c.tlsClient == nil {
		return nil, errors.New("the request need cert")
	}
	c.fillRequestData(params)
	body, err := Request(c.tlsClient, c.domain, "/secapi/pay/refund", mapToXML(params), params["nonce_str"])
	if err != nil {
		return nil, err
	}
	return c.processResponseXML(body, params["sign_type"])
}

// RefundQuery 退款查询
// 文档链接: https://pay.weixin.qq.com/wiki/doc/api/jsapi.php?chapter=9_5
func (c *WechatPay) RefundQuery(params map[string]string) (ret map[string]string, err error) {
	c.fillRequestData(params)
	body, err := Request(c.client, c.domain, "/pay/refundquery", mapToXML(params), params["nonce_str"])
	if err != nil {
		return nil, err
	}
	return c.processResponseXML(body, params["sign_type"])
}

// DownloadBill 下载交易账单
// 文档链接: https://pay.weixin.qq.com/wiki/doc/api/jsapi.php?chapter=9_6
func (c *WechatPay) DownloadBill(params map[string]string) ([]byte, error) {
	c.fillRequestData(params)
	return Request(c.client, c.domain, "/pay/downloadbill", mapToXML(params), params["nonce_str"])
}

// DownloadFundflow 下载资金账单
// 文档链接: https://pay.weixin.qq.com/wiki/doc/api/jsapi.php?chapter=9_18&index=7
func (c *WechatPay) DownloadFundflow(params map[string]string) ([]byte, error) {
	c.fillRequestData(params)
	return Request(c.client, c.domain, "/pay/downloadfundflow", mapToXML(params), params["nonce_str"])
}

// ShortURL 转换短链接
// 文档链接: https://pay.weixin.qq.com/wiki/doc/api/native.php?chapter=9_9&index=10
func (c *WechatPay) ShortURL(params map[string]string) (map[string]string, error) {
	c.fillRequestData(params)
	body, err := Request(c.client, c.domain, "/tools/shorturl", mapToXML(params), params["nonce_str"])
	if err != nil {
		return nil, err
	}
	return c.processResponseXML(body, params["sign_type"])
}

// AuthCodeToOpenid 授权码查询OPENID接口
// 文档链接: https://pay.weixin.qq.com/wiki/doc/api/micropay.php?chapter=9_13&index=9
func (c *WechatPay) AuthCodeToOpenid(params map[string]string) (map[string]string, error) {
	c.fillRequestData(params)
	body, err := Request(c.client, c.domain, "/tools/authcodetoopenid", mapToXML(params), params["nonce_str"])
	if err != nil {
		return nil, err
	}
	return c.processResponseXML(body, params["sign_type"])
}

// SignType 获取当前商户号的签名类型
func (c *WechatPay) SignType() string {
	return c.signType
}

// APISecret 获取当前商户号的API密钥
func (c *WechatPay) APISecret() string {
	return c.apiSecret
}
