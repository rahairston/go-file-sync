package common

import (
	"os"
	"runtime"
)

type SyncConstants struct {
	LoggingLocation string
	ConfigLocation  string
}

const (
	Separator       string = string(os.PathSeparator)
	LastRunFileName string = "last_run.conf"
)

func GetOSConstants() *SyncConstants {
	opsys := runtime.GOOS
	switch opsys {
	case "windows":
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return &SyncConstants{
			LoggingLocation: home + "\\AppData\\Local\\go-file-sync\\",
			ConfigLocation:  home + "\\AppData\\Local\\go-file-sync\\",
		}
	default:
		return &SyncConstants{
			LoggingLocation: "/var/log/go-file-sync/",
			ConfigLocation:  "/etc/go-file-sync/",
		}
	}
}
