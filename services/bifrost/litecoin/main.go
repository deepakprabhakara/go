package litecoin

import (
	"math/big"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
	"github.com/stellar/go/support/errors"
	"github.com/stellar/go/support/log"
	"github.com/tyler-smith/go-bip32"
)

var (
	eight = big.NewInt(8)
	ten   = big.NewInt(10)
	// litInLtc = 10^8
	litInLtc = new(big.Rat).SetInt(new(big.Int).Exp(ten, eight, nil))
)

// Listener listens for transactions using litecoin-core RPC. It calls TransactionHandler for each new
// transactions. It will reprocess the block if TransactionHandler returns error. It will
// start from the block number returned from Storage.GetLitecoinBlockToProcess or the latest block
// if it returned 0. Transactions can be processed more than once, it's TransactionHandler
// responsibility to ignore duplicates.
// Listener tracks only P2PKH payments.
// You can run multiple Listeners if Storage is implemented correctly.
type Listener struct {
	Enabled            bool
	Client             Client  `inject:""`
	Storage            Storage `inject:""`
	TransactionHandler TransactionHandler
	Testnet            bool

	chainParams *chaincfg.Params
	log         *log.Entry
}

type Client interface {
	GetBlockCount() (int64, error)
	GetBlockHash(blockHeight int64) (*chainhash.Hash, error)
	GetBlock(blockHash *chainhash.Hash) (*wire.MsgBlock, error)
}

// Storage is an interface that must be implemented by an object using
// persistent storage.
type Storage interface {
	// GetLitecoinBlockToProcess gets the number of Litecoin block to process. `0` means the
	// processing should start from the current block.
	GetLitecoinBlockToProcess() (uint64, error)
	// SaveLastProcessedLitecoinBlock should update the number of the last processed Litecoin
	// block. It should only update the block if block > current block in atomic transaction.
	SaveLastProcessedLitecoinBlock(block uint64) error
}

type TransactionHandler func(transaction Transaction) error

type Transaction struct {
	Hash       string
	TxOutIndex int
	// Value in lits
	ValueLit int64
	To       string
}

type AddressGenerator struct {
	masterPublicKey *bip32.Key
	chainParams     *chaincfg.Params
}

func LtcToLit(ltc string) (int64, error) {
	valueRat := new(big.Rat)
	_, ok := valueRat.SetString(ltc)
	if !ok {
		return 0, errors.New("Could not convert to *big.Rat")
	}

	// Calculate value in litoshi
	valueRat.Mul(valueRat, litInLtc)

	// Ensure denominator is equal `1`
	if valueRat.Denom().Cmp(big.NewInt(1)) != 0 {
		return 0, errors.New("Invalid precision, is value smaller than 1 litoshi?")
	}

	return valueRat.Num().Int64(), nil
}
