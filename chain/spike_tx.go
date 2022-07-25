package chain

import (
	"encoding/json"
	"fmt"
	"spike-blockchain-server/game"
)

type ERC20Tx struct {
	From    string `json:"from"`
	To      string `json:"to"`
	TxType  uint64 `json:"txType"`
	TxHash  string `json:"txHash"`
	Status  uint64 `json:"status"`
	PayTime int64  `json:"payTime"`
	Amount  string `json:"amount"`
}

type ERC721Tx struct {
	From    string `json:"from"`
	To      string `json:"to"`
	TxType  uint64 `json:"txType"`
	TxHash  string `json:"txHash"`
	Status  uint64 `json:"status"`
	PayTime int64  `json:"payTime"`
	TokenId uint64 `json:"tokenId"`
}

type SpikeTxMgr struct {
	erc20Notify    chan ERC20Tx
	erc721Notify   chan ERC721Tx
	rechargeNotify chan ERC20Tx
	close          chan struct{}
	mqApi          game.MqApi
}

func newSpikeTxMgr(client *game.KafkaClient, erc20Notify chan ERC20Tx, rechargeNotify chan ERC20Tx, erc721Notify chan ERC721Tx) *SpikeTxMgr {
	s := &SpikeTxMgr{
		erc20Notify:    erc20Notify,
		rechargeNotify: rechargeNotify,
		erc721Notify:   erc721Notify,
		mqApi:          client,
	}

	return s
}

func (s *SpikeTxMgr) run() {
	for {
		select {
		case erc20Tx := <-s.erc20Notify:
			txByte, err := json.Marshal(erc20Tx)
			if err != nil {
				log.Errorf("json marshal err : %+v, txByte : %s", err, txByte)
				break
			}
			log.Infof("erc20 value : %s", string(txByte))
			err = s.mqApi.SendMessage(game.Msg{
				Topic: game.ERC20TXTOPIC,
				Key:   erc20Tx.TxHash,
				Value: string(txByte),
			})
			if err != nil {
				fmt.Println("", err)
			}
		case erc20Tx := <-s.rechargeNotify:
			txByte, err := json.Marshal(erc20Tx)
			if err != nil {
				log.Error("json marshal err : ", err)
				break
			}
			log.Infof("value : %s", string(txByte))
			err = s.mqApi.SendMessage(game.Msg{
				Topic: game.RECHARGETXTOPIC,
				Key:   erc20Tx.TxHash,
				Value: string(txByte),
			})
			if err != nil {
				log.Error("", err)
			}
		case erc721Tx := <-s.erc721Notify:
			txByte, err := json.Marshal(erc721Tx)
			log.Infof("value : %s", string(txByte))
			if err != nil {
				log.Error(err)
				break
			}
			err = s.mqApi.SendMessage(game.Msg{
				Topic: game.ERC721TXTOPIC,
				Key:   erc721Tx.TxHash,
				Value: string(txByte),
			})
			if err != nil {
				log.Error("", err)
			}
		case <-s.close:
			//log
			return
		}
	}

}
