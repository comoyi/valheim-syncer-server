package app

import (
	"github.com/comoyi/valheim-syncer-server/config"
	"github.com/comoyi/valheim-syncer-server/server"
)

func Start() {
	config.LoadConfig()
	server.Start()
}
