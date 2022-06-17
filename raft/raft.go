package raft

import "sync"

type Raft struct {
	mu sync.Mutex

	me int //idx

	currentTerm int //当前Term

	VotedFor int //我投票给了谁

	Vote int //票数

	state int //标识了我leard 、follower

	Log []LogEntry

	node RaftNode

	CommitIndex int //目前确定已经提交的idx

	lastApplied int

	nextIndex map[int]int

	matchIndex map[int]int
}

type RaftNode struct {
	RaftId int

	Port int
}

func NewRaft(node RaftNode) *Raft {
	r := new(Raft)

	r.node = node
	r.currentTerm = -1

	r.VotedFor = -1

	r.state = -1

	r.mu = sync.Mutex{}

	return r

}

//投票
func (r *Raft) AddVoted() {

	r.mu.Lock()
	r.Vote += 1
	r.mu.Unlock()
}

func (r *Raft) setVoteFor(id int) {
	r.mu.Lock()
	r.VotedFor = id
	r.mu.Unlock()
}
func (r *Raft) setCurrentTerm(term int) {
	r.mu.Lock()
	r.currentTerm = term
	r.mu.Unlock()
}

type LogEntry struct {
	Term    int
	Index   int
	Command interface{}
}
