// Copyright 2018 The Vogo Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package identity

import "github.com/vogo/grpcapi/pkg/util/jsonutil"

//Identity authorization identity
type Identity struct {
	UserID string   `json:"id"`
	Roles  []string `json:"rol,omitempty"`
	Scopes []string `json:"scp,omitempty"`
}

//ToJSON convert to json string
func (id *Identity) ToJSON() string {
	return jsonutil.ToString(id)
}

//Parse string to identity
func Parse(json string) (*Identity, error) {
	id := &Identity{}
	if err := jsonutil.Parse(json, id); err != nil {
		return nil, err
	}
	return id, nil
}

//New identity
func New(id string, roles, scopes []string) *Identity {
	return &Identity{
		UserID: id,
		Roles:  roles,
		Scopes: scopes,
	}
}
