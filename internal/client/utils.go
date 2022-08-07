package client

import "fmt"

func (c *Config) ValidateFlags() error {
	if c.ClientId == "" || c.ClientSecret == "" {
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
