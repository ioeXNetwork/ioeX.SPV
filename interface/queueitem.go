package _interface

import (
	. "github.com/ioeXNetwork/ioeX.Utility/common"
)

type QueueItem struct {
	TxHash    Uint256
	BlockHash Uint256
	Height    uint32
}
