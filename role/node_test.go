package role

import (
	"fmt"
	"log"
	"testing"
)

func TestNode_GetFileAndUpdateStatus(t *testing.T) {

	node := &Node{
		NodeIndex: 1,
		Host:      "39.96.23.156",
		Passwd:    "ZZ@123123",
	}
	node.Init()

	node.GetFileAndUpdateStatus()
	for _, miner := range node.Miners {
		log.Print(miner.name)
	}
	for _, storage := range node.Storages {
		log.Print(storage.name)
	}
}

func TestNode_CheckErrorLog(t *testing.T) {
	node := &Node{
		NodeIndex: 1,
		Host:      "39.96.23.156",
		Passwd:    "ZZ@123123",
	}
	node.Init()
	node.Miners = append(
		node.Miners,
		&Miner{
			nodeIndex:   node.NodeIndex,
			name:        fmt.Sprintf("%v:miner-%v", node.Host, 0),
			remoteDir:   "",
			localDir:    "/home/chengze/go/src/github.com/multivactech/monitor/data/39.96.23.156/miner-0",
			lastErrSize: 0,
			shardNum:    0,
		},
	)
	node.CheckErrorLog()
}
