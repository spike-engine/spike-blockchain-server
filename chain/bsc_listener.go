package chain

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/go-redis/redis"
	"math/big"
	"spike-blockchain-server/cache"
	"spike-blockchain-server/game"
)

const BLOCKNUM = "blockNum"

type BscClient struct {
	rc *redis.Client
	ec *ethclient.Client
	TxFilter
	SpikeTxMgr
}

func NewBscListener(speedyNodeAddress string, targetAddr []string) (*BscClient, error) {
	bl := &BscClient{}
	bl.rc = cache.RedisClient

	client, err := ethclient.Dial(speedyNodeAddress)
	if err != nil {
		return nil, err
	}
	bl.ec = client

	bl.TxFilter = newTarget(targetAddr)
	kafkaClient := game.NewKafkaClient([]string{"KAFKA_ADDRESS"})
	bl.SpikeTxMgr = newSpikeTxMgr(kafkaClient)
	return bl, nil
}

func (bl *BscClient) Run() {
	//not set the key
	if bl.rc.Get(BLOCKNUM).Err() != redis.Nil {
		nowBlockNum, err := bl.ec.BlockNumber(context.Background())
		if err != nil {
			//log
			return
		}
		blockNum, err := bl.rc.Get(BLOCKNUM).Uint64()
		if blockNum <= nowBlockNum {
			for i := blockNum + 1; i <= nowBlockNum; i++ {
				go bl.PastBlockFilter(i)
			}
		}
	}
	go bl.NewBlockFilter()
}

func (bl *BscClient) NewBlockFilter() {
	newBlockChan := make(chan *types.Header)
	ethClient := bl.ec
	sub, err := ethClient.SubscribeNewHead(context.Background(), newBlockChan)
	if err != nil {
		//log
		return
	}
	for {
		select {
		case _ = <-sub.Err():
			//log
		case header := <-newBlockChan:

			block, err := ethClient.BlockByHash(context.Background(), header.ParentHash)
			if err != nil {
				//log
				break
			}
			bl.SingleBlockFilter(block)
			bl.rc.Set(BLOCKNUM, block.Number().Uint64()-1, 0)
		}
	}
}

func (bl *BscClient) PastBlockFilter(blockNum uint64) error {
	block, err := bl.ec.BlockByNumber(context.Background(), big.NewInt(int64(blockNum)))
	if err != nil {
		return err
	}
	err = bl.SingleBlockFilter(block)
	bl.rc.Set(BLOCKNUM, block.Number().Uint64(), 0)
	return err
}

func (bl *BscClient) SingleBlockFilter(block *types.Block) error {
	for _, tx := range block.Transactions() {
		if !bl.Accept(tx) {
			continue
		}
		chainID, err := bl.ec.NetworkID(context.Background())
		if err != nil {
			//log
			return err
		}
		var fromAddr string
		if msg, err := tx.AsMessage(types.NewEIP155Signer(chainID), nil); err != nil {
			fromAddr = msg.From().Hex()
		}
		spikeTx := SpikeTx{
			from:   fromAddr,
			to:     tx.To().Hex(),
			value:  tx.Value().String(),
			txHash: tx.Hash(),
		}
		bl.AddSpikeTx(block.Number().Uint64(), spikeTx)
	}
	return nil
}

func (bl *BscClient) QueryTxByHash(txHash string) (value string, isPending bool, err error) {
	tx, isPending, err := bl.ec.TransactionByHash(context.Background(), common.BytesToHash([]byte(txHash)))
	return tx.Value().String(), isPending, err
}
