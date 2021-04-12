package wxpay

import (
	"bytes"
	"encoding/xml"
	"io"
	"strings"
)

// xmlToMap XML格式字符串转换为map
func xmlToMap(body []byte) (map[string]string, error) {
	params := make(map[string]string)
	decoder := xml.NewDecoder(bytes.NewReader(body))
	var key string
	var value string
	var err error
	var t xml.Token
	for t, err = decoder.Token(); err == nil; t, err = decoder.Token() {
		switch token := t.(type) {
		case xml.StartElement:
			key = token.Name.Local
		case xml.CharData:
			value = strings.TrimSpace(string([]byte(token)))
		case xml.EndElement:
			if token.Name.Local == key && value != "" {
				params[key] = value
			}
		}
	}
	if err == io.EOF {
		return params, nil
	}
	return params, err
}

// mapToXML 将Map转换为XML格式的字符串
func mapToXML(params map[string]string) []byte {
	var buf bytes.Buffer
	buf.WriteString(`<xml>`)
	for k, v := range params {
		// buf.WriteString("<" + k + "><![CDATA[" + v + "]]></" + k + ">")
		buf.WriteString("<")
		buf.WriteString(k)
		buf.WriteString("><![CDATA[")
		buf.WriteString(v)
		buf.WriteString("]]></")
		buf.WriteString(k)
		buf.WriteString(">")
	}
	buf.WriteString(`</xml>`)
	return buf.Bytes()
}
