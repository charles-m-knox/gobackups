package logger

import (
	"fmt"
	"gobackups/commands"
	"gobackups/config"
	"gobackups/models"
	"io/ioutil"
	"log"
	"time"

	"encoding/json"
)

var currentlyWriting bool

func LogBackupLoopable(conf config.Config, backup models.Backup, backupTarget models.BackupTarget, logFile string, success bool, output string, errorOutput string) (err error) {
	currentlyWriting = true

	// TODO: use the GetLog function instead of this to reduce duplication
	// read from log file
	backupLogFileData, err := ioutil.ReadFile(logFile)
	if err != nil {
		log.Printf("failed to read backup log file %v, it may not exist: %v", logFile, err.Error())
	}

	// unmarshal into []BackupLog
	var backupLog []models.BackupLog

	// only unmarshal if log file exists
	if err == nil {
		err = json.Unmarshal(backupLogFileData, &backupLog)
		if err != nil {
			return fmt.Errorf("failed to parse backup log file %v: %v", logFile, err.Error())
		}
	}

	// push this result into the data
	timeNow := time.Now().Local().UTC()
	backupLog = append([]models.BackupLog{models.BackupLog{
		BackupID:  backup.ID,
		TargetID:  backupTarget.ID,
		Command:   commands.GetCommand(backup, backupTarget),
		Type:      backupTarget.Type,
		Output:    output,
		Success:   success,
		Timestamp: timeNow.Unix(),
		Date:      timeNow.Format(time.RFC3339),
	}}, backupLog...)

	// check if the log is getting too big
	if int64(len(backupLog)) > conf.MaxLogEntries {
		backupLog = backupLog[0:conf.MaxLogEntries]
	}

	// marshal for writing to log file
	backupLogJSON, err := json.Marshal(backupLog)
	if err != nil {
		return fmt.Errorf("failed to prep updated backup log file %v: %v", logFile, err.Error())
	}

	// write to log file
	err = ioutil.WriteFile(logFile, backupLogJSON, 0644) // TODO: make perm configurable
	if err != nil {
		return fmt.Errorf("failed to save updated backup log file %v: %v", logFile, err.Error())
	}
	currentlyWriting = false
	return nil
}

func LogBackup(conf config.Config, backup models.Backup, backupTarget models.BackupTarget, logFile string, success bool, output string, errorOutput string) (err error) {
	for {
		if !currentlyWriting {
			err = LogBackupLoopable(conf, backup, backupTarget, logFile, success, output, errorOutput)
			if err != nil {
				return err
			}
			return nil
		}
		log.Printf("sleeping 1 second to log backup entry to avoid concurrent writes to json log")
		time.Sleep(1 * time.Second)
	}
}

// GetLog reads the log file and returns it as a string as well as as a
// marshaled object
func GetLog(logFile string) (logStr string, backupLogs []models.BackupLog, err error) {
	// read from log file
	backupLogFileData, err := ioutil.ReadFile(logFile)
	if err != nil {
		return logStr, backupLogs, fmt.Errorf("failed to read backup log file %v: %v", logFile, err.Error())
	}

	// unmarshal into []BackupLog
	var backupLog []models.BackupLog
	err = json.Unmarshal(backupLogFileData, &backupLog)
	if err != nil {
		return logStr, backupLogs, fmt.Errorf("failed to parse backup log file %v: %v", logFile, err.Error())
	}

	return string(backupLogFileData), backupLog, nil
}
