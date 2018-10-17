// Copyright 2018 The Vogo Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package jsonutil

import (
	"encoding/json"

	"github.com/vogo/clog"
)

//ToString json to string
func ToString(o interface{}) string {
	b, err := json.Marshal(o)
	if err != nil {
		clog.Error(nil, "Failed to encode [%+v], error: %+v", o, err)
		return ""
	}
	return string(b)
}

//Parse json to object
func Parse(s string, o interface{}) error {
	return json.Unmarshal([]byte(s), o)
}
