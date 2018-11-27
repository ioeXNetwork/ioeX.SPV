package bloom

import (
	"fmt"
	"io"

	"github.com/ioeXNetwork/ioeX.Utility/common"
	"github.com/ioeXNetwork/ioeX.Utility/p2p/msg"
)

// maxFlagsPerMerkleProof is the maximum number of flag bytes that could
// possibly fit into a merkle proof.  Since each transaction is represented by
// a single bit, this is the max number of transactions per block divided by
// 8 bits per byte.  Then an extra one to cover partials.
const maxFlagsPerMerkleProof = msg.MaxTxPerBlock / 8

type MerkleProof struct {
	BlockHash    common.Uint256
	Height       uint32
	Transactions uint32
	Hashes       []*common.Uint256
	Flags        []byte
}

func (p *MerkleProof) Serialize(w io.Writer) error {
	// Read num transaction hashes and limit to max.
	numHashes := len(p.Hashes)
	if numHashes > msg.MaxTxPerBlock {
		str := fmt.Sprintf("too many transaction hashes for message "+
			"[count %v, max %v]", numHashes, msg.MaxTxPerBlock)
		return common.FuncError("MerkleProof.Serialize", str)
	}
	numFlagBytes := len(p.Flags)
	if numFlagBytes > maxFlagsPerMerkleProof {
		str := fmt.Sprintf("too many flag bytes for message [count %v, "+
			"max %v]", numFlagBytes, maxFlagsPerMerkleProof)
		return common.FuncError("MerkleProof.Serialize", str)
	}

	err := common.WriteElements(w, &p.BlockHash, p.Height, p.Transactions,
		uint32(numHashes))
	if err != nil {
		return err
	}

	for _, hash := range p.Hashes {
		if err := hash.Serialize(w); err != nil {
			return err
		}
	}

	return common.WriteVarBytes(w, p.Flags)
}

func (p *MerkleProof) Deserialize(r io.Reader) error {
	var numHashes uint32
	err := common.ReadElements(r,
		&p.BlockHash,
		&p.Height,
		&p.Transactions,
		&numHashes,
	)
	if err != nil {
		return err
	}

	if numHashes > msg.MaxTxPerBlock {
		return fmt.Errorf("MerkleProof.Deserialize too many transaction"+
			" hashes for message [count %v, max %v]", numHashes, msg.MaxTxPerBlock)
	}

	hashes := make([]common.Uint256, numHashes)
	p.Hashes = make([]*common.Uint256, 0, numHashes)
	for i := uint32(0); i < numHashes; i++ {
		hash := &hashes[i]
		if err := hash.Deserialize(r); err != nil {
			return err
		}
		p.Hashes = append(p.Hashes, hash)
	}

	p.Flags, err = common.ReadVarBytes(r, maxFlagsPerMerkleProof,
		"merkle proof flags size")
	return err
}
