package role

import (
	"fmt"
	"testing"
)

func TestShard_getMinerShardInfo(t *testing.T) {
	shard := &Shard{
		LastCheckSimpleBlockInfo: nil,
		Miners: map[string]*Miner{"/home/chengze/go/src/monitor/data/%d/39.96.64.104/miner-0": &Miner{}},
	}
	shardInfo, err := shard.getMinerShardInfo("/home/chengze/go/src/monitor/data/0/39.96.64.104/miner-0/data/1")
	if err != nil {
		t.Log(err)
	} else {
		content := shardInfo.content
		for _, line := range content {
			fmt.Println(line)
		}
	}
}

func TestShard_getSimpleBlockInfo(t *testing.T) {
	shard := &Shard{}
	line := "{round 0, leader 17d6aa0f6ae26d48d9f48d4a72c2d6956b8373db837cf582a7e3365c8d8d9511, prevHash 10aeeed0b2f266d99dd49c8b29db7039c6d3e407fb8fdd5687f793dae6575265, curHash 15246a48110cf943a899f7fd21c0044bf3a3b257ee2c204a223a032b0336fbca, txCount 1, source 0"
	simpleBlockInfo := shard.getSimpleBlockInfo(line)
	t.Log(simpleBlockInfo)
}
