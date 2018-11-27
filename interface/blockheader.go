package _interface

import (
	"time"

	"github.com/ioeXNetwork/ioeX.SPV/util"

	"github.com/ioeXNetwork/ioeX.MainChain/core"
	"github.com/ioeXNetwork/ioeX.Utility/common"
	"github.com/ioeXNetwork/ioeX.Utility/crypto"
)

func newBlockHeader() util.BlockHeader {
	return util.NewIOEXHeader(&core.Header{})
}

// GenesisHeader creates a specific genesis header by the given
// foundation address.
func GenesisHeader(foundation *common.Uint168) util.BlockHeader {
	// Genesis time
	genesisTime := time.Date(2017, time.December, 22, 10, 0, 0, 0, time.UTC)

	// header
	header := core.Header{
		Version:    core.BlockVersion,
		Previous:   common.EmptyHash,
		MerkleRoot: common.EmptyHash,
		Timestamp:  uint32(genesisTime.Unix()),
		Bits:       0x1d03ffff,
		Nonce:      core.GenesisNonce,
		Height:     uint32(0),
	}

	// IOEX coin
	ioeXCoin := &core.Transaction{
		TxType:         core.RegisterAsset,
		PayloadVersion: 0,
		Payload: &core.PayloadRegisterAsset{
			Asset: core.Asset{
				Name:      "IOEX",
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

	coinBase := &core.Transaction{
		TxType:         core.CoinBase,
		PayloadVersion: core.PayloadCoinBaseVersion,
		Payload:        new(core.PayloadCoinBase),
		Inputs: []*core.Input{
			{
				Previous: core.OutPoint{
					TxID:  common.EmptyHash,
					Index: 0x0000,
				},
				Sequence: 0x00000000,
			},
		},
		Attributes: []*core.Attribute{},
		LockTime:   0,
		Programs:   []*core.Program{},
	}

	coinBase.Outputs = []*core.Output{
		{
			AssetID:     ioeXCoin.Hash(),
			Value:       20000 * 10000 * 100000000,
			ProgramHash: *foundation,
		},
	}

	nonce := []byte{0x4d, 0x65, 0x82, 0x21, 0x07, 0xfc, 0xfd, 0x52}
	txAttr := core.NewAttribute(core.Nonce, nonce)
	coinBase.Attributes = append(coinBase.Attributes, &txAttr)

	transactions := []*core.Transaction{coinBase, ioeXCoin}
	hashes := make([]common.Uint256, 0, len(transactions))
	for _, tx := range transactions {
		hashes = append(hashes, tx.Hash())
	}
	header.MerkleRoot, _ = crypto.ComputeRoot(hashes)

	return util.NewIOEXHeader(&header)
}
