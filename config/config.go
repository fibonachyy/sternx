package config

import (
	"log"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

type Config struct {
	Postgres struct {
		Host           string `yaml:"Host"`
		User           string `yaml:"User"`
		Password       string `yaml:"Password"`
		DB             string `yaml:"DB"`
		MigrationsPath string `yaml:"MigrationsPath"`
	} `yaml:"Postgres"`
	Grpc struct {
		Port string `yaml:"Port"`
	}
	Tls struct {
		Cert string `yaml:"Cert"`
		Key  string `yaml:"Key"`
	}
}

func ReadConfig(configFile string) Config {
	c := &Config{}
	err := c.Unmarshal(c, configFile)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
	return *c
}

func (c *Config) Unmarshal(rawVal interface{}, fileName string) error {
	viper.SetConfigFile(fileName)
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}
	var input interface{} = viper.AllSettings()
	config := defaultDecoderConfig(rawVal)
	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return err
	}
	return decoder.Decode(input)
}

func defaultDecoderConfig(output interface{}) *mapstructure.DecoderConfig {
	c := &mapstructure.DecoderConfig{
		Metadata:         nil,
		Result:           output,
		WeaklyTypedInput: true,
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			mapstructure.StringToTimeDurationHookFunc(),
			mapstructure.StringToSliceHookFunc(","),
		),
	}
	return c
}
