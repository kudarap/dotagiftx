package envconf

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

// EnvPrefix default env prefix APP.
var EnvPrefix = "APP"

// Load parses .env values into a struct.
func Load(in interface{}) error {
	// Load env file.
	if err := godotenv.Load(); err != nil {
		return err
	}

	// Bind env values.
	if err := envconfig.Process(EnvPrefix, in); err != nil {
		return err
	}

	return nil
}
