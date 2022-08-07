package app

import (
	"encoding/csv"
	"encoding/json"
	"os"

	"github.com/pacholoamit/google-analytics-activity-monitor/internal/models"
)

func (app Application) writeJSONToCSV(c []models.ChangeHistoryEvent, header []string, destination string) error {

	outputFile, err := os.Create("./activity" + destination)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	writer := csv.NewWriter(outputFile)
	defer writer.Flush()

	if err := writer.Write(header); err != nil {
		return err
	}

	for _, a := range c {
		var record []string

		switch a.UserActorEmail {
		case "":
			record = append(record, "NOT AVAILABLE")
		default:
			record = append(record, a.UserActorEmail)
		}

		record = append(record, a.ChangeTime)

		switch a.ActorType {
		case "":
			record = append(record, "NOT AVAILABLE")
		default:
			record = append(record, a.ActorType)
		}

		chString, err := json.Marshal(a.Changes)

		if err != nil {
			return err
		}
		record = append(record, string(chString))

		if err := writer.Write(record); err != nil {
			return err
		}
	}
	return nil
}
