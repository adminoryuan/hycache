package raft

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
)

//监听rpc
func ListenRpc(r *Raft) {
	err := rpc.Register(r)
	if err == nil {
		log.Fatal(err.Error())
	}

	rpc.HandleHTTP()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", r.node.Port))
	if err != nil {
		panic(err)
	}
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

func Call(method string, args interface{}, rely bool) {


}
