package config

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/rahairston/go-file-sync/common"
)

func BuildSyncConfig(consts *common.SyncConstants) (*common.SyncObject, error) {
	_, err := os.Stat(consts.ConfigLocation)
	if err != nil {
		log.Fatal("config file not found")
		return nil, err
	}
	_, err = os.Stat(consts.LoggingLocation)
	if err != nil {
		log.Fatal("log file not found")
		return nil, err
	}

	return parseJSONConfig(consts)
}

func parseJSONConfig(consts *common.SyncConstants) (*common.SyncObject, error) {
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

func ParseLastModifiedFile(consts *common.SyncConstants, srcPath string) map[string]string {
	fileMod := make(map[string]string)
	filePath := strings.ReplaceAll(strings.ReplaceAll(srcPath, "/", "_"), "\\", "_")

	file, err := os.Open(consts.ConfigLocation + filePath + ".json")
	if err != nil {
		return fileMod
	}

	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return fileMod
	}

	err = json.Unmarshal([]byte(data), &fileMod)

	if err != nil {
		return fileMod
	}

	return fileMod
}

func WriteLastModifiedFile(consts *common.SyncConstants, srcPath string, modData map[string]string) {
	filePath := strings.ReplaceAll(strings.ReplaceAll(srcPath, "/", "_"), "\\", "_")
	file, err := os.OpenFile(consts.ConfigLocation+filePath+".json", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	jsonString, err := json.Marshal(modData)

	if err != nil {
		log.Fatal(err)
	}

	file.Write(jsonString)
}

func SetLoggingFile(consts *common.SyncConstants) *os.File {
	now := time.Now().UTC()

	logFile, err := os.OpenFile(consts.LoggingLocation+now.Format("2006-01-02")+".log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		log.Panic(err)
	}

	log.SetOutput(logFile)

	log.Println("Starting File Sync...")

	return logFile
}
