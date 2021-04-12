package wxpay

import (
	"io/ioutil"
	"net/http"
)

// ParesNotify 解析通知
func ParesNotify(r *http.Request) (map[string]string, error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	return xmlToMap(body)
}

// ResponseSuccess 返回通知成功
func ResponseSuccess() []byte {
	return mapToXML(map[string]string{"return_code": "SUCCESS", "return_msg": "OK"})
}

// ResponseFail 返回通知失败
func ResponseFail(errMsg string) []byte {
	return mapToXML(map[string]string{"return_code": "FAIL", "return_msg": errMsg})
}
