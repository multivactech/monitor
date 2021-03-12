package role

import (
	"log"
	"os"
)

type Storage struct {
	name        string
	remoteDir   string
	localDir    string
	lastErrSize int64
}

func (storage *Storage) CheckErrorLog(dir string) (string, error) {
	//log.Print(dir)
	if file, err := os.Open(dir); err != nil {
		return "", err
	} else {
		defer file.Close()
		fileInfo, _ := file.Stat()
		newErrSize := fileInfo.Size()
		log.Printf("%v size is : %v, preSize is : %v", dir, newErrSize, storage.lastErrSize)
		if newErrSize-storage.lastErrSize < 10 {
			storage.lastErrSize = newErrSize
			return "", nil
		}
		context := make([]byte, newErrSize-storage.lastErrSize)
		file.ReadAt(context, storage.lastErrSize)
		storage.lastErrSize = newErrSize
		return toString(context), nil
	}
}
