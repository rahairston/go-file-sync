package filesystem

import (
	"io"
	"log"

	"github.com/rahairston/go-file-sync/common"
)

type DirClient struct {
	sourceFs   common.FileSystem
	dstFs      common.FileSystem
	exclusions common.ExcludeObject
}

func BuildDirClient(syncConfig *common.SyncObject, sourceFs common.FileSystem, dstFs common.FileSystem) (*DirClient, error) {
	return &DirClient{
		sourceFs:   sourceFs,
		dstFs:      dstFs,
		exclusions: syncConfig.Exclusions,
	}, nil
}

func (dir DirClient) SyncFiles() {
	defer dir.sourceFs.Close()
	defer dir.dstFs.Close()

	fileNames := dir.sourceFs.GetFileNames(dir.sourceFs.GetPath(), dir.exclusions)

	c := make(chan string, len(fileNames))

	for _, fileName := range fileNames {
		go dir.SyncFile(fileName, c)
	}

	for i := 0; i < cap(c); i++ {
		log.Println(<-c)
	}
}

func (dir DirClient) SyncFile(fileName string, c chan string) {

	srcFile, err := dir.sourceFs.Open(fileName)

	if err != nil {
		panic(err)
	}

	dstPath := dir.dstFs.CorrectPathSeparator(dir.dstFs.GetPath() + fileName)

	dstFile, err := dir.dstFs.Open(fileName)

	if err != nil {
		panic(err)
	}

	if err != nil {
		file, err := dir.dstFs.Create(dstPath)

		if err != nil {
			panic(err)
		}
		io.Copy(file, srcFile)
	} else {
		if !common.DeepCompare(srcFile, dstFile) {
			io.Copy(dstFile, srcFile)
		}
	}

	c <- fileName
}
