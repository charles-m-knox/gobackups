package commands

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsBackupAwaiting(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		BackupID               string
		BackupsInProgress      map[string]bool
		BackupAwaitDefinitions map[string][]string
		ExpectedResult         bool
		TestName               string
	}{
		{
			BackupID: "a",
			BackupsInProgress: map[string]bool{
				"a": true,
			},
			BackupAwaitDefinitions: map[string][]string{
				"b": {"a"},
			},
			ExpectedResult: true,
			TestName:       "main id is in progress",
		},
		{
			BackupID: "a",
			BackupsInProgress: map[string]bool{
				"a": false,
			},
			BackupAwaitDefinitions: map[string][]string{
				"b": {"a"},
			},
			ExpectedResult: false,
			TestName:       "main id is not in progress",
		},
		{
			BackupID: "a",
			BackupsInProgress: map[string]bool{
				"a": true,
			},
			BackupAwaitDefinitions: map[string][]string{
				"b": {"a"},
			},
			ExpectedResult: true,
			TestName:       "b awaits a which is in progress",
		},
		{
			BackupID: "a",
			BackupsInProgress: map[string]bool{
				"d": true,
			},
			BackupAwaitDefinitions: map[string][]string{
				"a": {"b"},
				"b": {"c"},
				"c": {"d"},
			},
			ExpectedResult: true,
			TestName:       "a>b>c>d, d is in progress",
		},
		{
			BackupID: "b",
			BackupsInProgress: map[string]bool{
				"d": true,
				"c": true,
			},
			BackupAwaitDefinitions: map[string][]string{
				"a": {"b"},
				"b": {"c"},
				"c": {"d"},
			},
			ExpectedResult: true,
			TestName:       "a>b>c>d, d+c are in progress",
		},
		{
			BackupID: "a",
			BackupsInProgress: map[string]bool{
				"g": true,
			},
			BackupAwaitDefinitions: map[string][]string{
				"a": {"b", "c", "d", "e", "f"},
				"b": {"c", "d", "e", "f"},
				"c": {"d", "e", "f"},
				"d": {"e", "f"},
				"e": {"f"},
				"f": {"g"},
			},
			ExpectedResult: true,
			TestName:       "abcde->f, f->g, g is in progress",
		},
		{
			BackupID: "a",
			BackupsInProgress: map[string]bool{
				"g": false,
				"h": true,
			},
			BackupAwaitDefinitions: map[string][]string{
				"a": {"b", "c", "d", "e", "f"},
				"b": {"c", "d", "e", "f"},
				"c": {"d", "e", "f"},
				"d": {"e", "f"},
				"e": {"f"},
				"f": {"g"},
			},
			ExpectedResult: false,
			TestName:       "abcde->f, f->g, h is in progress",
		},
	}

	for _, test := range tests {
		assert.Equal(
			test.ExpectedResult,
			IsBackupAwaiting(
				test.BackupID,
				test.BackupsInProgress,
				test.BackupAwaitDefinitions,
			),
			test.TestName,
		)
	}
}
