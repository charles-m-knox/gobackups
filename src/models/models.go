package models

import "time"

type Backup struct {
	SourcePath  string         `yaml:"sourcePath"`
	Targets     []BackupTarget `yaml:"targets"`
	Interval    time.Duration  `yaml:"interval"`
	Description string         `yaml:"description"`
	ID          string         `yaml:"id"`
	Await       []string       `yaml:"await"`
}

type BackupTarget struct {
	TargetPath         string   `yaml:"targetPath"`
	ID                 string   `yaml:"id"`
	PassphraseFilePath string   `yaml:"passphraseFilePath"`
	Type               string   `yaml:"type"` // rsync, targz, gpgtargz
	Await              []string `yaml:"await"`
}

type BackupLog struct {
	BackupID    string `json:"backupId"`
	TargetID    string `json:"targetId"`
	Command     string `json:"command"`
	Type        string `json:"type"`
	Output      string `json:"output"`
	ErrorOutput string `json:"errorOutput"`
	Success     bool   `json:"success"`
	Timestamp   int64  `json:"timestamp"`
	Date        string `json:"date"`
}

// chartData is to structured format for the data returned when `chart` format is requested
type ChartData struct {
	Labels []string     `json:"labels"`
	Series []DataSeries `json:"series"`
}

// seriesData contains the data for each series in the chart data
type DataSeries struct {
	Label string  `json:"label"`
	Data  []int64 `json:"data"`
}
