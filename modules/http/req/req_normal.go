package req

import (
	"encoding/json"
	"github.com/gogf/gf/os/glog"
	"http2db/utils"
	"net/http"
)

type NormalClient struct {
	BaseServer string `json:"base_server"`
}

func NewNormalClient(b []byte) (n *NormalClient) {
	if err := json.Unmarshal(b, &n); err != nil {
		glog.Error("[NewNormalClient] unmarshal json to config struct error: ", err.Error())
		return nil
	}

	if empty := n.isEmpty(); empty {
		glog.Error("[NewNormalClient] NormalClient isEmpty , please check param")
		return nil
	}

	return
}

func (n *NormalClient) isEmpty() bool {
	return n.BaseServer == ""
}

func (n *NormalClient) NewRequest(method string, path string, body map[string]interface{}) *http.Request {
	buffer := utils.ToBuffer(body)
	req, _ := http.NewRequest(method, n.BaseServer+path, buffer)
	req.Header.Set("Content-Type", "application/json;charset=utf-8")
	req.Header.Set("Token", "_G7CU2FKpdD-ht7PCVp_bdXnCjFzl0_sVJvtGhupp13M1u2hiUIeRS5wBrn_f9M9")
	return req
}
