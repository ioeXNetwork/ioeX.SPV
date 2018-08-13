package db

import (
	"github.com/ioeX/ioeX.Utility/common"
)

type QueueItem struct {
	NotifyId common.Uint256
	TxId     common.Uint256
	Height   uint32
}
