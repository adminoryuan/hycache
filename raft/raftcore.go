package raft

import (
	"fmt"
	"math/rand"
	"net/rpc"
	"time"
)

//尝试变成候选者
//此时每个参与者都有机会
// 那么这个幸运儿是谁了?
func (r *Raft) ProCandidate() bool {

	time.Sleep(getRandomTime())
	time.Sleep(getRandomTime())
	r.mu.Lock()
	if r.VotedFor == -1 && r.leaderId == -1 {
		//还是保持不变的话 那么你有机会成为大哥

		fmt.Println("我成功成为候选人了")
		r.mu.Unlock()
		r.state = 1
		r.SetTerm(r.term + 1)
		r.SetVoteFor(r.node.RaftId)
		r.AddVoted()
		//r.SetCurrentLeader(r.node.RaftId)

		return true

	}
	r.mu.Unlock()
	return false

}
func getCurrSecoud() int64 {
	return time.Now().UnixNano() / int64(time.Second)
}

func (r *Raft) Start() {
	for {
		if r.ProCandidate() {
			if r.Election() {
				fmt.Println("当选")
				break
			} else {
				continue
			}
		} else {
			break
		}

	}
}

//开始选举
func (r *Raft) Election() bool {
	okChan := make(chan struct{})

	reqVoteobj := ReqVote{Term: r.term, CandidateId: r.leaderId, LastLogIndex: 0, LastLogTerm: 0}
	//var rs ReqVoteRes
	go r.ForWardCall("Raft.AskForNodeVote", reqVoteobj, new(ReqVoteRes), func(b interface{}) {
		if b == nil {
			return
		}
		result := b.(*ReqVoteRes)

		if result.VoteGranted {
			okChan <- struct{}{}
		}

	})
	for {
		select {
		case <-time.After(time.Second * time.Duration(r.EleTimeOut)):

			fmt.Println("选举超时")

			r.setDefault()

			return false
		case <-okChan:
			fmt.Println("选票+1")
			r.AddVoted()
			fmt.Println(r.Vote)
			if r.Vote > (3/2) && r.leaderId == -1 {
				//当选了

				r.state = 2
				r.SetCurrentLeader(r.node.RaftId)

				//广播给子节点 告诉他们我当选了
				//	var re bool
				r.ForWardCall("Raft.RecvLeaderTaskOffice", r.node, new(bool), func(b interface{}) {
					if b == nil {
						return
					}
					if *b.(*bool) {
						fmt.Println("通知成功")
					} else {
						fmt.Println("通知失败")
					}

				})
				go r.Heartbeat()
				return true

			}

			fmt.Println("未当选")
			return false
		}
	}
}

//心跳检测
func (r *Raft) DeteHeart() {
	for {
		time.Sleep(time.Microsecond * 4000)

		if r.lastApplied != 0 && (getCurrSecoud()-int64(r.lastApplied)) > int64(r.HeartSleep) {
			r.setDefault()
			r.lastApplied = 0
			fmt.Println("未收到心跳")

			r.mu.Lock()

			r.leaderId = -1

			r.VotedFor = -1

			r.state = -1

			r.Vote = 0

			r.lastApplied = 0

			r.mu.Unlock()

			go r.Start()

		}

	}
}

//向其它节点发送
func (r *Raft) Heartbeat() {
	for {
		select {
		case <-r.IsSendHeart:
			return
		default:
			r.ForWardCall("Raft.RecvHeart", r.node, new(ReqVoteRes), func(b interface{}) {
				//	fmt.Println("客户端收到心跳")
			})
			fmt.Println("发送心跳")
			time.Sleep(time.Duration(r.HeartSleep * int(time.Second)))

		}

	}
}

//生产一个随机时候
func getRandomTime() time.Duration {
	rand.Seed(time.Now().UnixNano())
	return time.Millisecond*time.Duration(rand.Intn(3000)) + 1500

}

//向所有节点追加日志
func (r *Raft) AppenEntry(log LogEntry, rely *ReqVoteRes) error {
	for _, node := range r.regisConfig.Globle {
		go r.CallRecvLogger(node)
	}
	return nil
}
func (r *Raft) CallRecvLogger(cf RaftNode) bool {

	cli, err := rpc.DialHTTP("tcp", "127.0.0.1"+cf.Port)
	if err != nil {
		return false
	}
	for {
		var re ReqVoteRes
		log := r.Log[r.nextIndex[cf.RaftId]]

		err := cli.Call("Raft.RecvLogger", log, &re)
		if err != nil {
			continue
		}
		if re.VoteGranted {
			r.nextIndex[cf.RaftId] -= 1
			continue
		} else {
			r.matchIndex[cf.RaftId] = r.nextIndex[cf.RaftId]
		}

	}

}
