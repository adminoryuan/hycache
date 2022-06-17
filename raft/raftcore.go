package raft

import (
	"fmt"
	"math/rand"
	"time"
)

//尝试变成候选者
//此时每个参与者都有机会
// 那么这个幸运儿是谁了?
func (r *Raft) ProCandidate() bool {
	time.Sleep(getRandomTime())

	if r.VotedFor == -1 && r.currentTerm == -1 {
		//还是保持不变的话 那么你有机会成为大哥
		fmt.Println("我成功成为候选人了")
		r.mu.Lock()

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

func (r *Raft) Start() {
	for r.currentTerm == -1 {
		if r.ProCandidate() {
			if r.Election() {
				fmt.Println("当选")
				//go r.DeteHeart()

			}
		}
	}
}

//开始选举
func (r *Raft) Election() bool {
	okChan := make(chan struct{})

	fmt.Println("开始选举")
	go r.ForWardCall("Raft.AskForNodeVote", r.node, true, func(b bool) {
		fmt.Println(b)
		if b {
			okChan <- struct{}{}
		}

	})
	for {
		select {
		case <-time.After(time.Second * time.Duration(r.EleTimeOut)):

			fmt.Println("选举超时")

			return false

		case <-okChan:
			r.mu.Lock()
			r.Vote += 1
			r.mu.Unlock()
			fmt.Println(r.Vote)
			if r.Vote > (len(r.regisConfig.Globle)/2) && r.currentTerm == -1 {
				//当选了

				r.state = 2
				r.SetCurrentTerm(r.node.RaftId)

				//广播给子节点 告诉他们我当选了
				r.ForWardCall("Raft.RecvLeaderTaskOffice", r.node, true, func(b bool) {
					fmt.Println("通知成功")
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
		if r.lastApplied != 0 && getCurrSecoud()-int64(r.lastApplied) > int64(r.HeartSleep) {
			r.setDefault()
			r.lastApplied = 0
			fmt.Println("未收到心跳")
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
			r.ForWardCall("Raft.RecvHeart", r.node, false, func(b bool) {
				//	fmt.Println("客户端收到心跳")
			})
			time.Sleep(time.Duration(r.HeartSleep))

		}

	}
}

//生产一个随机时候
func getRandomTime() time.Duration {
	return time.Duration(3000 + rand.Intn(10000)*int(time.Microsecond))
}
