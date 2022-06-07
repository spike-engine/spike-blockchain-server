package middleware

import (
	"fmt"
	"testing"
)

func TestETHParamsVerify(t *testing.T) {
	params := ethParams{
		Chain: "eth",
		Sign:  "0x64e34ad9914a2decc48ee445aea6ad6155789a6f1b912048500a81f5bbbb4411721ea559cc18fab0758d35918dfae1aa4adec41805a87caa11e817f1698d76681c",
		Nonce: "hello, spike",
		Addr:  "0x43a0D6FcD600f061B38e80D6580Ab900Ed67dFF1",
	}
	fmt.Println("ETH params.verify:", params.verify())
}
