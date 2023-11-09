package common

import (
	"io/fs"
	"os"
	"time"
)

type SharedFile interface {
	Close() error
	Name() string
	Read(b []byte) (n int, err error)
	Write([]byte) (int, error)
	Stat() (os.FileInfo, error)
}

type FileSystem interface {
	GetPath() string
	CorrectPathSeparator(path string) string
	Stat(fileName string) (fs.FileInfo, error)
	Create(fileName string) (SharedFile, error)
	OpenFile(fileName string, flag int) (SharedFile, error)
	GetFileNames(path string, exclusions ExcludeObject) []string
	ReadFile(fileName string) ([]byte, error)
	Truncate(fileName string, newsize int64) error
	Close()
}

type FileConfigType string

const (
	Local FileConfigType = "local"
	Smb   FileConfigType = "smb"
)

type Authentication struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type SmbConfig struct {
	Host           string         `json:"host"`
	Port           string         `json:"port"`
	Authentication Authentication `json:"authentication"`
	MountPoint     string         `json:"mountPoint"`
}

type ConnectionObject struct {
	Type      FileConfigType `json:"type"`
	Path      string         `json:"path"`
	SmbConfig SmbConfig      `json:"smbConfig"`
}

type ExcludeObject struct {
	Files   []string `json:"files"`
	Folders []string `json:"folders"`
}

type SyncObject struct {
	SourceConnections []ConnectionObject `json:"sources"`
	DstConnection     ConnectionObject   `json:"destination"`
	Exclusions        ExcludeObject      `json:"exclude"`
}

type LastModifiedObject struct {
	Path         string
	LastModified time.Time
}
