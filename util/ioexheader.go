package util

import (
	"github.com/ioeXNetwork/ioeX.MainChain/core"
	"github.com/ioeXNetwork/ioeX.Utility/common"
)

// Ensure SideHeader implement BlockHeader interface.
var _ BlockHeader = (*IOEXHeader)(nil)

type IOEXHeader struct {
	*core.Header
}

func (h *IOEXHeader) Previous() common.Uint256 {
	return h.Header.Previous
}

func (h *IOEXHeader) Bits() uint32 {
	return h.Header.Bits
}

func (h *IOEXHeader) MerkleRoot() common.Uint256 {
	return h.Header.MerkleRoot
}

func (h *IOEXHeader) PowHash() common.Uint256 {
	return h.AuxPow.ParBlockHeader.Hash()
}

func NewIOEXHeader(orgHeader *core.Header) BlockHeader {
	return &IOEXHeader{Header: orgHeader}
}
