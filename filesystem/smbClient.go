package filesystem

import (
	"errors"
	"io/fs"
	"net"
	"os"
	"strings"

	"github.com/rahairston/go-file-sync/common"

	"github.com/hirochachacha/go-smb2"
)

type SmbClient struct {
	s    *smb2.Session
	fs   *smb2.Share
	conn net.Conn
	Path string
}

func SmbConnect(connection common.ConnectionObject) (*SmbClient, error) {
	config := connection.SmbConfig
	path := connection.Path
	var port = config.Port
	if port == "" {
		port = "445"
	}

	conn, err := net.Dial("tcp", config.Host+":"+config.Port)

	if err != nil {
		return nil, err
	}

	d := smb2.Dialer{
		Initiator: &smb2.NTLMInitiator{
			User:     config.Authentication.Username,
			Password: config.Authentication.Password,
		},
	}

	s, err := d.Dial(conn)
	if err != nil {
		return nil, err
	}

	fs, err := s.Mount(config.MountPoint)

	if err != nil {
		return nil, err
	}

	var adjustedPath string = path
	if !strings.HasSuffix(path, "\\") { // Keep \\ since SMB is Windows file pathing
		adjustedPath = path + "\\"
	}

	return &SmbClient{
		s:    s,
		conn: conn,
		fs:   fs,
		Path: adjustedPath,
	}, ValidateSmbPath(fs, adjustedPath)
}

func ValidateSmbPath(fs *smb2.Share, path string) error {
	info, err := fs.Stat(path)

	if err != nil {
		panic(err)
	} else if !info.IsDir() {
		panic(errors.New("path provided must be a folder"))
	}

	return nil
}

func (smbClient SmbClient) GetPath() string {
	return smbClient.Path
}

func (smbClient SmbClient) CorrectPathSeparator(path string) string {
	return strings.ReplaceAll(path, "/", "\\")
}

func (smbClient SmbClient) Stat(fileName string) (fs.FileInfo, error) {
	return smbClient.fs.Stat(fileName)
}

func (smbClient SmbClient) Open(fileName string) (common.SharedFile, error) {
	return smbClient.fs.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0644)
}

func (smbClient SmbClient) Create(fileName string) (common.SharedFile, error) {
	return smbClient.fs.Create(fileName)
}

func (smbClient SmbClient) GetFileNames(path string, exclusions common.ExcludeObject) []string {
	var result []string
	var adjustedPath string = path
	if !strings.HasSuffix(path, "\\") { // Keep \\ since SMB is Windows file pathing
		adjustedPath = path + "\\"
	}

	files, _ := smbClient.fs.ReadDir(path)

	for _, file := range files {
		if strings.HasPrefix(file.Name(), ".") {
			continue
		} else if file.IsDir() && !common.ShouldBeExcluded(file.Name(), exclusions.Folders) {
			result = append(result, smbClient.GetFileNames(adjustedPath+file.Name(), exclusions)...)
		} else if !file.IsDir() && !common.ShouldBeExcluded(file.Name(), exclusions.Files) {
			result = append(result, adjustedPath+file.Name())
		}
	}

	return result
}

func (smbClient SmbClient) ReadFile(fileName string) ([]byte, error) {
	return smbClient.fs.ReadFile(fileName)
}

func (smbClient SmbClient) Close() {
	smbClient.s.Logoff()
	smbClient.conn.Close()
}
