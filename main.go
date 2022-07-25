package main

import (
	logger "github.com/ipfs/go-log"
	"os"
	"spike-blockchain-server/chain"
	"spike-blockchain-server/config"
	"spike-blockchain-server/server"
)

func main() {
	logger.SetLogLevel("*", "INFO")
	config.Init()
	bscClient, err := chain.NewBscListener(os.Getenv("MORALIS_SPEEDY_NODE"), os.Getenv("TARGET_ADDR"))
	if err != nil {
		//log
		return
	}
	bscClient.Run()

	r := server.NewRouter(bscClient)
	r.Run(":3000")
}
