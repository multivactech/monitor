package monitor

import (
	"log"
	"time"

	"github.com/multivactech/monitor/config"
	"github.com/multivactech/monitor/errorInfo"
	"github.com/multivactech/monitor/mail"
	"github.com/multivactech/monitor/role"
)

type Monitor struct {
	mail     mail.Mail
	nodePool []*role.Node
	shards   []*role.Shard
}

func (m *Monitor) init() {

	m.mail.Init()

	errorInfo.AllError = make(chan *errorInfo.ErrorInfo, 1000)

	for nodeIndex, sshNode := range config.Config.SshNodes {
		node := &role.Node{
			NodeIndex: nodeIndex,
			Host:      sshNode.Host,
			Passwd:    sshNode.Passwd,
		}
		m.nodePool = append(m.nodePool, node)
		node.Init()
	}
}

func (monitor *Monitor) Start() {
	monitor.init()
	log.Print("monitor init success.")

	go monitor.recvError()

	for checkIndex := 1; ; checkIndex++ {

		log.Printf("%v CheckRound start", checkIndex)

		monitor.startCheck()

		log.Printf("%v CheckRound end, then wait for 2 minutes", checkIndex)
		time.Sleep(2 * time.Minute)
	}
}

func (monitor *Monitor) startCheck() {
	for nodeIndex, node := range monitor.nodePool {
		node.GetFileAndUpdateStatus()
		node.CheckErrorLog()
		monitor.checkShardStatus(nodeIndex)
	}
	log.Print("check error.log compelet.")
	for _, shard := range monitor.shards {
		log.Print("now check shard ", shard.ShardIndex)
		shard.CheckAll()
	}
}

// listen to channel: AllError, if recv msg, package it,
//  and after checkDuration, if have packaged msg, send mail to us
func (monitor *Monitor) recvError() {
	checkDuration := 10 * time.Minute
	errContent := ""
	//time.Sleep(30 * time.Second)
	for {
		select {
		case error := <-errorInfo.AllError:
			errContent += error.Host + ": " + ": " + error.ErrorContent + "\r\n"
		default:
			if errContent != "" {
				monitor.mail.Send("测试网报警", errContent)
				log.Print("send email success!")

				// log.Print("测试网报警: ", errContent)
				errContent = ""
			}
			time.Sleep(checkDuration)
		}
	}
}

// check every miners' shard info, then collect the miner to corresponding shard
// TODO: test for this function
func (monitor *Monitor) checkShardStatus(nodeIndex int) {
	node := monitor.nodePool[nodeIndex]
	for _, miner := range node.Miners {
		for shardIndex := 0; shardIndex < miner.GetShardNum(); shardIndex++ {
			if shardIndex >= len(monitor.shards) {

				monitor.shards = append(monitor.shards, &role.Shard{
					ShardIndex: shardIndex,
					Miners:     map[string]*role.Miner{miner.GetLocalDir(): miner},
					LastCheckSimpleBlockInfo: nil,
				})
			}
			if _, ok := monitor.shards[shardIndex].Miners[miner.GetLocalDir()]; !ok {
				monitor.shards[shardIndex].Miners[miner.GetLocalDir()] = miner
			}
		}
	}
}
