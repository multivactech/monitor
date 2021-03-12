package role

import (
	"log"
	"os"
	"unsafe"
)

type Miner struct {
	nodeIndex   int
	name        string
	remoteDir   string
	localDir    string
	lastErrSize int64
	shardNum    int
}

// check miner's error.log and if it have update,
// we return the content something new.
func (miner *Miner) CheckErrorLog(dir string) (string, error) {
	if file, err := os.Open(dir); err != nil {
		return "", err
	} else {
		defer file.Close()
		fileInfo, _ := file.Stat()
		newErrSize := fileInfo.Size()
		log.Printf("%v size is : %v, preSize is : %v", dir, newErrSize, miner.lastErrSize)
		if newErrSize-miner.lastErrSize < 10 {
			miner.lastErrSize = newErrSize
			return "", nil
		}
		context := make([]byte, newErrSize-miner.lastErrSize)
		file.ReadAt(context, miner.lastErrSize)
		miner.lastErrSize = newErrSize
		return toString(context), nil
	}
}

func (miner *Miner) GetShardNum() int {
	return miner.shardNum
}

// convert []byte to string
func toString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// getMinerLocalDir
func (miner *Miner) GetLocalDir() string {
	return miner.localDir
}
