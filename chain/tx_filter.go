package chain

import "github.com/ethereum/go-ethereum/core/types"

type TxFilter interface {
	Accept(transaction *types.Transaction) bool
}
