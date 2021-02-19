package handlers

import (
	"encoding/json"
	"gobackups/config"
	"gobackups/constants"
	"gobackups/logger"
	"gobackups/models"
	"sort"
	"strconv"

	"log"
	"net/http"
)

var conf *config.Config

func SetConfig(config *config.Config) {
	conf = config
}

// HTTPFunc is a function type that matches what is typically passed
// into an HTTP API endpoint
type HTTPFunc func(http.ResponseWriter, *http.Request)

// ProcessCORS is a decorator function that enables CORS preflight responses and
// sets CORS headers.
// Usage:
// func GetSomething(writer http.ResponseWriter, req *http.Request) {
//     helpers.ProcessCORS(writer, req, func(writer http.ResponseWriter, req *http.Request) {
//         // do some logic with writer and req
//     })
// }
func ProcessCORS(writer http.ResponseWriter, req *http.Request, fn HTTPFunc) {
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	writer.Header().Set("Access-Control-Allow-Headers", constants.AllowedHeaders)

	// handle preflight options requests
	if req.Method == "OPTIONS" {
		writer.WriteHeader(http.StatusOK)
		return
	}

	fn(writer, req)
}

func LogViewHandler(w http.ResponseWriter, req *http.Request) {
	ProcessCORS(w, req, func(writer http.ResponseWriter, req *http.Request) {
		// return the log file
		backupLogFileStr, backupLogData, err := logger.GetLog(conf.LogFile)
		if err != nil {
			log.Printf("failed to retrieve logs: %v", err.Error())
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			_, err := w.Write([]byte("[]"))
			if err != nil {
				log.Printf("failed to write http response: %v", err.Error())
			}
			return
		}

		logEntryResultLimit := 100

		queryLimit := req.URL.Query().Get("c")
		if (len(queryLimit)) != 0 {
			queryConverted, err := strconv.Atoi(queryLimit)
			if err != nil {
				log.Printf("bad user request for limit query: %v", err.Error())
				w.Header().Set("Content-Type", "text/plain")
				w.WriteHeader(http.StatusBadRequest)
				_, err := w.Write([]byte("invalid 'c' query parameter, please specify an integer"))
				if err != nil {
					log.Printf("failed to write http response: %v", err.Error())
				}
				return
			}

			if queryConverted > -1 {
				logEntryResultLimit = queryConverted
			}
		}

		viewType := req.URL.Query().Get("t")
		if (len(viewType)) != 0 {
			// the user has specified a specific view
			switch viewType {
			case "chart":
				// return as series and labels
				chartData := models.ChartData{}

				// track all unique date labels
				uniqueDates := make(map[string]string)
				uniqueBackupTargetIDs := make(map[string][]string)
				uniqueBackupTargetIDSuccessValues := make(map[string][]bool)
				for _, backupLogEntry := range backupLogData {
					uniqueDates[backupLogEntry.Date] = backupLogEntry.Date
					// keep track of the dates that have data for each unique backup
					// target ID
					uniqueBackupTargetIDs[backupLogEntry.TargetID] = append(uniqueBackupTargetIDs[backupLogEntry.TargetID], backupLogEntry.Date)

					// keep track of the success state for the backup log entry
					uniqueBackupTargetIDSuccessValues[backupLogEntry.TargetID] = append(uniqueBackupTargetIDSuccessValues[backupLogEntry.TargetID], backupLogEntry.Success)
				}

				// get keys from unique labels map
				for key := range uniqueDates {
					chartData.Labels = append(chartData.Labels, key)
				}

				// sort ascending
				sort.Strings(chartData.Labels)

				// restrict to fall within query limit length
				if logEntryResultLimit > 0 && len(chartData.Labels) > logEntryResultLimit {
					chartData.Labels = chartData.Labels[len(chartData.Labels)-logEntryResultLimit : len(chartData.Labels)]
				}

				for tgtID, tgtDates := range uniqueBackupTargetIDs {
					seriesValues := []int64{}
					// iterate over each unique backup log target to build a
					// series for it
					for _, backupDate := range chartData.Labels {
						// check if a backupDate exists for this backup target
						backupSuccessAsInt := 0

						for i, tgtBackupDate := range tgtDates {
							if tgtBackupDate == backupDate {
								// found a backup date match, so now
								// retrieve its success value from the
								// uniqueBackupTargetIDSuccessValues map
								if uniqueBackupTargetIDSuccessValues[tgtID][i] {
									backupSuccessAsInt = 1
								} else {
									backupSuccessAsInt = -1
								}
							}
						}
						seriesValues = append(seriesValues, int64(backupSuccessAsInt))
					}
					chartData.Series = append(chartData.Series, models.DataSeries{
						Label: tgtID,
						Data:  seriesValues,
					})
				}

				// marshal to json
				chartDataOutput, err := json.Marshal(chartData)
				if err != nil {
					log.Printf("failed to marshal chart data: %v", err.Error())

					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusInternalServerError)
					_, err = w.Write([]byte("[]"))
					if err != nil {
						log.Printf("failed to write http response: %v", err.Error())
					}
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, err = w.Write(chartDataOutput)
				if err != nil {
					log.Printf("failed to write http response: %v", err.Error())
				}
				return
			}
		}

		// restrict to fall within query limit length
		if logEntryResultLimit > 0 && len(backupLogData) > logEntryResultLimit {
			limitedLogData, err := json.Marshal(backupLogData[0:logEntryResultLimit])
			if err != nil {
				log.Printf("failed to limit logs: %v", err.Error())
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				_, err := w.Write([]byte("[]"))
				if err != nil {
					log.Printf("failed to write http response: %v", err.Error())
				}
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, err = w.Write(limitedLogData)
			if err != nil {
				log.Printf("failed to write http response: %v", err.Error())
			}
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err = w.Write([]byte(backupLogFileStr))
		if err != nil {
			log.Printf("failed to write http response: %v", err.Error())
		}
	})
}
