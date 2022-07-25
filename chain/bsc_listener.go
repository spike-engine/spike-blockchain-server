package chain

import (
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/go-redis/redis"
	logger "github.com/ipfs/go-log"
	"os"
	"spike-blockchain-server/cache"
	"spike-blockchain-server/game"
)

var log = logger.Logger("chain")

const (
	BNB_BLOCKNUM   = "bnb_blockNum"
	USDC_BLOCKNUM  = "usdc_blockNum"
	SKK_BLOCKNUM   = "skk_blockNum"
	SKS_BLOCKNUM   = "sks_blockNum"
	AUNFT_BLOCKNUM = "aunft_blockNum"
)

type BscListener struct {
	ec *ethclient.Client
	rc *redis.Client
	l  map[TokenType]Listener
}

func NewBscListener(speedyNodeAddress string, targetWalletAddr string) (*BscListener, error) {
	log.Infof("bsc listener start")
	bl := &BscListener{}

	client, err := ethclient.Dial(speedyNodeAddress)
	if err != nil {
		log.Error("eth client dial err : ", err)
		return nil, err
	}
	bl.rc = cache.RedisClient
	bl.ec = client
	erc20Notify := make(chan ERC20Tx, 10)
	rechargeNotify := make(chan ERC20Tx, 10)
	erc721Notify := make(chan ERC721Tx, 10)
	l := make(map[TokenType]Listener)
	l[BNB] = newBNBListener(newBNBTarget(targetWalletAddr), bl.ec, bl.rc, erc20Notify, rechargeNotify)
	l[SKK] = newERC20Listener(newSKKTarget(targetWalletAddr), SKKContractAddress, SKK_BLOCKNUM, bl.ec, bl.rc, erc20Notify, rechargeNotify, getABI(SKKContractAbi))
	l[SKS] = newERC20Listener(newSKSTarget(targetWalletAddr), SKSContractAddress, SKS_BLOCKNUM, bl.ec, bl.rc, erc20Notify, rechargeNotify, getABI(SKSContractAbi))
	l[AUNFT] = newAUNFTListener(bl.ec, bl.rc, erc721Notify, getABI(AUNFTAbi))
	bl.l = l
	spikeTxMgr := newSpikeTxMgr(game.NewKafkaClient(os.Getenv("KAFKA_ADDR")), erc20Notify, rechargeNotify, erc721Notify)
	go spikeTxMgr.run()
	return bl, nil
}

func (bl *BscListener) Run() {
	for _, listener := range bl.l {
		go func(l Listener) {
			l.run()
		}(listener)

	}
}
