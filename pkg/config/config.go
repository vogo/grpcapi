package config

import (
	"flag"
	"strconv"

	"github.com/go2s/o2m"
	"github.com/golang/glog"
	"github.com/spf13/viper"
)

var (
	configFile = flag.String("config", "config", "config file path")

	//SecretKey secret key
	SecretKey = "ADkajfdJIALDSFJJkadf"

	//APIGatewayAddress api gateway address
	APIGatewayAddress = ":8080"

	//EchoServiceAddress echo service address
	EchoServiceAddress = ":9001"

	//HelloServiceAddress hello service address
	HelloServiceAddress = ":9002"

	//ServiceKeyHello hello
	ServiceKeyHello = "hello"

	//ServiceKeyEcho echo
	ServiceKeyEcho = "echo"
)

//Config config definition
type Config struct {
	Mongo     o2m.MongoConfig   `mapstructure:"mongo"`
	Log       LogConfig         `mapstructure:"log"`
	Debug     bool              `mapstructure:"debug"`
	Endpoints map[string]string `mapstructure:"endpoints"`
	SignKey   string            `mapstructure:"sign_key"`
}

//LogConfig log config
type LogConfig struct {
	logDir      string `mapstructure:"log_dir"`
	logtostderr bool   `mapstructure:"logtostderr"`
	v           int    `mapstructure:"v"`
}

//DefaultConfig default config
func DefaultConfig() *Config {
	return &Config{
		Debug:   true,
		SignKey: "grpcapi-38ASD(*DFL@S",
		Endpoints: map[string]string{
			"echo":  "localhost:9001",
			"hello": "localhost:9002",
		},
		Mongo: o2m.MongoConfig{
			Addrs:     []string{"127.0.0.1:27017"},
			Database:  "oauth2",
			Username:  "oauth2",
			Password:  "oauth2",
			PoolLimit: 10,
		},
		Log: LogConfig{
			logtostderr: true,
			v:           -1,
		},
	}
}

//LoadConfig load default config
func LoadConfig() *Config {
	flag.Parse()
	return LoadConfigFile(*configFile)
}

//LoadConfigFile load config file
func LoadConfigFile(file string) *Config {
	viper.SetConfigFile(file)
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()
	if err != nil {
		glog.Info(err)
	}
	c := DefaultConfig()
	err = viper.Unmarshal(c)
	if err != nil {
		glog.Fatalf("unable to decode into struct, %v", err)
	}

	configLog(c)

	return c
}
func configLog(cfg *Config) {
	if cfg.Log.logDir != "" {
		flag.Set("log_dir", cfg.Log.logDir)
	}
	flag.Set("logtostderr", strconv.FormatBool(cfg.Log.logtostderr))
	if cfg.Log.v >= 0 {
		flag.Set("v", strconv.Itoa(cfg.Log.v))
	}

}
