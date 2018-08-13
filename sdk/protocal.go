package sdk

import (
	"time"

	"github.com/ioeX/ioeX.Utility/common"
	"github.com/ioeX/ioeX.Utility/crypto"
	"github.com/ioeX/ioeX.Utility/p2p"

	ioex "github.com/ioeX/ioeX.MainChain/core"
)

const (
	ProtocolVersion = p2p.EIP001Version // The protocol version implemented SPV protocol
	MaxMsgSize      = 1024 * 1024 * 8   // The max size of a message in P2P network
	OpenService     = 1 << 2
)

func GenesisHeader(foundation *common.Uint168) *ioex.Header {
	// header
	header := ioex.Header{
		Version:    ioex.BlockVersion,
		Previous:   common.EmptyHash,
		MerkleRoot: common.EmptyHash,
		Timestamp:  uint32(time.Unix(time.Date(2017, time.December, 22, 10, 0, 0, 0, time.UTC).Unix(), 0).Unix()),
		Bits:       0x1d03ffff,
		Nonce:      ioex.GenesisNonce,
		Height:     uint32(0),
	}

	// IOEX coin
	ioeXCoin := &ioex.Transaction{
		TxType:         ioex.RegisterAsset,
		PayloadVersion: 0,
		Payload: &ioex.PayloadRegisterAsset{
			Asset: ioex.Asset{
				Name:      "IOEX",
				Precision: 0x08,
				AssetType: 0x00,
			},
			Amount:     0 * 100000000,
			Controller: common.Uint168{},
		},
		Attributes: []*ioex.Attribute{},
		Inputs:     []*ioex.Input{},
		Outputs:    []*ioex.Output{},
		Programs:   []*ioex.Program{},
	}

	coinBase := &ioex.Transaction{
		TxType:         ioex.CoinBase,
		PayloadVersion: ioex.PayloadCoinBaseVersion,
		Payload:        new(ioex.PayloadCoinBase),
		Inputs: []*ioex.Input{
			{
				Previous: ioex.OutPoint{
					TxID:  common.EmptyHash,
					Index: 0x0000,
				},
				Sequence: 0x00000000,
			},
		},
		Attributes: []*ioex.Attribute{},
		LockTime:   0,
		Programs:   []*ioex.Program{},
	}

	coinBase.Outputs = []*ioex.Output{
		{
			AssetID:     ioeXCoin.Hash(),
			Value:       3300 * 10000 * 100000000,
			ProgramHash: *foundation,
		},
	}

	nonce := []byte{0x4d, 0x65, 0x82, 0x21, 0x07, 0xfc, 0xfd, 0x52}
	txAttr := ioex.NewAttribute(ioex.Nonce, nonce)
	coinBase.Attributes = append(coinBase.Attributes, &txAttr)

	transactions := []*ioex.Transaction{coinBase, ioeXCoin}
	hashes := make([]common.Uint256, 0, len(transactions))
	for _, tx := range transactions {
		hashes = append(hashes, tx.Hash())
	}
	header.MerkleRoot, _ = crypto.ComputeRoot(hashes)

	return &header
}
