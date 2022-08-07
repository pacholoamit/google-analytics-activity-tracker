package app

import "fmt"

type Config struct {
	CsvFile string
}

func (c *Config) ValidateFlags() error {
	if c.CsvFile == "" {
		return fmt.Errorf(`please provide all required flags:
			-clientId
				Google Client ID
			-clientSecret
				Google Client Secret
			-csvFile
				CSV File
			`)
	}
	return nil
}
