package chain

import (
	"context"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/event"
	"github.com/go-redis/redis"
	"math/big"
	"time"
)

type AUNFTListener struct {
	erc721Notify chan ERC721Tx
	ec           *ethclient.Client
	rc           *redis.Client
	abi          abi.ABI
}

func newAUNFTListener(ec *ethclient.Client, rc *redis.Client, erc721Notify chan ERC721Tx, abi abi.ABI) *AUNFTListener {
	return &AUNFTListener{
		erc721Notify,
		ec,
		rc,
		abi,
	}
}

func (al *AUNFTListener) run() {
	go al.NewEventFilter(AUNFTContractAddress)
	if al.rc.Get(AUNFT_BLOCKNUM).Err() == redis.Nil {
		log.Infof("aunft_blockNum is not exist")
		return
	}
	nowBlockNum, err := al.ec.BlockNumber(context.Background())
	if err != nil {
		log.Errorf("query now aunft_blockNum err : %+v", err)
		return
	}
	blockNum, err := al.rc.Get(AUNFT_BLOCKNUM).Uint64()
	if blockNum < nowBlockNum {
		al.PastEventFilter(AUNFTContractAddress, blockNum, nowBlockNum)
	}
}

func (al *AUNFTListener) NewEventFilter(contractAddr string) error {
	newEventChan := make(chan types.Log)
	ethClient := al.ec
	contractAddress := common.HexToAddress(contractAddr)
	query := ethereum.FilterQuery{
		Addresses: []common.Address{contractAddress},
	}
	sub, err := ethClient.SubscribeFilterLogs(context.Background(), query, newEventChan)
	if err != nil {
		log.Error("nft subscribeFilterLogs err : ", err)
		return err
	}
	for {
		select {
		case err = <-sub.Err():
			log.Error("nft sub err : ", err)
			sub = event.Resubscribe(time.Millisecond, func(ctx context.Context) (event.Subscription, error) {
				return ethClient.SubscribeFilterLogs(context.Background(), query, newEventChan)
			})
		case l := <-newEventChan:
			switch l.Topics[0].String() {
			case eventSignHash(TransferTopic):
				var status uint64
				recp, err := al.ec.TransactionReceipt(context.Background(), l.TxHash)
				status = recp.Status
				if err != nil {
					log.Error("nft TransactionReceipt err : ", err)
					status = 0
				}
				block, err := al.ec.BlockByNumber(context.Background(), big.NewInt(int64(l.BlockNumber)))
				if err != nil {
					status = 0
				}
				fromAddr := common.HexToAddress(l.Topics[1].Hex()).String()
				toAddr := common.HexToAddress(l.Topics[2].Hex()).String()

				al.erc721Notify <- ERC721Tx{
					From:    fromAddr,
					To:      toAddr,
					TxType:  AUNFT_TRANSFER,
					TxHash:  l.TxHash.Hex(),
					Status:  status,
					PayTime: int64(block.Time() * 1000),
					TokenId: l.Topics[3].Big().Uint64(),
				}
			}
			al.rc.Set(AUNFT_BLOCKNUM, l.BlockNumber, 0)
		}
	}
}

func (al *AUNFTListener) PastEventFilter(contractAddr string, fromBlockNum, toBlockNum uint64) error {
	ethClient := al.ec
	contractAddress := common.HexToAddress(contractAddr)

	query := ethereum.FilterQuery{
		Addresses: []common.Address{contractAddress},
		FromBlock: big.NewInt(int64(fromBlockNum)),
		ToBlock:   big.NewInt(int64(toBlockNum)),
	}

	sub, err := ethClient.FilterLogs(context.Background(), query)
	if err != nil {
		log.Error("nft subscribe event log, err : ", err)
		return err
	}
	for _, l := range sub {
		switch l.Topics[0].String() {
		case eventSignHash(TransferTopic):
			var status uint64
			recp, err := al.ec.TransactionReceipt(context.Background(), l.TxHash)
			status = recp.Status
			if err != nil {
				log.Error("nft TransactionReceipt err : ", err)
				status = 0
			}
			block, err := al.ec.BlockByNumber(context.Background(), big.NewInt(int64(l.BlockNumber)))
			if err != nil {
				status = 0
			}

			fromAddr := common.HexToAddress(l.Topics[1].Hex()).String()
			toAddr := common.HexToAddress(l.Topics[2].Hex()).String()
			al.erc721Notify <- ERC721Tx{
				From:    fromAddr,
				To:      toAddr,
				TxType:  AUNFT_TRANSFER,
				TxHash:  l.TxHash.Hex(),
				Status:  status,
				PayTime: int64(block.Time() * 1000),
				TokenId: l.Topics[3].Big().Uint64(),
			}
		}
	}
	return err
}
