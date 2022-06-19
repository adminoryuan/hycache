package raft

import (
	"sync"
)

type Raft struct {
	mu sync.Mutex

	HeartSleep int //间隔

	EleTimeOut int

	leaderId int //当前领导人id

	term int // 任期号

	VotedFor int //我投票给了谁

	Vote int //累计票数

	state int //标识了我leard 、follower

	Log []LogEntry

	node RaftNode

	IsSendHeart chan struct{}

	CommitIndex int //目前确定已经提交的idx

	lastApplied int64

	nextIndex map[int]int //每个节点下一个发送的日志index

	matchIndex map[int]int //节点已经提交日志的index

	regisConfig Config
}

type RaftNode struct {
	RaftId int

	Port string

	term int
}

type ReqVote struct {
	Term int

	CandidateId int

	LastLogIndex int

	LastLogTerm int
}

type ReqVoteRes struct {
	Term int

	Currid int

	VoteGranted bool
}

func NewRaft(node RaftNode) *Raft {
	r := new(Raft)

	r.node = node
	r.leaderId = -1

	r.VotedFor = -1

	r.state = -1

	r.nextIndex = make(map[int]int, 0)
	r.Vote = 0

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

	rr.state = 0

	rr.SetVoteFor(-1)

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

func (r *Raft) SetTerm(trem int) {
	r.mu.Lock()

	r.term = trem

	r.mu.Unlock()

}
func (r *Raft) SetCurrentLeader(term int) {
	r.mu.Lock()
	r.leaderId = term
	r.mu.Unlock()
}

type LogEntry struct {
	Term    int
	Index   int
	Command interface{}
}
