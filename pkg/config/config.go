package config

import (
	"flag"
	"strconv"

	"github.com/go2s/o2m"
	"github.com/golang/glog"
	"github.com/spf13/viper"
)

var (
	//SecretKey secret key
	SecretKey = "ADkajfdJIALDSFJJkadf"

	//APIGatewayAddress api gateway address
	APIGatewayAddress = ":8080"

	//TokenServiceAddress token service address
	TokenServiceAddress = ":9001"

	//EchoServiceAddress echo service address
	EchoServiceAddress = ":9002"

	//HelloServiceAddress hello service address
	HelloServiceAddress = ":9003"

	//EchoServiceEndpoint echo service endpoint
	EchoServiceEndpoint = flag.String("echo_endpoint", EchoServiceAddress, "endpoint of echo service")

	//HelloServiceEndpoint hello service endpoint
	HelloServiceEndpoint = flag.String("hello_endpoint", HelloServiceAddress, "endpoint of hello service")
)

//Config config definition
type Config struct {
	Mongo o2m.MongoConfig `mapstructure:"mongo"`
	Log   LogConfig       `mapstructure:"log"`
	Debug bool            `mapstructure:"debug"`
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
		Debug: true,
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
	return LoadConfigFile("config")
}

//LoadConfigFile load config file
func LoadConfigFile(filename string) *Config {
	flag.Parse()

	viper.SetConfigName(filename)
	viper.AddConfigPath("/etc/grpcapi/")
	viper.AddConfigPath("$HOME/.grpcapi")
	viper.AddConfigPath(".")
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
