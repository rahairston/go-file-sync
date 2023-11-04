package config

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"time"

	"github.com/rahairston/go-file-sync/common"
)

func BuildBackupConfig(consts *common.BackupConstants) (*common.SyncObject, error) {
	_, err := os.Stat(consts.ConfigLocation)
	if err != nil {
		return nil, err
	}
	_, err = os.Stat(consts.LoggingLocation)
	if err != nil {
		return nil, err
	}

	return parseJSONConfig(consts)
}

func parseJSONConfig(consts *common.BackupConstants) (*common.SyncObject, error) {
	jsonFile, err := os.Open(consts.ConfigLocation + "config.json")

	if err != nil {
		return nil, err
	}

	defer jsonFile.Close()

	jsonData, err := io.ReadAll(jsonFile)

	if err != nil {
		return nil, err
	}

	config := common.SyncObject{}

	json.Unmarshal(jsonData, &config)

	return &config, nil
}

func SetLoggingFile(consts *common.BackupConstants) *os.File {
	now := time.Now().UTC()

	logFile, err := os.OpenFile(consts.LoggingLocation+now.Format("2006-01-02")+".log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		log.Panic(err)
	}

	log.SetOutput(logFile)

	log.Println("Starting...")

	return logFile
}
