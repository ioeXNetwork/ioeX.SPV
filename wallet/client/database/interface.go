package database

import (
	"github.com/ioeXNetwork/ioeX.SPV/wallet/sutil"

	"github.com/ioeXNetwork/ioeX.Utility/common"
)

type Database interface {
	AddAddress(address *common.Uint168, script []byte, addrType int) error
	GetAddress(address *common.Uint168) (*sutil.Addr, error)
	GetAddrs() ([]*sutil.Addr, error)
	DeleteAddress(address *common.Uint168) error
	GetAddressUTXOs(address *common.Uint168) ([]*sutil.UTXO, error)
	GetAddressSTXOs(address *common.Uint168) ([]*sutil.STXO, error)
	BestHeight() uint32
	Clear() error
}
