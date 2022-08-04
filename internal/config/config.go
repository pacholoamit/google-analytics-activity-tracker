package config

type Config struct {
	clientId     string
	clientSecret string
	redirectURL  string
}

func New(clientId string, clientSecret string, redirectUrl string) *Config {
	return &Config{
		clientId:     clientId,
		clientSecret: clientSecret,
		redirectURL:  redirectUrl,
	}
}

func (c *Config) GetClientId() string {
	return c.clientId
}

func (c *Config) GetClientSecret() string {
	return c.clientSecret
}

func (c *Config) GetRedirectURL() string {
	return c.redirectURL
}
