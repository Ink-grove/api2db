package utils

import (
	"github.com/gogf/gf/net/ghttp"
	"net/http"
)

type BaseRouter struct {
}

type Resp struct {
	Code int         `json:"code"`
	Msg  interface{} `json:"msg"`
}

type RespData struct {
	Code    int         `json:"code"`
	Encrypt int         `json:"encrypt"`
	Msg     string      `json:"msg"`
	Data    interface{} `json:"data"`
}

func InfoSuccess(r *ghttp.Request) {
	r.Response.WriteJson(Resp{http.StatusOK, "success"})
	r.Exit()
}

func Info(r *ghttp.Request, msg interface{}) {
	r.Response.WriteJson(Resp{http.StatusOK, msg})
	r.Exit()
}

func Data(r *ghttp.Request, data interface{}) {
	r.Response.WriteJson(RespData{
		Code: http.StatusOK,
		Data: data,
	})
	r.Exit()
}
