package cmd

import (
	"http2db/ctl"
)

var (
	G_dsu = ctl.Controller{}
)

func StartDsu() {
	go func() {
		G_dsu.Init()
		G_dsu.Run()
	}()
}
