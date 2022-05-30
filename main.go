package main

import (
	"spike-blockchain-server/config"
	"spike-blockchain-server/server"
)

func main() {
	config.Init()

	r := server.NewRouter()
	r.Run(":3000")
}
