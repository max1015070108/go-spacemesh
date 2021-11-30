package system

import (
	"context"

	"github.com/spacemeshos/go-spacemesh/common/types"
)

//go:generate mockgen -package=mocks -destination=./mocks/fetcher.go -source=./fetcher.go

// Fetcher is a general interface that defines a component capable of fetching data from remote peers.
type Fetcher interface {
	AtxFetcher
	BlockFetcher
	PoetProofFetcher
	BallotFetcher
	TxFetcher
}

// BlockFetcher defines an interface for fetching blocks from remote peers.
type BlockFetcher interface {
	GetBlocks(context.Context, []types.BlockID) error
	FetchBlock(context.Context, types.BlockID) error
}

// AtxFetcher defines an interface for fetching ATXs from remote peers.
type AtxFetcher interface {
	GetAtxs(context.Context, []types.ATXID) error
	FetchAtx(context.Context, types.ATXID) error
}

// TxFetcher defines an interface for fetching transactions from remote peers.
type TxFetcher interface {
	GetTxs(context.Context, []types.TransactionID) error
}

// PoetProofFetcher defines an interface for fetching PoET proofs from remote peers.
type PoetProofFetcher interface {
	GetPoetProof(context.Context, types.Hash32) error
}

// BallotFetcher defines an interface for fetching Proposal and Ballot from remote peers.
type BallotFetcher interface {
	GetBallots(context.Context, []types.BallotID) error
}