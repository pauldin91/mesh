package utils

import "github.com/spf13/viper"

type Config struct {
	CertPath          string `mapstructure:"CERTIFICATE_PATH"`
	CertFile          string `mapstructure:"CERTIFICATE_FILE"`
	CertKey           string `mapstructure:"CERTIFICATE_KEY"`
	HttpServerAddress string `mapstructure:"HTTP_SERVER_ADDRESS"`
	WsReadLimit       int64  `mapstructure:"READ_LIMIT"`
	WsWriteLimit      int64  `mapstructure:"WRITE_LIMIT"`
	WsServerAddress   string `mapstructure:"WS_SERVER_ADDRESS"`
	Peers             string `mapstructure:"PEERS"`
}

func LoadConfig(path string) (Config, error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	var config Config
	err := viper.ReadInConfig()
	if err != nil {
		return config, err
	}

	err = viper.Unmarshal(&config)

	return config, err
}
