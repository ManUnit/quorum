package raft

import (
	"github.com/ethereum/quorum/core/types"
)

type InvalidRaftOrdering struct {
	// Current head of the chain
	headBlock *types.Block

	// New block that should point to the head, but doesn't
	invalidBlock *types.Block
}
