package phantasm

const (
	defaultConfigAddr       = "http://localhost:8000/crawler/phantasm"
	defaultConfigWebhookURL = "http://localhost:8000/webhook/phantasm"
	defaultConfigSecret     = "reality-rift"
	defaultConfigPath       = "./.localdata/phantasm"
)

type Config struct {
	Addrs      []string
	WebhookURL string `envconfig:"WEBHOOK_URL"`
	Secret     string
	Path       string
}

func (c Config) setDefault() Config {
	if len(c.Addrs) == 0 {
		c.Addrs = []string{defaultConfigAddr}
	}
	if c.WebhookURL == "" {
		c.WebhookURL = defaultConfigWebhookURL
	}
	if c.Secret == "" {
		c.Secret = defaultConfigSecret
	}
	if c.Path == "" {
		c.Path = defaultConfigPath
	}
	return c
}
