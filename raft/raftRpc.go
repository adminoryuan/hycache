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
func (r *Raft) ForWardCall(method string, args interface{}, fun func(bool)) {
	for _, cf := range r.regisConfig.Globle {
		if cf.RaftId == r.node.RaftId {

			//不需call 自己
			continue
		}
		cli, err := rpc.DialHTTP("tcp", "127.0.0.1"+cf.Port)
		if err != nil {
			//	//log.Fatalf("链接出错")
			//	fmt.Println("链接出错")
			fun(false)
			continue
		}

		var rely bool
		err = cli.Call(method, args, &rely)

		if err != nil {
			fun(false)
			continue
		}
		cli.Close()
		fun(rely)

	}

}

//接收心跳Rpc 函数
func (r *Raft) RecvHeart(node RaftNode, re *bool) error {
	if r.currentTerm == -1 {
		r.SetCurrentTerm(node.RaftId)
	}
	r.lastApplied = getCurrSecoud()
	*re = true
	fmt.Println("收到心跳")
	return nil
}

//向主节点投票
func (r *Raft) AskForNodeVote(node RaftNode, rely *bool) error {

	//fmt.Println("调用")
	//r.mu.Lock()
	if r.currentTerm != -1 || r.VotedFor != -1 {
		*rely = false
		fmt.Println("投票失败")
		return nil
	}

	r.SetVoteFor(node.RaftId)
	*rely = true
	//r.mu.Unlock()
	return nil
}

//大哥上任了
func (r *Raft) RecvLeaderTaskOffice(node RaftNode, rely *bool) error {

	fmt.Printf("大哥是%d \n", node.RaftId)

	r.SetCurrentTerm(node.RaftId)

	r.setDefault()
	*rely = true
	return nil
}
