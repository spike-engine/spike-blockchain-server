package chain

const (
	SKK TokenType = iota
	SKS
	USDC
	BNB
	AUGT
	AUNFT
)

const (
	SKK_RECHARGE = iota + 1
	SKS_RECHARGE
	USDC_RECHARGE
	BNB_RECHARGE
	SKK_WITHDRAW
	SKS_WITHDRAW
	USDC_WITHDRAW
	BNB_WITHDRAW
	AUNFT_TRANSFER
	AUGT_RECHARGE
	AUGT_WITHDRAW
	NOT_EXIST
)

var recharge = map[int]struct{}{
	SKK_RECHARGE:  {},
	SKS_RECHARGE:  {},
	USDC_RECHARGE: {},
	BNB_RECHARGE:  {},
}

type TokenType int

type TxFilter interface {
	Accept(fromAddr, toAddr string) (bool, uint64)
}

func check(txType int) bool {
	_, ok := recharge[txType]
	return ok
}
