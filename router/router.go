package router

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"http2db/api"
	"http2db/middleware"
)

func InitAPIRouter() {
	s := g.Server()
	//应用接口
	s.Use(middleware.Cors)
	s.Group("http2db/api", func(g *ghttp.RouterGroup) {
		action := api.NewAction()
		g.ALL("/getTaskStatus", action.GetTaskStatus)
		g.ALL("/getAllCronJob", action.GetAllCronJob)
		g.ALL("/stop", action.Stop)
		g.ALL("/stopAll", action.StopAll)
		g.ALL("/pause", action.Pause)
		g.ALL("/keepOn", action.KeepOn)

	})
}

func APPStart() []byte {
	return nil
}
