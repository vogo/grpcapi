// Copyright 2018 The Vogo Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package config

import (
	"flag"
	"fmt"
	"strconv"

	"github.com/go2s/o2m"
	"github.com/golang/glog"
	"github.com/spf13/viper"
)

var (
	configFile = flag.String("config", "config", "config file path")
)

const (
	//HostPrefix host name prefix
	HostPrefix = "grpc-"
	//ServiceEcho echo service name
	ServiceEcho = "echo"
	//ServiceHello hello service name
	ServiceHello = "hello"
)

const (
	//PortAPIGateway api gateway port
	PortAPIGateway = 8080
	//PortEcho echo service port
	PortEcho = 9001
	//PortHello hello service port
	PortHello = 9002
)

const (
	//HostEcho echo service host
	HostEcho = HostPrefix + ServiceEcho
	//HostHello hello service host
	HostHello = HostPrefix + ServiceHello
)

var (
	//EndpointEcho echo service address
	EndpointEcho = fmt.Sprintf("%s:%d", HostEcho, PortEcho)
	//EndpointHello hello service address
	EndpointHello = fmt.Sprintf("%s:%d", HostHello, PortHello)
)

//Config config definition
type Config struct {
	Mongo   o2m.MongoConfig `mapstructure:"mongo"`
	Log     LogConfig       `mapstructure:"log"`
	Debug   bool            `mapstructure:"debug"`
	SignKey string          `mapstructure:"sign-key"`
}

//LogConfig log config
type LogConfig struct {
	logDir      string `mapstructure:"log-dir"`
	logtostderr bool   `mapstructure:"logtostderr"`
	v           int    `mapstructure:"v"`
}

//DefaultConfig default config
func DefaultConfig() *Config {
	return &Config{
		Debug:   true,
		SignKey: "grpcapi-38ASD(*DFL@S",
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
