package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	c := LoadConfigFile("config_test")
	assert.Equal(t, "test", c.Mongo.Database)
	assert.Equal(t, "oauth2", c.Mongo.Username)
	assert.Equal(t, "oauth2", c.Mongo.Password)
	l := len(c.Mongo.Addrs)
	assert.Equal(t, 2, l)
	if l > 1 {
		assert.Equal(t, "127.0.0.1:27017", c.Mongo.Addrs[0])
		assert.Equal(t, "127.0.0.1:27018", c.Mongo.Addrs[1])
	}
}
