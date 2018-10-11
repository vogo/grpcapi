package config

import (
	"flag"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	flag.Parse()
	c := LoadConfigFile("config_test.yml")
	assert.Equal(t, "test", c.Mongo.Database)
	assert.Equal(t, "oauth2", c.Mongo.Username)
	assert.Equal(t, "oauth2", c.Mongo.Password)
	l := len(c.Mongo.Addrs)
	assert.Equal(t, 2, l)
	if l > 1 {
		assert.Equal(t, "127.0.0.1:27017", c.Mongo.Addrs[0])
		assert.Equal(t, "127.0.0.1:27018", c.Mongo.Addrs[1])
	}

	assert.True(t, c.Debug)
	assert.Equal(t, "localhost:9001", c.Endpoints["echo"])
	assert.Equal(t, "localhost:9002", c.Endpoints["hello"])
}
