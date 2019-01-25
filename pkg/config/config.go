// Copyright 2018 The Vogo Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package config

import (
	"context"
	"flag"
	"fmt"
	"path/filepath"
	"strconv"

	"github.com/go2s/o2s/o2m"
	"github.com/spf13/viper"
	"github.com/vogo/clog"
	"github.com/vogo/grpcapi/pkg/util/ctxutil"
	"gopkg.in/natefinch/lumberjack.v2"
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

	//ServiceMongodb mongodb service name
	ServiceMongodb = "mongodb"
)

const (
	//PortAPIGateway api gateway port
	PortAPIGateway = 8080

	//PortEcho echo service port
	PortEcho = 9001

	//PortHello hello service port
	PortHello = 9002

	//PortMongodb mongodb service port
	PortMongodb = 27017
)

const (
	//HostEcho echo service host
	HostEcho = HostPrefix + ServiceEcho

	//HostHello hello service host
	HostHello = HostPrefix + ServiceHello

	//HostMongodb mongodb service host
	HostMongodb = HostPrefix + ServiceMongodb
)

var (
	//EndpointEcho echo service address
	EndpointEcho = fmt.Sprintf("%s:%d", HostEcho, PortEcho)

	//EndpointHello hello service address
	EndpointHello = fmt.Sprintf("%s:%d", HostHello, PortHello)

	//EndpointMongodb mongodb service address
	EndpointMongodb = fmt.Sprintf("%s:%d", HostMongodb, PortMongodb)
)

//Config config definition
type Config struct {
	Mongo   o2m.MongoConfig `mapstructure:"mongo"`
	Log     LogConfig       `mapstructure:"log"`
	SignKey string          `mapstructure:"sign-key"`
}

//LogConfig log config
type LogConfig struct {
	Dir         string `mapstructure:"dir"`
	Filename    string `mapstructure:"filename"`
	Level       string `mapstructure:"level"`
	MaxSize     int    `mapstructure:"max-size"` //MB
	MaxBackups  int    `mapstructure:"max-backups"`
	MaxAge      int    `mapstructure:"max-age"` //days
	Compress    bool   `mapstructure:"compress"`
	Logtostderr bool   `mapstructure:"logtostderr"`
	V           int    `mapstructure:"v"`
}

//DefaultConfig default config
func DefaultConfig() *Config {
	return &Config{
		SignKey: "grpcapi-38ASD(*DFL@S",
		Mongo: o2m.MongoConfig{
			Address:   fmt.Sprintf("mongodb://%s", EndpointMongodb),
			Database:  "oauth2",
			Username:  "oauth2",
			Password:  "oauth2",
			PoolLimit: 10,
		},
		Log: LogConfig{
			Logtostderr: true,
			V:           -1,
			MaxSize:     10,
			MaxBackups:  10,
			MaxAge:      30,
			Compress:    false,
			Level:       "info",
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
	cfg := DefaultConfig()

	viper.SetConfigFile(file)
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()
	if err != nil {
		clog.Info(nil, "no config file provided, using default config instead. [error] %v", err)
	} else {
		err = viper.Unmarshal(cfg)
		if err != nil {
			clog.Fatal(nil, "unable to decode into struct, %v", err)
		}
	}

	configLog(cfg)

	return cfg
}
func configLog(cfg *Config) {
	if cfg.Log.Dir != "" {
		flag.Set("log_dir", cfg.Log.Dir)
		filename := cfg.Log.Filename
		if filename == "" {
			filename = "grpc.log"
		}
		clog.SetOutput(&lumberjack.Logger{
			Filename:   filepath.Join(cfg.Log.Dir, filename),
			MaxSize:    cfg.Log.MaxSize,
			MaxBackups: cfg.Log.MaxBackups,
			MaxAge:     cfg.Log.MaxAge,
			Compress:   cfg.Log.Compress})
	}

	clog.SetContextFommatter(func(ctx context.Context) string {
		if rid := ctxutil.GetRequestID(ctx); rid != "" {
			return rid
		}
		return "-"
	})

	clog.SetLevelByString(cfg.Log.Level)

	flag.Set("logtostderr", strconv.FormatBool(cfg.Log.Logtostderr))
	if cfg.Log.V >= 0 {
		flag.Set("v", strconv.Itoa(cfg.Log.V))
	}

}
