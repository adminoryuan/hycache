package raft

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"time"
)

//监听rpc
func ListenRpc(r *Raft) {
	err := rpc.Register(r)
	if err == nil {
		log.Fatal(err.Error())
	}

	rpc.HandleHTTP()

	lis, err := net.Listen("tcp", r.node.Port)
	if err != nil {
		panic(err)

	}

	fmt.Println("rpc 启动成功")
	go func() {
		for {
			con, err := lis.Accept()
			if err != nil {
				panic(err)
			}

			rpc.ServeConn(con)
		}
	}()

}

//转发rpc 调用
func (r *Raft) ForWardCall(method string, args interface{}, rely bool, fun func(bool)) {
	for _, cf := range r.regisConfig.Globle {
		if cf.RaftId == r.node.RaftId {
			fmt.Println("call 自己")
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
		err = cli.Call(method, args, rely)
		if err != nil {
			fun(false)
			continue
		}
		cli.Close()
		fun(true)

	}

}

//接收心跳Rpc 函数
func (r *Raft) RecvHeart(node RaftNode, re *bool) {
	r.SetCurrentTerm(node.RaftId)

	r.lastApplied = time.Now().Minute()
	*re = true

}

//向主节点投票
func (r *Raft) AskForNodeVote(node RaftNode) {

	fmt.Println("调用")
	r.mu.Lock()

	r.SetVoteFor(node.RaftId)
	r.mu.Unlock()
}

//大哥上任了
func (r *Raft) RecvLeaderTaskOffice(node RaftNode) {
	r.mu.Lock()
	r.SetCurrentTerm(node.RaftId)
	fmt.Printf("大哥是%d", node.RaftId)
	r.mu.Unlock()

}
