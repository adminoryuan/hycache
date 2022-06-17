package raft

import "sync"

type Raft struct {
	mu sync.Mutex

	HeartSleep int //间隔

	EleTimeOut int

	currentTerm int //当前Term

	VotedFor int //我投票给了谁

	Vote int //累计票数

	state int //标识了我leard 、follower

	Log []LogEntry

	node RaftNode

	IsSendHeart chan struct{}

	CommitIndex int //目前确定已经提交的idx

	lastApplied int

	nextIndex map[int]int

	matchIndex map[int]int

	regisConfig Config
}

type RaftNode struct {
	RaftId int

	Port string
}

func NewRaft(node RaftNode) *Raft {
	r := new(Raft)

	r.node = node
	r.currentTerm = -1

	r.VotedFor = -1

	r.state = -1

	r.mu = sync.Mutex{}

	raftConfig := []RaftNode{RaftNode{RaftId: 1, Port: ":9000"}, RaftNode{RaftId: 2, Port: ":9001"}, RaftNode{RaftId: 3, Port: ":9002"}}
	r.regisConfig = Config{Globle: raftConfig}
	r.EleTimeOut = 10
	r.HeartSleep = 3

	return r

}

func (c *Raft) AddConfig(nodes RaftNode) {
	c.regisConfig.Globle = append(c.regisConfig.Globle, nodes)
}
func (rr *Raft) setDefault() {
	rr.mu.Lock()
	rr.state = -1

	rr.SetVoteFor(-1)

	rr.SetCurrentTerm(-1)
	rr.mu.Unlock()
}

//投票
func (r *Raft) AddVoted() {

	r.mu.Lock()
	r.Vote += 1
	r.mu.Unlock()
}

func (r *Raft) SetVoteFor(id int) {
	r.mu.Lock()
	r.VotedFor = id
	r.mu.Unlock()
}

func (r *Raft) SetCurrentTerm(term int) {
	r.mu.Lock()
	r.currentTerm = term
	r.mu.Unlock()
}

type LogEntry struct {
	Term    int
	Index   int
	Command interface{}
}
