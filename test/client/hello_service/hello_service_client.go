// Code generated by go-swagger; DO NOT EDIT.

package hello_service

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"
)

// New creates a new hello service API client.
func New(transport runtime.ClientTransport, formats strfmt.Registry) *Client {
	return &Client{transport: transport, formats: formats}
}

/*
Client for hello service API
*/
type Client struct {
	transport runtime.ClientTransport
	formats   strfmt.Registry
}

/*
Hello hello API
*/
func (a *Client) Hello(params *HelloParams, authInfo runtime.ClientAuthInfoWriter) (*HelloOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewHelloParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "Hello",
		Method:             "POST",
		PathPattern:        "/api/v1/hello",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http", "https"},
		Params:             params,
		Reader:             &HelloReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	return result.(*HelloOK), nil

}

// SetTransport changes the transport on the client
func (a *Client) SetTransport(transport runtime.ClientTransport) {
	a.transport = transport
}
