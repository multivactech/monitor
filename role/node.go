package role

import (
	"fmt"
	"log"
	"strconv"

	"github.com/multivactech/monitor/config"

	"github.com/multivactech/monitor/connect"
	"github.com/multivactech/monitor/errorInfo"
)

// A Node means a remote ssh computer,
// it maybe have more than one miner or storage...
type Node struct {
	Host      string
	Passwd    string
	conn      *connect.Connection
	Miners    []*Miner
	Storages  []*Storage
	NodeIndex int
}

func (node *Node) Init() {
	node.conn = &connect.Connection{}
	if err := node.conn.Init(node.Host, node.Passwd); err != nil {
		panic(err)
	}
}

// check whether node' s miner and storage have output error.log and update it.
func (node *Node) CheckErrorLog() {
	//remoteDir := config.Config.RemoteDir
	//localDir:= "/home/chengze/go/src/github.com/multivactech/monitor/data"

	log.Print("checkErrorLog")
	for _, miner := range node.Miners {
		localLogDir := fmt.Sprintf("%v/error.log", miner.localDir)
		errorContent, err := miner.CheckErrorLog(localLogDir)
		log.Printf("%v check localLogDir: %v, errorContent: %v", miner.name, localLogDir, errorContent)
		if err != nil || len(errorContent) < 10 {
			continue
		}
		// log.Printf("%v have error", miner.name)
		errorInfo.FindErrorFromString(errorContent)
	}

	for _, storage := range node.Storages {
		localLogDir := fmt.Sprintf("%v/error.log", storage.localDir)
		errorContent, err := storage.CheckErrorLog(localLogDir)
		log.Printf("%v check localLogDir: %v, errorContent: %v", storage.name, localLogDir, errorContent)
		if err != nil || len(errorContent) < 10 {
			continue
		}
		// log.Printf("%v have error", storage.name)
		errorInfo.FindErrorFromString(errorContent)
	}
}

// transver all miners and storage to check if it have been collected to shard.Miners and shard.Storages
// download all error.log and minerShardLog(data/simnet/0~15)
func (node *Node) GetFileAndUpdateStatus() {
	// traversal all file folder，miner-0 ~ 63, storage-0 ~ 63
	localDir := config.Config.LocalDir
	remoteDir := config.Config.RemoteDir
	//remoteDir := "/root/go/src/github.com/multivactech/MultiVAC/test"
	//localDir := "/home/chengze/go/src/github.com/multivactech/monitor/data"

	for index := 0; index < 64; index++ {
		// traversal all miner
		remoteMinerDir := fmt.Sprintf("%v/miner-%v", remoteDir, index)
		remoteStorageDir := fmt.Sprintf("%v/storage-%v", remoteDir, index)

		remoteMinerErrorLogDir := fmt.Sprintf("%v/miner-%v/logs/simnet/error.log", remoteDir, index)
		remoteStorageErrorLogDir := fmt.Sprintf("%v/storage-%v/logs/simnet/error.log", remoteDir, index)

		localMinerDir := fmt.Sprintf("%v/%v/miner-%v", localDir, node.Host, index)
		localStorageDir := fmt.Sprintf("%v/%v/storage-%v", localDir, node.Host, index)

		if node.conn.IsExist(remoteMinerDir) {
			log.Printf("%v:miner-%v exist", node.Host, index)
			if index < len(node.Miners) {
				// TODO: trace miner status
			} else {
				// add miner
				node.Miners = append(
					node.Miners,
					&Miner{
						nodeIndex:   node.NodeIndex,
						name:        fmt.Sprintf("%v:miner-%v", node.Host, index),
						remoteDir:   remoteMinerDir,
						localDir:    localMinerDir,
						lastErrSize: 0,
						shardNum:    0,
					},
				)
			}

			if node.conn.IsExist(remoteMinerErrorLogDir) {
				node.conn.GetFile(remoteMinerErrorLogDir, localMinerDir, "error.log")
			}

			for shardIndex := 0; shardIndex < 64; shardIndex++ {
				remoteFileDir := fmt.Sprintf("%v/data/simnet/%v", remoteMinerDir, shardIndex)
				if node.conn.IsExist(remoteFileDir) {
					node.conn.GetFile(remoteFileDir, localMinerDir+"/shardInfo", strconv.Itoa(shardIndex))
				} else {
					node.Miners[index].shardNum = shardIndex
					break
				}
			}
		}

		// traversal all storage
		if node.conn.IsExist(remoteStorageDir) {
			log.Printf("%v:storage-%v exist", node.Host, index)
			if index < len(node.Storages) {
				// TODO: if trace storage status?
			} else {
				// add miner number
				node.Storages = append(
					node.Storages,
					&Storage{
						name:        fmt.Sprintf("%v:storage-%v", node.Host, index),
						remoteDir:   remoteStorageDir,
						localDir:    localStorageDir,
						lastErrSize: 0,
					},
				)
			}
			//log.Print(remoteErrlogDir, "   ", localStorageDir)
			if node.conn.IsExist(remoteStorageErrorLogDir) {
				//log.Print(localStorageDir)
				node.conn.GetFile(remoteStorageErrorLogDir, localStorageDir, "error.log")
			}
		}
	}
}
