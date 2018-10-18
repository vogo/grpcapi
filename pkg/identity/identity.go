// Copyright 2018 The Vogo Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

// Benchmark result:
// ------------------
// BenchmarkToString-4      	 3000000	       424 ns/op
// BenchmarkToJSON-4        	 2000000	       743 ns/op
// BenchmarkParseString-4   	10000000	       214 ns/op
// BenchmarkParseJSON-4     	 1000000	      2277 ns/op
//
// string encoding has better performance

package identity

import (
	"bytes"
	"errors"
	"strings"

	"github.com/vogo/grpcapi/pkg/util/jsonutil"
)

// FieldSpliter field splitor
const FieldSpliter = '|'

//ItemSpliter item splitor
const ItemSpliter = ','

//Identity authorization identity
type Identity struct {
	UserID string   `json:"id"`
	Roles  []string `json:"rol,omitempty"`
	Scopes []string `json:"scp,omitempty"`
}

//JSON convert to json
func (id *Identity) JSON() string {
	return jsonutil.ToString(id)
}

//String convert to string
func (id *Identity) String() string {
	buf := bytes.NewBufferString(id.UserID)

	if len(id.Roles) == 0 && len(id.Scopes) == 0 {
		return buf.String()
	}

	buf.WriteByte(FieldSpliter)
	more := false
	for _, item := range id.Roles {
		if more {
			buf.WriteByte(ItemSpliter)
		}
		buf.WriteString(item)
		more = true
	}

	if len(id.Scopes) > 0 {
		buf.WriteByte(FieldSpliter)
		more = false
		for _, item := range id.Scopes {
			if more {
				buf.WriteByte(ItemSpliter)
			}
			buf.WriteString(item)
			more = true
		}
	}

	return buf.String()
}

//Parse string to identity
func Parse(s string) (*Identity, error) {
	if s == "" {
		return nil, errors.New("nil identity")
	}

	id := &Identity{}

	ParseString(id, s)

	return id, nil
}

//ParseString parse string
func ParseString(id *Identity, s string) {
	index := strings.IndexByte(s, FieldSpliter)
	if index < 0 {
		id.UserID = s
		return
	}
	id.UserID = s[:index]

	s = s[index+1:]
	index = strings.IndexByte(s, FieldSpliter)
	arrString := s
	if index >= 0 {
		arrString = s[:index]
	}
	if arrString != "" {
		id.Roles = strings.Split(arrString, string(ItemSpliter))
	}
	if index < 0 {
		return
	}

	arrString = s[index+1:]
	if arrString != "" {
		id.Scopes = strings.Split(arrString, string(ItemSpliter))
	}
	return
}

//ParseJSON parse json
func ParseJSON(id *Identity, j string) error {
	return jsonutil.Parse(j, id)
}

//New identity
func New(id string, roles, scopes []string) *Identity {
	return &Identity{
		UserID: id,
		Roles:  roles,
		Scopes: scopes,
	}
}
