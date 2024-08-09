package main

import (
	"http2db/cmd"
	"http2db/router"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	router.InitAPIRouter()
	cmd.StartDsu()
	wait()
}

func wait() {
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	os.Exit(0)
}
