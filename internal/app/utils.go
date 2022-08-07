package app

import (
	"encoding/csv"
	"encoding/json"
	"os"

	"github.com/pacholoamit/google-analytics-activity-monitor/internal/models"
)

func (app Application) writeJSONToCSV(c []models.ChangeHistoryEvent, header []string, destination string) error {

	outputFile, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	writer := csv.NewWriter(outputFile)
	defer writer.Flush()

	if err := writer.Write(header); err != nil {
		return err
	}

	for _, r := range c {
		var csvRow []string
		chString, err := json.Marshal(r.Changes)
		if err != nil {
			return err
		}

		csvRow = append(csvRow, r.UserActorEmail, r.ChangeTime, r.ActorType, string(chString))
		if err := writer.Write(csvRow); err != nil {
			return err
		}
	}
	return nil
}
