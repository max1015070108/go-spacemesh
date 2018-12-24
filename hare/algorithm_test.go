package hare

import (
	"github.com/spacemeshos/go-spacemesh/crypto"
	"github.com/spacemeshos/go-spacemesh/hare/config"
	"github.com/spacemeshos/go-spacemesh/hare/pb"
	"github.com/spacemeshos/go-spacemesh/p2p/service"
	"github.com/stretchr/testify/assert"
	"testing"
)

var cfg = config.DefaultConfig()

func generatePubKey(t *testing.T) crypto.PublicKey {
	_, pub, err := crypto.GenerateKeyPair()

	if err != nil {
		assert.Fail(t, "failed generating key")
		t.FailNow()
	}

	return pub
}

// test that a message to a specific set id is delivered by the broker
func TestConsensusProcess_StartTwice(t *testing.T) {
	sim := service.NewSimulator()
	n1 := sim.NewNode()

	broker := NewBroker(n1)
	s := NewEmptySet(cfg.SetSize)
	oracle := NewMockOracle()
	signing := NewMockSigning()

	proc := NewConsensusProcess(cfg, generatePubKey(t), *setId1, *s, oracle, signing, n1)
	broker.Register(setId1, proc)
	err := proc.Start()
	assert.Equal(t, nil, err)
	err = proc.Start()
	assert.Equal(t, "instance already started", err.Error())
}

func TestConsensusProcess_eventLoop(t *testing.T) {
	sim := service.NewSimulator()
	n1 := sim.NewNode()
	n2 := sim.NewNode()

	broker := NewBroker(n1)
	s := NewEmptySet(cfg.SetSize)
	oracle := NewMockOracle()
	signing := NewMockSigning()

	proc := NewConsensusProcess(cfg, generatePubKey(t), *setId1, *s, oracle, signing, n1)
	broker.Register(setId1, proc)
	go proc.eventLoop()
	n2.Broadcast(ProtoName, []byte{})

	proc.Close()
	<-proc.CloseChannel()
}

func TestConsensusProcess_handleMessage(t *testing.T) {
	sim := service.NewSimulator()
	n1 := sim.NewNode()

	broker := NewBroker(n1)
	s := NewEmptySet(cfg.SetSize)
	oracle := NewMockOracle()
	signing := NewMockSigning()

	proc := NewConsensusProcess(cfg, generatePubKey(t), *setId1, *s, oracle, signing, n1)
	broker.Register(setId1, proc)

	m := NewMessageBuilder().SetIteration(0).SetSetId(*setId1).SetPubKey(generatePubKey(t)).Sign(proc.signing).Build()

	proc.handleMessage(m)
}

func TestConsensusProcess_nextRound(t *testing.T) {
	sim := service.NewSimulator()
	n1 := sim.NewNode()

	broker := NewBroker(n1)
	s := NewEmptySet(cfg.SetSize)
	oracle := NewMockOracle()
	signing := NewMockSigning()

	proc := NewConsensusProcess(cfg, generatePubKey(t), *setId1, *s, oracle, signing, n1)
	broker.Register(setId1, proc)

	proc.nextRound()
	assert.Equal(t, uint32(1), proc.k)
	proc.nextRound()
	assert.Equal(t, uint32(2), proc.k)
}

func generateConsensusProcess(t *testing.T) *ConsensusProcess {
	sim := service.NewSimulator()
	n1 := sim.NewNode()

	s := NewEmptySet(cfg.SetSize)
	oracle := NewMockOracle()
	signing := NewMockSigning()

	return NewConsensusProcess(cfg, generatePubKey(t), *setId1, *s, oracle, signing, n1)
}

func TestConsensusProcess_DoesMatchRound(t *testing.T) {
	s := NewEmptySet(cfg.SetSize)
	pub := generatePubKey(t)
	cp := generateConsensusProcess(t)

	msgs := make([]*pb.HareMessage, 5, 5)
	msgs[0] = BuildPreRoundMsg(pub, s)
	msgs[1] = BuildStatusMsg(pub, s)
	msgs[2] = BuildProposalMsg(pub, s)
	msgs[3] = BuildCommitMsg(pub, s)
	msgs[4] = BuildNotifyMsg(pub, s)

	rounds := make([][4]bool, 5) // index=round
	rounds[0] = [4]bool{true, true, true, true}
	rounds[1] = [4]bool{true, false, false, false}
	rounds[2] = [4]bool{false, true, true, false}
	rounds[3] = [4]bool{false, false, true, false}
	rounds[4] = [4]bool{true, true, true, true}

	for j:=0;j<len(msgs);j++ {
		for i := 0; i < 4; i++ {
			assert.Equal(t, rounds[j][i], cp.doesMessageMatchRound(msgs[j]))
			cp.nextRound()
		}
	}
}
