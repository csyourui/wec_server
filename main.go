package main

import (
	"github.com/csyourui/wechat_server/command"
	"github.com/csyourui/wechat_server/pkg/comm"
	"github.com/csyourui/wechat_server/pkg/log"
)

func main() {
	log.Logger.Debug("Start")
	comm.Execute(command.Root)
}
