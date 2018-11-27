package util

import (
	"io"

	"github.com/ioeXNetwork/ioeX.Utility/common"
)

type BlockHeader interface {
	Previous() common.Uint256
	Bits() uint32
	MerkleRoot() common.Uint256
	Hash() common.Uint256
	PowHash() common.Uint256
	Serialize(w io.Writer) error
	Deserialize(r io.Reader) error
}

type Transaction interface {
	Hash() common.Uint256
	Serialize(w io.Writer) error
	Deserialize(r io.Reader) error
}
