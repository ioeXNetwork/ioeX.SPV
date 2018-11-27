package sync

import (
	"github.com/ioeXNetwork/ioeX.SPV/blockchain"
	"github.com/ioeXNetwork/ioeX.SPV/bloom"
	"github.com/ioeXNetwork/ioeX.SPV/util"
)

const (
	defaultMaxPeers = 125
)

// Config is a configuration struct used to initialize a new SyncManager.
type Config struct {
	Chain *blockchain.BlockChain

	MaxPeers int

	UpdateFilter        func() *bloom.Filter
	TransactionAnnounce func(tx util.Transaction)
}

func NewDefaultConfig(chain *blockchain.BlockChain,
	updateFilter func() *bloom.Filter) *Config {
	return &Config{
		Chain:        chain,
		MaxPeers:     defaultMaxPeers,
		UpdateFilter: updateFilter,
	}
}
