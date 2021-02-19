package commands

import (
	"gobackups/config"
	"gobackups/constants"
	"gobackups/models"

	"bytes"
	"fmt"
	"os/exec"
)

// GetCommand builds the command for the specified backup and its backup
// target.
// TODO: integrate this into the below functions properly
func GetCommand(backup models.Backup, target models.BackupTarget) string {
	switch target.Type {
	case constants.BackupTargetTypeRsync:
		return fmt.Sprintf("rsync -avP %v %v", backup.SourcePath, target.TargetPath)
	case constants.BackupTargetTypeTargz:
		return fmt.Sprintf("tar -czvf %v %v", target.TargetPath, backup.SourcePath)
	case constants.BackupTargetTypeGpgTargz:
		return fmt.Sprintf(`tar -czvf - "%v" | gpg -c --batch --passphrase-file="%v" >"%v"`, backup.SourcePath, target.PassphraseFilePath, target.TargetPath)
	}
	return ""
}

// RunCommand allows arbitrary commands to be passed in for execution
func RunCommand(commandName string, commandArgs []string) (output string, err error) {
	cmd := exec.Command(commandName, commandArgs...)
	var out bytes.Buffer
	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		return "", fmt.Errorf("failure running cmd: %v", err.Error())
	}
	return out.String(), nil
}

// BackupRsync performs a common rsync command to do a backup of sourcePath
// to targetPath
func BackupRsync(conf config.Config, sourcePath string, targetPath string) (output string, err error) {
	// output, err = RunCommand("rsync", []string{"-avP", sourcePath, targetPath})
	// if err != nil {
	// 	return output, fmt.Errorf("failed to run rsync backup: %v", err.Error())
	// }
	// return output, nil

	cmd := fmt.Sprintf(
		`ls %v && cd $(dirname %v) && rsync -avP %v %v`,
		sourcePath,
		targetPath,
		sourcePath,
		targetPath,
	)
	fmt.Printf("%v\n", cmd)
	outputBytes, err := exec.Command(conf.ShellCommand, "-c", cmd).CombinedOutput()
	output = string(outputBytes)
	if err != nil {
		return output, fmt.Errorf("failed to run rsync tar backup: %v", err.Error())
	}
	return output, nil
}

// BackupTar performs a common tar.gz compression command to do a compressed
// archive backup of sourcePath to targetPath
func BackupTar(conf config.Config, sourcePath string, targetPath string) (output string, err error) {
	// output, err = RunCommand("tar", []string{"-czvf", targetPath, sourcePath})
	// if err != nil {
	// 	return output, fmt.Errorf("failed to run tar backup: %v", err.Error())
	// }
	// return output, nil

	cmd := fmt.Sprintf(
		`ls %v && cd $(dirname %v) && tar -czf %v %v`,
		sourcePath,
		targetPath,
		targetPath,
		sourcePath,
	)
	fmt.Printf("%v\n", cmd)
	outputBytes, err := exec.Command(conf.ShellCommand, "-c", cmd).CombinedOutput()
	output = string(outputBytes)
	if err != nil {
		return output, fmt.Errorf("failed to run rsync tar backup: %v", err.Error())
	}
	return output, nil
}

// BackupEncryptedTar performs a symmetric gpg-encrypted tar.gz compression
// command to do a compressed archive backup of sourcePath to targetPath,
// using a passphrase contained in a file
func BackupEncryptedTar(conf config.Config, sourcePath string, targetPath string, passphraseFilePath string) (output string, err error) {
	cmd := fmt.Sprintf(
		`ls %v && cd $(dirname %v) && ls %v && tar -czf - "%v" | gpg -c --pinentry loopback --batch --passphrase-file="%v" >"%v"`,
		sourcePath,
		targetPath,
		passphraseFilePath,
		sourcePath,
		passphraseFilePath,
		targetPath,
	)
	fmt.Printf("%v\n", cmd)
	outputBytes, err := exec.Command(conf.ShellCommand, "-c", cmd).CombinedOutput()
	output = string(outputBytes)
	if err != nil {
		return output, fmt.Errorf("failed to run encrypted tar backup: %v", err.Error())
	}
	return output, nil
}

/*
a: depends on b
b: depends on c

a depends on c and on b

awaits[a] = [b]
awaits[b] = [c]

for each awaitId in awaits[a]:
    for each awaitSubId in awaits[awaitId]:
        check if backupsInProgress[awaitSubId] is busy
        if busy, return true
        if not busy, return false
        for each awaitSubSubId in awaits[awaitSubId]:
            if there are no more dependencies, then
			return
*/
func IsBackupAwaiting(backupID string, backupsInProgress map[string]bool, backupAwaitDefinitions map[string][]string) bool {
	for backupIDInProgress, isInProgress := range backupsInProgress {
		// case: if the current backup ID is in progress, return true
		if backupID == backupIDInProgress && isInProgress {
			return true
		}
		// case: recursively check all the await definitions for backupID
		for _, backupDependencyID := range backupAwaitDefinitions[backupID] {
			if IsBackupAwaiting(backupDependencyID, backupsInProgress, backupAwaitDefinitions) {
				return true
			}
		}
	}
	return false
}
