package chain

import "github.com/ethereum/go-ethereum/core/types"

type target struct {
	toAddress []string
}

func newTarget(toAddress []string) *target {
	return &target{
		toAddress: toAddress,
	}
}

func (t *target) Accept(tx *types.Transaction) bool {
	if tx.To() == nil {
		return false
	}
	for _, to := range t.toAddress {
		if to == tx.To().Hex() {
			return true
		}
	}
	return false
}
