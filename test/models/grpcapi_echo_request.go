// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/swag"
)

// GrpcapiEchoRequest grpcapi echo request
// swagger:model grpcapiEchoRequest
type GrpcapiEchoRequest struct {

	// value
	Value string `json:"value,omitempty"`
}

// Validate validates this grpcapi echo request
func (m *GrpcapiEchoRequest) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *GrpcapiEchoRequest) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *GrpcapiEchoRequest) UnmarshalBinary(b []byte) error {
	var res GrpcapiEchoRequest
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
