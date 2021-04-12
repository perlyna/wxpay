package wxpay

import (
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"hash"
	"math/rand"
	"sort"
	"strings"
	"time"
)

// FieldSign 请求参数中的签名字段
const FieldSign = "sign"

//SignType 签名类型
type SignType = string

const (
	// MD5 MD5签名类型, 默认
	MD5 SignType = "MD5"
	// HMACSHA256 SHA256 签名类型
	HMACSHA256 SignType = "HMAC-SHA256"
)

// Sign 获取参数签名值
func Sign(params map[string]string, apiKey string, signType SignType) string {
	var keys = make([]string, 0, len(params))
	for k := range params {
		if k != FieldSign { // 排除 sign 字段
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)
	var buf bytes.Buffer
	for _, k := range keys {
		if len(params[k]) > 0 {
			buf.WriteString(k)
			buf.WriteString("=")
			buf.WriteString(params[k])
			buf.WriteString("&")
		}
	}
	buf.WriteString("key=")
	buf.WriteString(apiKey)
	var h hash.Hash
	if signType == HMACSHA256 {
		h = hmac.New(sha256.New, []byte(apiKey))
	} else {
		h = md5.New() // 签名默认是 MD5
	}
	h.Write(buf.Bytes())
	return strings.ToUpper(hex.EncodeToString(h.Sum(nil)))
}

// IsSignatureValid 判断签名是否有效
func IsSignatureValid(params map[string]string, apiKey string, signType SignType) bool {
	if sign, ok := params[FieldSign]; ok {
		return sign == Sign(params, apiKey, signType)
	}
	return false
}

// nonceStr 随机字符串
const nonceStr = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// nonceStrLength 随机字符串长度
const nonceStrLength = len(nonceStr)

// GenerateNonceStr 获取随机字符串
func GenerateNonceStr(length int) string {
	result := make([]byte, length)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		result[i] = nonceStr[r.Intn(nonceStrLength)]
	}
	return string(result)
}
