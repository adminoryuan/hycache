package main

import (
	"SmiRaftDis/raft"
	"flag"
	"sync"
)

func main() {
	port := flag.String("port", "9000", "输入监听的端口")
	id := flag.Int("id", -1, "输入监听的端口")

	flag.Parse()

	w := sync.WaitGroup{}

	w.Add(1)

	raftobj := raft.NewRaft(raft.RaftNode{RaftId: *id, Port: *port})
	//raftobj.AddConfig(raft.RaftNode{RaftId: *id, Port: *port})

	go raft.ListenRpc(raftobj)

	go raftobj.DeteHeart()

	raftobj.Start()

	w.Wait()
}
