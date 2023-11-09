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

	conf, err := config.BuildSyncConfig(consts)

	if err != nil {
		log.Panic(err)
	}

	dstFs, err := filesystem.Connect(conf.DstConnection)

	if err != nil {
		log.Panic(err)
	}

	for _, source := range conf.SourceConnections {
		srcFs, err := filesystem.Connect(source)

		if err != nil {
			log.Panic(err)
		}

		lastFileMod := config.ParseLastModifiedFile(consts, srcFs.GetPath())

		dirClient, err := filesystem.BuildDirClient(conf, srcFs, dstFs)

		if err != nil {
			log.Panic(err)
		}

		updateFileMod := dirClient.SyncFiles(lastFileMod)

		config.WriteLastModifiedFile(consts, srcFs.GetPath(), updateFileMod)
	}

	log.Println("Sync Complete.")
}
