// Phantasm Crawler
//
// "Drawing on his battles fought across many worlds and many times, phantasms of the Chaos Knight rise up to quell all
// who oppose him"
//
// "Summons several phantasmal copies of the Chaos Knight from alternate dimensions. The phantasms are illusions that
// deal 100% damage, but take 350% damage."
//
// Phantasm crawls inventory for item and delivery tracking. Hopefully, by summoning multiple instances of the crawler
// will provide better steam inventory raw data.
//
// crawler.go
//	- script is intended for serverless functions to work around with ip rate limits during peak usage.
// 	- publishes raw inventory data to target webhook url.

package phantasm

type Config struct {
	WebhookURL string `envconfig:"DG_PHANTASM_WEBHOOK_URL"`
	Secret     string
}
