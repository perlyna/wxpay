package wxpay

import "net/http"

// Option 设置微信支付的其他属性
type Option func(*WechatPay)

// WithMD5 设置签名类型为 MD5 ,默认为 MD5
func WithMD5() Option {
	return func(c *WechatPay) {
		c.signType = MD5
	}
}

// WithHMACSHA256 设置签名类型为 HMACSHA256
func WithHMACSHA256() Option {
	return func(c *WechatPay) {
		c.signType = HMACSHA256
	}
}

// WithTLSClient API证书的请求客户端
// 如果不设置, 需要证书的请求无法正常运行
func WithTLSClient(tlsClient *http.Client) Option {
	return func(c *WechatPay) {
		c.tlsClient = tlsClient
	}
}

// WithNotifyURL 支付通知地址
func WithNotifyURL(notifyURL string) Option {
	return func(c *WechatPay) {
		c.notifyURL = notifyURL
	}
}
