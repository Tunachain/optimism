package store

import (
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/peerstore"
)

type TopicScores struct {
	TimeInMesh               float64 `json:"timeInMesh"` // in seconds
	FirstMessageDeliveries   uint64  `json:"firstMessageDeliveries"`
	MeshMessageDeliveries    uint64  `json:"meshMessageDeliveries"`
	InvalidMessageDeliveries uint64  `json:"invalidMessageDeliveries"`
}

type GossipScores struct {
	Total              float64     `json:"total"`
	Blocks             TopicScores `json:"blocks"` // fully zeroed if the peer has not been in the mesh on the topic
	IPColocationFactor float64     `json:"IPColocationFactor"`
	BehavioralPenalty  float64     `json:"behavioralPenalty"`
}

func (g GossipScores) Apply(rec *scoreRecord) {
	rec.PeerScores.Gossip = g
}

type PeerScores struct {
	Gossip      GossipScores `json:"gossip"`
	ReqRespSync float64      `json:"reqRespSync"`
}

// ScoreDatastore defines a type-safe API for getting and setting libp2p peer score information
type ScoreDatastore interface {
	// GetPeerScores returns the current scores for the specified peer
	GetPeerScores(id peer.ID) (PeerScores, error)

	// SetScore applies the given store diff to the specified peer
	SetScore(id peer.ID, diff ScoreDiff) error
}

// ScoreDiff defines a type-safe batch of changes to apply to the peer-scoring record of the peer.
// The scoreRecord the diff is applied to is private: diffs can only be defined in this package,
// to ensure changes to the record are non-breaking.
type ScoreDiff interface {
	Apply(score *scoreRecord)
}

// ExtendedPeerstore defines a type-safe API to work with additional peer metadata based on a libp2p peerstore.Peerstore
type ExtendedPeerstore interface {
	peerstore.Peerstore
	ScoreDatastore
	peerstore.CertifiedAddrBook
}
