package config

import (
	"github.com/spf13/viper"
	"log"
)

type influxdb struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	Timeout int    `json:"timeout"`
}

type server struct {
	Address string `json:"address"`
	Tls     bool   `json:"tls"`
}

type config struct {
	RestServer server   `json:"restServer"`
	Influxdb   influxdb `json:"influxdb"`
}

var (
	Config config = config{

		RestServer: server{
			Address: "localhost:8082",
		},
		Influxdb: influxdb{
			Name:    "SportEvents",
			Address: "http://influxdb:8086",
			Timeout: 1,
		},
	}

	viperCfg *viper.Viper = viper.New()
)

// InitConfig  Initialize all the configuration
func InitConfig() error {
	viperCfg.AddConfigPath("../config")

	if err := viperCfg.ReadInConfig(); err != nil {
		// It's not necessarily a problem when no configuration file is found
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
	}

	if err := viperCfg.Unmarshal(&Config); err != nil {
		log.Printf("Error parsing configuration: %s\n", err)
		return err
	}

	return nil
}
