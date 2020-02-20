package config

import "github.com/kelseyhightower/envconfig"

// Config structures
type Config struct {
	Client struct {
		ID     string `envconfig:"ClIENT_ID"`
		Secret string `envconfig:"CLIENT_SECRET"`
	}
}

// SpotifyConfig returns configration variables in a structure
func SpotifyConfig() {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		panic(err)
	}
	return cfg
}
