package filesystem

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

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

func (dir DirClient) SyncFiles(lastFileMod map[string]string) map[string]string {
	defer dir.sourceFs.Close()
	defer dir.dstFs.Close()

	fileMod := make(map[string]string)

	fileNames := dir.sourceFs.GetFileNames(dir.sourceFs.GetPath(), dir.exclusions)

	c := make(chan struct {
		string
		time.Time
	}, len(fileNames))

	for _, fileName := range fileNames {
		go dir.SyncFile(fileName, lastFileMod[fileName], c)
	}

	for i := 0; i < cap(c); i++ {
		file := <-c
		log.Println(file.string)
		fileMod[file.string] = file.Time.String()
	}

	return fileMod
}

func (dir DirClient) SyncFile(fileName string, lastModifiedString string, c chan struct {
	string
	time.Time
}) {
	srcFile, err := dir.sourceFs.Open(fileName)
	srcInfo, _ := srcFile.Stat()
	srcMod := srcInfo.ModTime()

	if err != nil {
		panic(err)
	}

	lastModifiedDt, err := time.Parse("2006-01-02 15:04", lastModifiedString)
	if err != nil {
		lastModifiedDt = srcMod // SUB
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
	} else if srcMod.After(lastModifiedDt) && !common.DeepCompare(srcFile, dstFile) {
		_, err := io.Copy(dstFile, srcFile)
		if err != nil {
			panic(err)
		}
	}

	c <- struct {
		string
		time.Time
	}{fileName, srcMod}
}
