package chain

import (
	"encoding/json"
	"github.com/ethereum/go-ethereum/common"
	"spike-blockchain-server/game"
	"sync"
)

type SpikeTx struct {
	from   string
	to     string
	value  string
	txHash common.Hash
}

type SpikeTxMgr interface {
	AddSpikeTx(blockNum uint64, st SpikeTx) error
	RemoveSpikeTx(blockNum uint64, hash common.Hash) error
	NotifySpikeTx() error
}

type spikeTxMgr struct {
	txs   map[uint64][]SpikeTx //blockNum
	lock  sync.Mutex
	mqApi game.MqApi
}

func newSpikeTxMgr(client *game.KafkaClient) *spikeTxMgr {
	txs := make(map[uint64][]SpikeTx)
	return &spikeTxMgr{
		txs:   txs,
		mqApi: client,
	}
}

func (s *spikeTxMgr) AddSpikeTx(blockNum uint64, st SpikeTx) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	if _, have := s.txs[blockNum]; have {
		spikeTxList := append(s.txs[blockNum], st)
		s.txs[blockNum] = spikeTxList
		return nil
	}
	var spikeTxList []SpikeTx
	spikeTxList = append(spikeTxList, st)
	return nil
}

func (s *spikeTxMgr) RemoveSpikeTx(blockNum uint64, hash common.Hash) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	var remainingTx []SpikeTx
	if txs, have := s.txs[blockNum]; have {
		for _, tx := range txs {
			if tx.txHash == hash {
				continue
			}
			remainingTx = append(remainingTx, tx)
		}
		s.txs[blockNum] = remainingTx
	}
	return nil
}

func (s *spikeTxMgr) NotifySpikeTx() error {
	s.lock.Lock()
	defer s.lock.Unlock()
	for blockNum, txs := range s.txs {
		if len(txs) == 0 {
			delete(s.txs, blockNum)
		}
		for _, tx := range txs {
			txByte, err := json.Marshal(tx)
			if err != nil {
				break
			}
			err = s.mqApi.SendMessage(game.Msg{
				Topic: game.TXTOPIC,
				Key:   string(blockNum),
				Value: string(txByte),
			})
			if err != nil {
				break
			}
			s.RemoveSpikeTx(blockNum, tx.txHash)
		}
	}

	return nil
}
