package common

import (
	"os"
	"runtime"
)

type BackupConstants struct {
	LoggingLocation string
	ConfigLocation  string
}

const (
	Separator       string = string(os.PathSeparator)
	LastRunFileName string = "last_run.conf"
)

func GetOSConstants() *BackupConstants {
	opsys := runtime.GOOS
	switch opsys {
	case "windows":
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return &BackupConstants{
			LoggingLocation: home + "\\AppData\\Local\\go-file-sync\\",
			ConfigLocation:  home + "\\AppData\\Local\\go-file-sync\\",
		}
	default:
		return &BackupConstants{
			LoggingLocation: "/var/log/go-file-sync/",
			ConfigLocation:  "/etc/go-file-sync/",
		}
	}
}
