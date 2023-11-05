package filesystem

import (
	"errors"
	"io/fs"
	"log"
	"os"
	"strings"

	"github.com/rahairston/go-file-sync/common"
)

type LocalClient struct {
	Path string
}

func NewLocalClient(path string) (*LocalClient, error) {
	var adjustedPath string = path
	if !strings.HasSuffix(path, common.Separator) {
		adjustedPath = path + common.Separator
	}
	return &LocalClient{
		Path: adjustedPath,
	}, ValidatePath(adjustedPath)
}

func ValidatePath(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		panic(err)
	} else if !info.IsDir() {
		panic(errors.New("path provided must be a folder"))
	}

	return nil
}

func (lc LocalClient) GetPath() string {
	return lc.Path
}

func (lc LocalClient) Stat(fileName string) (fs.FileInfo, error) {
	return os.Stat(fileName)
}

func (lc LocalClient) Create(fileName string) (common.SharedFile, error) {
	return os.Create(fileName)
}

func (lc LocalClient) Open(fileName string) (common.SharedFile, error) {
	return os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0644)
}

func (lc LocalClient) CorrectPathSeparator(path string) string {
	if common.Separator == "\\" {
		return strings.ReplaceAll(path, "/", common.Separator)
	} else {
		return strings.ReplaceAll(path, "\\", common.Separator)
	}
}

func (lc LocalClient) GetFileNames(path string, exclusions common.ExcludeObject) []string {
	var result []string
	var adjustedPath string = path
	if !strings.HasSuffix(path, common.Separator) {
		adjustedPath = path + common.Separator
	}

	entries, err := os.ReadDir(adjustedPath)

	if err != nil {
		log.Fatal(err)
	}

	for _, e := range entries {
		if strings.HasPrefix(e.Name(), ".") {
			continue
		} else if e.IsDir() && !common.ShouldBeExcluded(e.Name(), exclusions.Folders) {
			result = append(result, lc.GetFileNames(adjustedPath+e.Name(), exclusions)...)
		} else if !e.IsDir() && !common.ShouldBeExcluded(e.Name(), exclusions.Files) {
			result = append(result, adjustedPath+e.Name())
		}
	}

	return result
}

func (lc LocalClient) ReadFile(fileName string) ([]byte, error) {
	return os.ReadFile(fileName)
}

func (lc LocalClient) Close() {
}
