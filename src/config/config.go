package config

import (
	"fmt"
	"gobackups/models"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Backups          []models.Backup `yaml:"backups"`
	LogFile          string          `yaml:"logFile"`
	MaxLogEntries    int64           `yaml:"maxLogEntries"`
	ShellCommand     string          `yaml:"shellCommand"`
	AwaitDefinitions map[string][]string
}

// LoadConfig reads from a provided yaml-formatted configuration filename
func LoadConfig(logFile string) (conf Config, err error) {
	// read from config file
	confData, err := ioutil.ReadFile(logFile)
	if err != nil {
		return conf, fmt.Errorf("failed to read config file %v: %v", logFile, err.Error())
	}

	err = yaml.Unmarshal(confData, &conf)
	if err != nil {
		return conf, fmt.Errorf("failed to parse config file %v: %v", logFile, err.Error())
	}

	conf.AwaitDefinitions = conf.BuildAwaitDefinitions()

	return conf, nil
}

// BuildAwaitDefinitions looks through every single configured target and source,
// and if that target/source has an "Await" rule defined,
// it has to wait for all of its dependencies to be completed
func (conf *Config) BuildAwaitDefinitions() map[string][]string {
	result := make(map[string][]string)
	for _, backupSource := range conf.Backups {
		if len(backupSource.Await) > 0 {
			result[backupSource.ID] = append(result[backupSource.ID], backupSource.Await...)
		}
		for _, backupTarget := range conf.Backups {
			if len(backupTarget.Await) > 0 {
				result[backupTarget.ID] = append(result[backupTarget.ID], backupTarget.Await...)
			}
		}
	}
	return result
}
