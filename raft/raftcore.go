package raft

import "time"

//尝试候选者
//此时每个参与者都有机会
// 那么这个幸运儿是谁了?
func (r *Raft) ProCandidate() {
	time.Sleep(getRandomTime())

	if r.VotedFor == -1 && r.currentTerm == -1 {
		//还是保持不变的话 那么你有机会成为大哥

	}

}

//生产一个随机时候
func getRandomTime() time.Duration {
	return 8000 * time.Microsecond
}

//向其它节点索引选票
func (r *Raft) AskForNodeVote() {

}
