package main

import (
	"os"
	"spike-blockchain-server/chain"
	"spike-blockchain-server/config"
	"spike-blockchain-server/server"
)

func main() {
	config.Init()
	bscClient, _ := chain.NewBscListener(os.Getenv("MORALIS_SPEEDY_NODE"), []string{"addr1"})
	bscClient.Run()

	r := server.NewRouter()
	r.Run(":3000")
}
