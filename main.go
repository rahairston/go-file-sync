package main

import (
	"log"

	"github.com/rahairston/go-file-sync/common"
	"github.com/rahairston/go-file-sync/config"
	"github.com/rahairston/go-file-sync/filesystem"
)

func main() {
	consts := common.GetOSConstants()

	logFile := config.SetLoggingFile(consts)
	defer logFile.Close()

	conf, err := config.BuildBackupConfig(consts)

	if err != nil {
		log.Panic(err)
	}

	if err != nil {
		log.Panic(err)
	}

	srcFs, err := filesystem.Connect(conf.SourceConnection)

	if err != nil {
		log.Panic(err)
	}

	dstFs, err := filesystem.Connect(conf.SourceConnection)

	if err != nil {
		log.Panic(err)
	}

	dirClient, err := filesystem.BuildDirClient(conf, srcFs, dstFs)

	if err != nil {
		log.Panic(err)
	}

	dirClient.SyncFiles()

	log.Println("Sync Complete.")
}
