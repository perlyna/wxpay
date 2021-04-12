package wxpay

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

// DefaultClient is the default Client
var DefaultClient = &http.Client{
	Transport: &http.Transport{
		Dial: func(network, addr string) (net.Conn, error) {
			conn, err := net.DialTimeout(network, addr, 6*time.Second)
			if err != nil {
				return nil, err
			}
			conn.SetDeadline(time.Now().Add(8 * time.Second))
			return conn, nil
		},
		MaxIdleConnsPerHost: http.DefaultMaxIdleConnsPerHost, // host最大复用个数 默认2
	},
}

// NewCertClient 创建使用API证书的请求客服端
// 由于绝大部分操作系统已内置了微信支付服务器证书的根CA证书,  2018年3月6日后, 不再提供CA证书文件（rootca.pem）下载,
// rootca 参数为 nil 即使用系统的根CA证书
func NewCertClient(cert, key, rootca []byte) (client *http.Client, err error) {
	var tlsCert tls.Certificate
	if tlsCert, err = tls.X509KeyPair(cert, key); err != nil {
		return
	}
	var RootCAs *x509.CertPool
	if len(rootca) > 0 {
		RootCAs = x509.NewCertPool()
		if ok := RootCAs.AppendCertsFromPEM(rootca); !ok {
			err = errors.New("failed to parse root certificate")
			return
		}
	}
	client = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				Certificates: []tls.Certificate{tlsCert},
				RootCAs:      RootCAs,
			},
			Dial: func(network, addr string) (net.Conn, error) {
				conn, err := net.DialTimeout(network, addr, 6*time.Second)
				if err != nil {
					return nil, err
				}
				conn.SetDeadline(time.Now().Add(8 * time.Second))
				return conn, nil
			},
			MaxIdleConnsPerHost: http.DefaultMaxIdleConnsPerHost, // host最大复用个数 默认2
		},
	}
	return
}

// Request request
func Request(client *http.Client, domain, urlSuffix string, body []byte, uuid string) ([]byte, error) {
	request, err := http.NewRequest(http.MethodPost, "https://"+domain+urlSuffix, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	ret, err := RequestOnce(client, request)
	return ret, err
}

// RequestOnce request once
func RequestOnce(client *http.Client, request *http.Request) ([]byte, error) {
	request.Header.Set("Content-Type", "text/xml")
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return body, err
	}
	if response.StatusCode >= http.StatusBadRequest { // 返回状态码错误
		return body, fmt.Errorf("response code: %d", response.StatusCode)
	}
	return body, nil
}
