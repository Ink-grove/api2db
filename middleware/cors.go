package middleware

import (
	"github.com/gogf/gf/net/ghttp"
	"net/http"
)

func Cors(r *ghttp.Request) {
	r.Response.CORSDefault()
	r.Middleware.Next()
}

func Cors2(r *ghttp.Request) {
	method := r.Request.Method
	origin := r.Request.Header.Get("Origin")
	r.Header.Set("Access-Control-Allow-Origin", origin)
	r.Header.Set("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token,X-Token,X-User-Id")
	r.Header.Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS,DELETE,PUT")
	r.Header.Set("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type, New-Token, New-Expires-At")
	r.Header.Set("Access-Control-Allow-Credentials", "true")

	//r.Response.CORSDefault()
	// 放行所有OPTIONS方法
	if method == "OPTIONS" {
		r.Response.Writer.Status = http.StatusNoContent
		r.Exit()
	}

	r.Middleware.Next()
}
