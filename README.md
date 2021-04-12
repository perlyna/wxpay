# wxpay
微信支付V2 Go SDK
------

对[微信支付开发者文档](https://pay.weixin.qq.com/wiki/doc/api/index.html)中给出的API进行了封装。

提供了对应的方法：

|方法名 | 说明 |
|--------|--------|
|MicroPay| 刷卡支付 |
|UnifiedOrder | 统一下单|
|OrderQuery | 查询订单 |
|Reverse | 撤销订单 |
|CloseOrder|关闭订单|
|Refund|申请退款|
|RefundQuery|查询退款|
|DownloadBill|下载对账单|
|ShortUrl|转换短链接|
|AuthCodeToOpenid|授权码查询openid|

## 示例
配置 WechatPay
```golang
import github.com/perlyna/wxpay

tlsClient, _ := wxpay.NewCertClient(pk, pem, nil)

wechatPay := wxpay.New(appID, mchID, apiSecret,
 wxpay.WithNotifyURL(notifyURL),
 wxpay.WithTLSClient(tlsClient))
```

统一下单：
```golang
order := make(map[string]string)
order["body"] = "腾讯充值中心-QQ会员充值"
order["out_trade_no"] = "2016090910595900000012"
order["device_info"] = ""
order["fee_type"] = "CNY"
order["total_fee"] = "1"
order["spbill_create_ip"] = "123.12.12.123"
order["notify_url"] = `http://www.example.com/wxpay/notify`
order["trade_type"] = `NATIVE`
order["product_id"] = `12`
ret, err := wechatPay.UnifiedOrder(order)
```
订单查询：
```golang
order := make(map[string]string)
order["out_trade_no"] = "2016090910595900000012"
ret, err := wechatPay.OrderQuery(order)
```
