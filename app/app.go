package app

import (
	"github.com/comoyi/valheim-syncer-server/config"
	"github.com/comoyi/valheim-syncer-server/server"
)

func Start() {
	config.LoadConfig()

	go func() {
		server.Start()
	}()

	if config.Conf.Gui != "OFF" {
		server.StartGUI()
	} else {
		c := make(chan struct{})
		<-c
	}
}
