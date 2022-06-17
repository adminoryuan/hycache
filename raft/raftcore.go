package raft

import (
	"fmt"
	"time"
)

//尝试变成候选者
//此时每个参与者都有机会
// 那么这个幸运儿是谁了?
func (r *Raft) ProCandidate() bool {
	time.Sleep(getRandomTime())

	if r.VotedFor == -1 && r.currentTerm == -1 {
		//还是保持不变的话 那么你有机会成为大哥
		r.mu.Unlock()

		r.state = 1
		r.currentTerm = -1
		r.VotedFor = r.node.RaftId
		r.Vote += 1 // 为自己投票

		r.mu.Unlock()

		return true
	}
	return false
}
func getCurrSecoud() int64 {
	return time.Now().UnixNano() / int64(time.Second)
}

//开始选举
func (r *Raft) Election() bool {
	okChan := make(chan struct{})

	go r.ForWardCall("Raft.AskForNodeVote", r.node, true, func(b bool) {
		fmt.Println("请候选人投票完毕")
		okChan <- struct{}{}
	})
	for {
		select {
		case <-time.After(time.Second * time.Duration(r.EleTimeOut)):
			return false

		case <-okChan:
			if r.Vote > len(r.regisConfig.Globle)/2 && r.currentTerm == -1 {
				//当选了

				r.state = 2
				r.SetCurrentTerm(r.node.RaftId)

				//广播给子节点 告诉他们我当选了
				r.ForWardCall("Raft.RecvLeaderTaskOffice", r.node, true, func(b bool) {
					fmt.Println("成功")
				})
				go r.Heartbeat()
				return true

			}
		}
	}
}

//心跳检测
func (r *Raft) DeteHeart() {
	for {
		if r.lastApplied != 0 && getCurrSecoud()-int64(r.lastApplied) > int64(r.HeartSleep) {
			r.setDefault()
			r.lastApplied = 0
			func() {
				for {
					if r.ProCandidate() {
						if r.Election() {
							break
						}
					} else {
						break
					}
				}
			}()
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
			r.ForWardCall("Raft.RecvHeart", r.node, false, func(b bool) {
				fmt.Println("客户端收到心跳")
			})
			time.Sleep(time.Duration(r.HeartSleep))

		}

	}
}

//生产一个随机时候
func getRandomTime() time.Duration {
	return 1234
}
