// Copyright 2018 The Vogo Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package main

import (
	"grpcapi/pkg/apigateway"
	"grpcapi/pkg/config"
)

func main() {
	cfg := config.LoadConfig()
	apigateway.Serve(cfg)
}
