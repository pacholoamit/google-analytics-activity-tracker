package app

import (
	"encoding/csv"
	"os"

	"github.com/pacholoamit/google-analytics-activity-monitor/internal/models"
)

func (app Application) writeJSONToCSV(c []models.ChangeHistoryEvent, header []string, destination string) error {

	outputFile, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	// 4. Write the header of the CSV file and the successive rows by iterating through the JSON struct array
	writer := csv.NewWriter(outputFile)
	defer writer.Flush()

	if err := writer.Write(header); err != nil {
		return err
	}

	for _, r := range c {
		var csvRow []string
		csvRow = append(csvRow, r.ChangeTime, r.UserActorEmail, r.ActorType)
		if err := writer.Write(csvRow); err != nil {
			return err
		}
	}
	return nil
}
