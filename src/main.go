package main

import (
	"gobackups/commands"
	"gobackups/config"
	"gobackups/constants"
	"gobackups/handlers"
	"gobackups/logger"
	"gobackups/models"

	"fmt"
	"log"
	"net/http"
	"time"
)

var backupIDsInProgress map[string]bool

func monitorBackup(backupConfig models.Backup, conf config.Config) {
	log.Printf("starting backup %v with %v interval", backupConfig.ID, backupConfig.Interval)
	for {
		// check if there is a backup in progress that this source depends on
		// TODO: configurable await time
		if commands.IsBackupAwaiting(backupConfig.ID, backupIDsInProgress, conf.AwaitDefinitions) {
			log.Printf("backup source %v is awaiting another backup, checking again soon...", backupConfig.ID)
			time.Sleep(5 * time.Second)
			continue
		}
		log.Printf("backup %v starting...", backupConfig.ID)
		// set this backup ID as "in progress"
		backupIDsInProgress[backupConfig.ID] = true
		for _, backupTarget := range backupConfig.Targets {
			for {
				// check if there is a backup in progress that this source depends on
				// TODO: configurable await time & limited retries
				if commands.IsBackupAwaiting(backupTarget.ID, backupIDsInProgress, conf.AwaitDefinitions) {
					log.Printf("backup target %v is awaiting another backup, checking again soon...", backupConfig.ID)
					time.Sleep(5 * time.Second)
					continue
				}
				break
			}
			// set this backup ID as "in progress"
			backupIDsInProgress[backupTarget.ID] = true

			// decision tree based on what kind of backup this is
			switch backupTarget.Type {
			case constants.BackupTargetTypeRsync:
				errorOutput := ""
				success := true
				output, err := commands.BackupRsync(conf, backupConfig.SourcePath, backupTarget.TargetPath)
				if err != nil {
					errorOutput = err.Error()
					log.Printf("error: %v", errorOutput)
					success = false
				}
				err = logger.LogBackup(conf, backupConfig, backupTarget, conf.LogFile, success, output, errorOutput)
				if err != nil {
					log.Printf(`backup %v %v target "%v" failed to log: %v`, backupConfig.ID, backupTarget.Type, backupTarget.TargetPath, err.Error())
				}
			case constants.BackupTargetTypeTargz:
				errorOutput := ""
				success := true
				output, err := commands.BackupTar(conf, backupConfig.SourcePath, backupTarget.TargetPath)
				if err != nil {
					errorOutput = err.Error()
					log.Printf("error: %v", errorOutput)
					success = false
				}
				err = logger.LogBackup(conf, backupConfig, backupTarget, conf.LogFile, success, output, errorOutput)
				if err != nil {
					log.Printf(`backup %v %v target "%v" failed to log: %v`, backupConfig.ID, backupTarget.Type, backupTarget.TargetPath, err.Error())
				}
			case constants.BackupTargetTypeGpgTargz:
				errorOutput := ""
				success := true

				if backupTarget.PassphraseFilePath == "" {
					log.Printf("warning: backup %v target %v has no passphraseFilePath defined", backupConfig.ID, backupTarget.ID)
					return
				}

				output, err := commands.BackupEncryptedTar(conf, backupConfig.SourcePath, backupTarget.TargetPath, backupTarget.PassphraseFilePath)
				if err != nil {
					errorOutput = err.Error()
					log.Printf("error: %v", errorOutput)
					success = false
				}
				err = logger.LogBackup(conf, backupConfig, backupTarget, conf.LogFile, success, output, errorOutput)
				if err != nil {
					log.Printf(`backup %v %v target "%v" failed to log: %v`, backupConfig.ID, backupTarget.Type, backupTarget.TargetPath, err.Error())
				}
			}
			// this specific backup target is completed
			backupIDsInProgress[backupTarget.ID] = false
		}
		log.Printf("backup %v finished, waiting %v...", backupConfig.ID, backupConfig.Interval)
		// all backup targets for this source are finished
		backupIDsInProgress[backupConfig.ID] = false
		time.Sleep(time.Duration(backupConfig.Interval))
	}
}

func main() {
	backupIDsInProgress = make(map[string]bool)
	log.Printf("starting at %v", time.Now())

	// load config
	conf, err := config.LoadConfig("config.yml")
	if err != nil {
		log.Fatal(err)
	}

	// spawn a goroutine that checks each backup
	for _, backupSource := range conf.Backups {
		go monitorBackup(backupSource, conf)
	}

	// push the config to the handlers module so it has acccess to config
	// information
	handlers.SetConfig(&conf)

	http.HandleFunc(fmt.Sprintf("%v/logs", constants.ContentPrefixURL), handlers.LogViewHandler)

	err = http.ListenAndServe(constants.ServerInterfaceBindURL, nil)
	if err != nil {
		log.Fatal(err)
	}
}
