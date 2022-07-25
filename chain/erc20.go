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
	"strings"
	"time"
)

type SKSTarget struct {
	txAddress string
}

func newSKSTarget(address string) *SKSTarget {
	return &SKSTarget{
		txAddress: address,
	}
}

func (t *SKSTarget) Accept(fromAddr, toAddr string) (bool, uint64) {
	if strings.ToLower(t.txAddress) == strings.ToLower(toAddr) {
		return true, SKS_RECHARGE
	}

	if strings.ToLower(t.txAddress) == strings.ToLower(fromAddr) {
		return true, SKS_WITHDRAW
	}

	return false, NOT_EXIST
}

type SKKTarget struct {
	txAddress string
}

func newSKKTarget(address string) *SKKTarget {
	return &SKKTarget{
		txAddress: address,
	}
}

func (t *SKKTarget) Accept(fromAddr, toAddr string) (bool, uint64) {
	if strings.ToLower(t.txAddress) == strings.ToLower(toAddr) {
		return true, SKK_RECHARGE
	}

	if strings.ToLower(t.txAddress) == strings.ToLower(fromAddr) {
		return true, SKK_WITHDRAW
	}

	return false, NOT_EXIST
}

type ERC20Listener struct {
	TxFilter
	contractAddr   string
	cacheBlockNum  string
	erc20Notify    chan ERC20Tx
	rechargeNotify chan ERC20Tx
	ec             *ethclient.Client
	rc             *redis.Client
	abi            abi.ABI
}

func newERC20Listener(filter TxFilter, contractAddr string, cacheBlockNum string, ec *ethclient.Client, rc *redis.Client, erc20Notify chan ERC20Tx, rechargeNotify chan ERC20Tx, abi abi.ABI) *ERC20Listener {
	return &ERC20Listener{
		filter,
		contractAddr,
		cacheBlockNum,
		erc20Notify,
		rechargeNotify,
		ec,
		rc,
		abi,
	}
}

func (el *ERC20Listener) run() {
	go el.NewEventFilter(el.contractAddr, el.TxFilter)
	if el.rc.Get(el.cacheBlockNum).Err() == redis.Nil {
		log.Infof("%s is not exist", el.cacheBlockNum)
		return
	}
	nowBlockNum, err := el.ec.BlockNumber(context.Background())
	if err != nil {
		log.Errorf("query now %s blockNum err : %+v", el.cacheBlockNum, err)
		return
	}

	blockNum, err := el.rc.Get(el.cacheBlockNum).Uint64()
	if blockNum < nowBlockNum {
		el.PastEventFilter(el.contractAddr, blockNum, nowBlockNum, el.TxFilter)
	}
}

func (el *ERC20Listener) NewEventFilter(contractAddr string, filter TxFilter) error {
	newEventChan := make(chan types.Log)
	ethClient := el.ec
	contractAddress := common.HexToAddress(contractAddr)
	query := ethereum.FilterQuery{
		Addresses: []common.Address{contractAddress},
	}
	sub, err := ethClient.SubscribeFilterLogs(context.Background(), query, newEventChan)
	if err != nil {
		//log
		return err
	}
	for {
		select {
		case err := <-sub.Err():
			log.Errorf("erc 20 %s sub err : %+v", el.contractAddr, err)
			sub = event.Resubscribe(time.Millisecond, func(ctx context.Context) (event.Subscription, error) {
				return ethClient.SubscribeFilterLogs(context.Background(), query, newEventChan)
			})
		case log := <-newEventChan:
			switch log.Topics[0].String() {
			case eventSignHash(TransferTopic):
				intr, err := el.abi.Events["Transfer"].Inputs.Unpack(log.Data)
				if err != nil {
					break
				}
				fromAddr := common.HexToAddress(log.Topics[1].Hex()).String()
				toAddr := common.HexToAddress(log.Topics[2].Hex()).String()
				accept, txType := el.Accept(fromAddr, toAddr)
				if !accept {
					continue
				}
				var status uint64
				recp, err := el.ec.TransactionReceipt(context.Background(), log.TxHash)
				status = recp.Status
				if err != nil {
					status = 0
				}
				block, err := el.ec.BlockByNumber(context.Background(), big.NewInt(int64(log.BlockNumber)))
				if err != nil {
					status = 0
				}
				tx := ERC20Tx{
					From:    fromAddr,
					To:      toAddr,
					TxType:  txType,
					TxHash:  log.TxHash.Hex(),
					Status:  status,
					PayTime: int64(block.Time() * 1000),
					Amount:  intr[0].(*big.Int).String(),
				}

				if check(int(txType)) {
					el.rechargeNotify <- tx
				} else {
					el.erc20Notify <- tx
				}
			}
			el.rc.Set(el.cacheBlockNum, log.BlockNumber, 0)
		}
	}
}

func (el *ERC20Listener) PastEventFilter(contractAddr string, fromBlockNum, toBlockNum uint64, filter TxFilter) error {
	log.Infof("erc20 past event filter, type : %s, fromBlock : %d, toBlock : %d ", el.cacheBlockNum, fromBlockNum, toBlockNum)
	ethClient := el.ec
	contractAddress := common.HexToAddress(contractAddr)

	query := ethereum.FilterQuery{
		Addresses: []common.Address{contractAddress},
		FromBlock: big.NewInt(int64(fromBlockNum)),
		ToBlock:   big.NewInt(int64(toBlockNum)),
	}

	sub, err := ethClient.FilterLogs(context.Background(), query)
	if err != nil {
		log.Error("erc20 subscribe err : ", err)
		return err
	}
	for _, logEvent := range sub {
		switch logEvent.Topics[0].String() {
		case eventSignHash(TransferTopic):
			intr, err := el.abi.Events["Transfer"].Inputs.Unpack(logEvent.Data)
			if err != nil {
				log.Error("erc20 data unpack err : ", err)
				break
			}
			fromAddr := common.HexToAddress(logEvent.Topics[1].Hex()).String()
			toAddr := common.HexToAddress(logEvent.Topics[2].Hex()).String()
			accept, txType := el.Accept(fromAddr, toAddr)
			if !accept {
				continue
			}
			var status uint64
			recp, err := el.ec.TransactionReceipt(context.Background(), logEvent.TxHash)
			status = recp.Status
			if err != nil {
				status = 0
			}
			block, err := el.ec.BlockByNumber(context.Background(), big.NewInt(int64(logEvent.BlockNumber)))
			if err != nil {
				status = 0
			}
			el.erc20Notify <- ERC20Tx{
				From:    fromAddr,
				To:      toAddr,
				TxType:  txType,
				TxHash:  logEvent.TxHash.Hex(),
				Status:  status,
				PayTime: int64(block.Time() * 1000),
				Amount:  intr[0].(*big.Int).String(),
			}
		}
	}
	return err
}
