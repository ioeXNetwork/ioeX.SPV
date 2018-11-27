package main

import (
	"fmt"

	"github.com/ioeXNetwork/ioeX.SPV/util"
	"github.com/ioeXNetwork/ioeX.SPV/wallet"

	"github.com/ioeXNetwork/ioeX.MainChain/core"
	"github.com/ioeXNetwork/ioeX.Utility/common"
)

var Version string

func main() {
	url := fmt.Sprint("http://127.0.0.1:", config.JsonRpcPort, "/spvwallet")
	wallet.RunClient(Version, url, getSystemAssetId(), func() util.BlockHeader {
		return util.NewIOEXHeader(&core.Header{})
	})
}

func getSystemAssetId() common.Uint256 {
	systemToken := &core.Transaction{
		TxType:         core.RegisterAsset,
		PayloadVersion: 0,
		Payload: &core.PayloadRegisterAsset{
			Asset: core.Asset{
				Name:      "ELA",
				Precision: 0x08,
				AssetType: 0x00,
			},
			Amount:     0 * 100000000,
			Controller: common.Uint168{},
		},
		Attributes: []*core.Attribute{},
		Inputs:     []*core.Input{},
		Outputs:    []*core.Output{},
		Programs:   []*core.Program{},
	}
	return systemToken.Hash()
}
