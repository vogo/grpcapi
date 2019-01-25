// Copyright 2018 The Vogo Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

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
	assert.Equal(t, "mongodb://127.0.0.1:27017,127.0.0.1:27018", c.Mongo.Address)
}
