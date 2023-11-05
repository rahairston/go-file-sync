package filesystem

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

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

	baseFileName := strings.TrimPrefix(fileName, dir.sourceFs.GetPath())
	dstPath := dir.dstFs.CorrectPathSeparator(dir.dstFs.GetPath() + baseFileName)
	dstFile, err := dir.dstFs.Open(dstPath)

	if err != nil {
		os.MkdirAll(filepath.Dir(dstPath), os.ModePerm)
		file, err := dir.dstFs.Create(dstPath)

		if err != nil {
			panic(err)
		}
		io.Copy(file, srcFile)
	} else {
		if !common.DeepCompare(srcFile, dstFile) {
			_, err := io.Copy(dstFile, srcFile)
			if err != nil {
				panic(err)
			}
		}
	}

	c <- fileName
}
