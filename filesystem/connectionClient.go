package filesystem

import (
	"github.com/rahairston/go-file-sync/common"
)

func Connect(connectionConfig common.ConnectionObject) (common.FileSystem, error) {

	switch connectionConfig.Type {
	case common.Smb:
		return SmbConnect(connectionConfig)
	case common.Local:
		return NewLocalClient(connectionConfig.Path)
	default:
		return nil, nil
	}
}
