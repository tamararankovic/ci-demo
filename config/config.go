package config

import "os"

type Config struct {
	Address       string
	JaegerAddress string
}

func GetConfig() Config {
	return Config{
		Address:       os.Getenv("PRODUCT_SERVICE_ADDRESS"),
		JaegerAddress: os.Getenv("JAEGER_ADDRESS"),
	}
}
