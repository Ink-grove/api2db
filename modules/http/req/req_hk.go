package req

import (
	"encoding/json"
	"github.com/gogf/gf/os/glog"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/util/gconv"
	"github.com/google/uuid"
	"http2db/utils"
	"net/http"
	"strings"
)

type HkClient struct {
	Key        string `json:"r_auth"`
	Secret     string `json:"r_security"`
	BaseServer string `json:"base_server"`
}

func NewHkClient(b []byte) (c *HkClient) {
	if err := json.Unmarshal(b, &c); err != nil {
		glog.Error("[NewHkClient] unmarshal json to config struct error: ", err.Error())
		return nil
	}

	return
}

func (h *HkClient) NewRequest(method string, path string, body map[string]interface{}) *http.Request {
	buffer := utils.ToBuffer(body)
	method = strings.ToUpper(method)

	req, _ := http.NewRequest(method, h.BaseServer+path, buffer)

	headers := ""
	httpHeader := method + "\n"
	acceptStr := "*/*"
	req.Header.Set("Accept", acceptStr)
	httpHeader += acceptStr + "\n"

	contentTypeStr := "application/json"
	req.Header.Set("Content-Type", contentTypeStr)
	httpHeader += contentTypeStr + "\n"

	req.Header.Set("x-ca-key", h.Key)
	httpHeader += "x-ca-key:" + h.Key + "\n"
	uuidStr := uuid.New().String()
	req.Header.Set("x-ca-nonce", uuidStr)
	httpHeader += "x-ca-nonce:" + uuidStr + "\n"
	timestampStr := gconv.String(gtime.Now().UnixNano() / 1e6)
	req.Header.Set("x-ca-timestamp", timestampStr)
	httpHeader += "x-ca-timestamp:" + timestampStr + "\n" + path
	headers = "x-ca-key,x-ca-nonce,x-ca-timestamp"

	signature := utils.ComputeHmac256(httpHeader, h.Secret)
	req.Header.Set("x-ca-signature-headers", headers)
	req.Header.Set("x-ca-signature", signature)

	//glog.Printf("\nhttpHeader:\n%s\nsignature:\n%s\n", httpHeader, signature)

	return req
}
