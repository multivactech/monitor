package role

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/multivactech/monitor/errorInfo"
)

// store per shard info of every miner
type Shard struct {
	ShardIndex               int
	LastCheckSimpleBlockInfo *SimpleBlockInfo
	Miners                   map[string]*Miner
}

type MinerShardLogInfo struct {
	shardIndex int
	minerInfo  string
	miner      *Miner
	content    []*SimpleBlockInfo
}

type SimpleBlockInfo struct {
	round    int
	leader   string
	prevHash string
	curHash  string
	txCount  int
	source   int
}

// convert line string to simpleBlcokInfo，
// then compare them all.
func (shard *Shard) CheckAll() {
	var shardLog []*MinerShardLogInfo

	for localMinerDir, miner := range shard.Miners {
		minerInShardInfo, err := shard.getMinerShardInfo(fmt.Sprintf("%v/shardInfo/%v", localMinerDir, shard.ShardIndex))
		minerInShardInfo.miner = miner
		// log.Printf("localDir:%v  miner: %v, minerInShardInfo: %v", localDir, miner, minerInShardInfo)
		if err != nil {
			log.Printf("can't get minerShardInfo of %v. %v", localMinerDir, err)
			continue
		}
		log.Printf("checkMinerShardInfo, miner: %v, shardIndex: %v", minerInShardInfo.miner.name, shard.ShardIndex)
		shard.checkMinerShardInfo(minerInShardInfo)
		// log.Printf("checkMinerShardInfo complete, miner: %v, shardIndex: %v", minerInShardInfo.miner.name, shard.ShardIndex)

		shardLog = append(shardLog, minerInShardInfo)
	}
	if len(shardLog) > 0 {
		shard.checkMinersLog(shardLog)
	} else {
		time.Sleep(5 * time.Second)
	}
}

// judge prevHash == preRound.curRound
func (shard *Shard) checkMinerShardInfo(simpleBlockInfos *MinerShardLogInfo) {
	if len(simpleBlockInfos.content) == 0 {
		log.Printf("shard %v, miner %v not update", simpleBlockInfos.shardIndex, simpleBlockInfos.miner.name)
		//TODO: miner die and not update 5 times, we will send email.
		return
	}
	prevHash := ""

	if shard.LastCheckSimpleBlockInfo == nil {
		prevHash = simpleBlockInfos.content[0].prevHash
	} else {
		prevHash = shard.LastCheckSimpleBlockInfo.curHash
	}

	for _, blockInfo := range simpleBlockInfos.content {
		// log.Print(prevHash, " ", blockInfo.prevHash, " ", strings.Compare(prevHash, blockInfo.prevHash))
		if strings.Compare(prevHash, blockInfo.prevHash) != 0 {
			log.Printf("矿工%s,分片%d,round为%v时preHash与round为%v时curHash不相同, %v, %v", simpleBlockInfos.miner.name, simpleBlockInfos.shardIndex, blockInfo.round-1, blockInfo.round, prevHash, blockInfo.prevHash)
			errorInfo.AllError <- &errorInfo.ErrorInfo{
				Host:         "",
				ErrorContent: fmt.Sprintf("矿工%s,分片%d,round为%d时preHash与round为%d时curHash不相同, %v, %v", simpleBlockInfos.miner.name, simpleBlockInfos.shardIndex, blockInfo.round-1, blockInfo.round, prevHash, blockInfo.prevHash),
				Node:         simpleBlockInfos.miner.name,
			}
			break
		}
		prevHash = blockInfo.curHash
	}
}

// judge every miner in current round info should equal
func (shard *Shard) checkMinersLog(shardLog []*MinerShardLogInfo) {
	log.Printf("now start checkMinersLog, shard %v", shard.ShardIndex)

	for i := 0; i < len(shardLog[0].content); i++ {
		target := shardLog[0].content[i]
		// log.Printf("round: %v", target.round)
		for _, minerShardInfo := range shardLog {
			if len(minerShardInfo.content) == 0 {
				errorInfo.AllError <- &errorInfo.ErrorInfo{
					Host:         minerShardInfo.minerInfo,
					ErrorContent: fmt.Sprintf("矿工%v在分片%v未更新信息", minerShardInfo.minerInfo, shard.ShardIndex),
				}
				continue
			}
			if i >= len(minerShardInfo.content) {
				shard.LastCheckSimpleBlockInfo = minerShardInfo.content[i-1]
				return
			}
			if !target.equal(minerShardInfo.content[i]) {
				errorInfo.AllError <- &errorInfo.ErrorInfo{
					Host:         minerShardInfo.minerInfo,
					ErrorContent: fmt.Sprintf("%v矿工和%v矿工在%v分片%v轮信息不同。", shardLog[0].miner.name, minerShardInfo.miner.name, shard.ShardIndex, target.round),
				}
			}
		}
	}
	// log.Print(len(shardLog[0].content))
	if len(shardLog) > 0 && len(shardLog[0].content) > 0 {
		shard.LastCheckSimpleBlockInfo = shardLog[0].content[len(shardLog[0].content)-1]
	}
	log.Printf("shardLastCheckSimpleBlockInfo: %v", shard.LastCheckSimpleBlockInfo)
}

// according to logDir, read miner's shard info,
// then according to LastCheckRound,
// collect information in MinerShardLogInfo
func (shard *Shard) getMinerShardInfo(logDir string) (*MinerShardLogInfo, error) {
	minerShardLogInfo := &MinerShardLogInfo{
		minerInfo:  logDir,
		shardIndex: shard.ShardIndex,
	}
	logFile, err := os.Open(logDir)
	if err != nil {
		log.Printf("read logFile %v, err: %v", logDir, err)
		return minerShardLogInfo, err
	}
	defer logFile.Close()

	rd := bufio.NewReader(logFile)
	for index := 0; ; index++ {
		line, err := rd.ReadString('\n')
		if err != nil || io.EOF == err {
			return minerShardLogInfo, nil
		}
		if shard.LastCheckSimpleBlockInfo == nil {
			simpleBlockInfo := shard.getSimpleBlockInfo(line)
			if simpleBlockInfo != nil {
				log.Print("simpleBlockInfo have not complete")
			} else {
				minerShardLogInfo.content = append(minerShardLogInfo.content, simpleBlockInfo)
			}
		} else {
			if index < shard.LastCheckSimpleBlockInfo.round {
				continue
			}
			simpleBlockInfo := shard.getSimpleBlockInfo(line)
			if simpleBlockInfo != nil {
				log.Print("simpleBlockInfo have not complete")
			} else {
				if simpleBlockInfo.round <= shard.LastCheckSimpleBlockInfo.round {
					continue
				}
				minerShardLogInfo.content = append(minerShardLogInfo.content, simpleBlockInfo)
			}
		}
	}
}

// convent line to SimpleBlockInfo
func (shard *Shard) getSimpleBlockInfo(line string) *SimpleBlockInfo {
	ss := strings.Split(line, ",")
	if len(ss) < 5 {
		return nil
	}
	ssplitToRound := strings.TrimPrefix(ss[0], "{round")
	ssplitToRound = strings.TrimSpace(ssplitToRound)
	round, _ := strconv.Atoi(ssplitToRound)

	ssplitToLeader := strings.TrimPrefix(ss[1], " leader")
	ssplitToLeader = strings.TrimSpace(ssplitToLeader)

	ssplitToPrevHash := strings.TrimPrefix(ss[2], " prevHash")
	ssplitToPrevHash = strings.TrimSpace(ssplitToPrevHash)

	ssplitToCurHash := strings.TrimPrefix(ss[3], " curHash")
	ssplitToCurHash = strings.TrimSpace(ssplitToCurHash)

	ssplitToTxCount := strings.TrimPrefix(ss[4], " txCount")
	ssplitToTxCount = strings.TrimSpace(ssplitToTxCount)
	txCount, _ := strconv.Atoi(ssplitToTxCount)

	ssplitToSource := strings.TrimPrefix(ss[5], " source")
	ssplitToSource = strings.TrimSpace(ssplitToSource)
	source, _ := strconv.Atoi(ssplitToSource)

	return &SimpleBlockInfo{
		round:    round,
		leader:   ssplitToLeader,
		prevHash: ssplitToPrevHash,
		curHash:  ssplitToCurHash,
		txCount:  txCount,
		source:   source,
	}
}

func (sbi1 SimpleBlockInfo) equal(sbi2 *SimpleBlockInfo) bool {
	// ignore source , because maybe source is sync
	if sbi1.round == sbi2.round {
		if sbi1.prevHash == sbi2.prevHash {
			if sbi1.curHash == sbi2.curHash {
				if sbi1.txCount == sbi2.txCount {
					return true
				}
			}
		}
	}

	return false
}
