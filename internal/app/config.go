package app

import "fmt"

type Config struct {
	File string
}

func (c *Config) ValidateFlags() error {
	if c.File == "" {
		return fmt.Errorf(`please provide all required flags:
			-clientId
				Google Client ID
			-clientSecret
				Google Client Secret
			-File
				File output path
			`)
	}
	return nil
}
