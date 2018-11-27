package main

import (
	"io"
	"os"

	"github.com/ioeXNetwork/ioeX.SPV/blockchain"
	"github.com/ioeXNetwork/ioeX.SPV/peer"
	"github.com/ioeXNetwork/ioeX.SPV/sdk"
	"github.com/ioeXNetwork/ioeX.SPV/sync"
	"github.com/ioeXNetwork/ioeX.SPV/wallet"
	"github.com/ioeXNetwork/ioeX.SPV/wallet/store"

	"github.com/ioeXNetwork/ioeX.Utility/elalog"
	"github.com/ioeXNetwork/ioeX.Utility/http/jsonrpc"
	"github.com/ioeXNetwork/ioeX.Utility/p2p/addrmgr"
	"github.com/ioeXNetwork/ioeX.Utility/p2p/connmgr"
	"github.com/ioeXNetwork/ioeX.Utility/p2p/server"
)

const LogPath = "./logs-spv/"

// log is a logger that is initialized with no output filters.  This
// means the package will not perform any logging by default until the caller
// requests it.
var (
	fileWriter = elalog.NewFileWriter(
		LogPath,
		config.MaxPerLogSize,
		config.MaxLogsSize,
	)
	level   = elalog.Level(config.PrintLevel)
	backend = elalog.NewBackend(io.MultiWriter(os.Stdout, fileWriter),
		elalog.Llongfile)

	admrlog = backend.Logger("ADMR", elalog.LevelOff)
	cmgrlog = backend.Logger("CMGR", elalog.LevelOff)
	bcdblog = backend.Logger("BCDB", level)
	synclog = backend.Logger("SYNC", level)
	peerlog = backend.Logger("PEER", level)
	spvslog = backend.Logger("SPVS", elalog.LevelInfo)
	srvrlog = backend.Logger("SRVR", level)
	rpcslog = backend.Logger("RPCS", level)
	waltlog = backend.Logger("WALT", level)
)

func init() {
	addrmgr.UseLogger(admrlog)
	connmgr.UseLogger(cmgrlog)
	blockchain.UseLogger(bcdblog)
	sdk.UseLogger(spvslog)
	jsonrpc.UseLogger(rpcslog)
	peer.UseLogger(peerlog)
	server.UseLogger(srvrlog)
	store.UseLogger(bcdblog)
	sync.UseLogger(synclog)
	wallet.UseLogger(waltlog)
}
