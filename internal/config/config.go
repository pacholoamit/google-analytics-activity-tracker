package config

type Config struct {
	clientId     string
	clientSecret string
	redirectURL  string
	code         string
}

func New(clientId string, clientSecret string, redirectUrl string) *Config {
	return &Config{
		clientId:     clientId,
		clientSecret: clientSecret,
		redirectURL:  redirectUrl,
		code:         "",
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

func (c *Config) SetCode(code string) {
	c.code = code
}

func (c *Config) GetCode() string {
	return c.code
}
