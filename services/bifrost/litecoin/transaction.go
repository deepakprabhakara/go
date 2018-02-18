package litecoin

import (
	"math/big"

	"github.com/stellar/go/services/bifrost/common"
)

func (t Transaction) ValueToStellar() string {
	ValueLit := new(big.Int).SetInt64(t.ValueLit)
	valueLtc := new(big.Rat).Quo(new(big.Rat).SetInt(ValueLit), litInLtc)
	return valueLtc.FloatString(common.StellarAmountPrecision)
}
