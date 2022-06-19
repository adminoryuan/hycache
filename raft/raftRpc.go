package raft

import (
	"fmt"
	"log"
	"net/http"
	"net/rpc"
)

//?????
//监听rpc
func ListenRpc(Raft *Raft) {
	err := rpc.Register(Raft)
	if err != nil {
		log.Fatal(err.Error())
	}

	rpc.HandleHTTP()

	err = http.ListenAndServe(Raft.node.Port, nil)
	if err != nil {
		fmt.Println("rpc 注入失败")
	}
	fmt.Println("rpc 服务启动成功")

}

//转发rpc 调用
func (r *Raft) ForWardCall(method string, args interface{}, res interface{}, fun func(interface{})) {
	for _, cf := range r.regisConfig.Globle {
		if cf.RaftId == r.node.RaftId {

			//不需call 自己
			continue
		}
		cli, err := rpc.DialHTTP("tcp", "127.0.0.1"+cf.Port)
		if err != nil {
			//	//log.Fatalf("链接出错")

			fun(nil)
			continue
		}

		err = cli.Call(method, args, res)
		if err != nil {
			fmt.Println(err.Error())
			fmt.Println("call 失败")
			fun(nil)
			continue
		}

		cli.Close()

		fun(res)

	}

}

//接收心跳Rpc 函数
func (r *Raft) RecvHeart(node RaftNode, re *ReqVoteRes) error {
	re.Term = node.term
	if r.term < node.term {
		re.VoteGranted = false
	} else if r.leaderId == -1 {
		r.SetCurrentLeader(node.RaftId)
	} else {
		r.lastApplied = getCurrSecoud()

		re.VoteGranted = true
	}
	fmt.Println("收到心跳")
	return nil
}

//向主节点投票
func (r *Raft) AskForNodeVote(node ReqVote, vote *ReqVoteRes) error {

	if node.Term < r.term {
		*vote = ReqVoteRes{
			VoteGranted: false,
			Term:        r.term,
		}

		return nil
	}

	fmt.Printf("%d -% d", r.leaderId, r.VotedFor)
	if r.leaderId != -1 || r.VotedFor != -1 {
		//vote.voteGranted = false
		//vote.term = r.term
		*vote = ReqVoteRes{
			VoteGranted: false,
			Term:        r.term,
		}

		fmt.Println("投票失败")
		return nil
	}

	r.SetVoteFor(node.CandidateId)
	*vote = ReqVoteRes{
		VoteGranted: true,
		Term:        r.term,
	}

	return nil
}

//大哥上任了
func (r *Raft) RecvLeaderTaskOffice(node RaftNode, rely *bool) error {

	fmt.Printf("大哥是%d \n", node.RaftId)

	r.SetCurrentLeader(node.RaftId)

	r.setDefault()
	*rely = true
	return nil
}
func (r *Raft) RecvLogger(lrty LogEntry, res ReqVoteRes, rely *bool) {

}
