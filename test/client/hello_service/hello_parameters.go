// Code generated by go-swagger; DO NOT EDIT.

package hello_service

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"
	"time"

	"golang.org/x/net/context"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"

	strfmt "github.com/go-openapi/strfmt"

	models "grpcapi/test/models"
)

// NewHelloParams creates a new HelloParams object
// with the default values initialized.
func NewHelloParams() *HelloParams {
	var ()
	return &HelloParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewHelloParamsWithTimeout creates a new HelloParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewHelloParamsWithTimeout(timeout time.Duration) *HelloParams {
	var ()
	return &HelloParams{

		timeout: timeout,
	}
}

// NewHelloParamsWithContext creates a new HelloParams object
// with the default values initialized, and the ability to set a context for a request
func NewHelloParamsWithContext(ctx context.Context) *HelloParams {
	var ()
	return &HelloParams{

		Context: ctx,
	}
}

// NewHelloParamsWithHTTPClient creates a new HelloParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewHelloParamsWithHTTPClient(client *http.Client) *HelloParams {
	var ()
	return &HelloParams{
		HTTPClient: client,
	}
}

/*HelloParams contains all the parameters to send to the API endpoint
for the hello operation typically these are written to a http.Request
*/
type HelloParams struct {

	/*Body*/
	Body *models.GrpcapiHelloRequest

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the hello params
func (o *HelloParams) WithTimeout(timeout time.Duration) *HelloParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the hello params
func (o *HelloParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the hello params
func (o *HelloParams) WithContext(ctx context.Context) *HelloParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the hello params
func (o *HelloParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the hello params
func (o *HelloParams) WithHTTPClient(client *http.Client) *HelloParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the hello params
func (o *HelloParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithBody adds the body to the hello params
func (o *HelloParams) WithBody(body *models.GrpcapiHelloRequest) *HelloParams {
	o.SetBody(body)
	return o
}

// SetBody adds the body to the hello params
func (o *HelloParams) SetBody(body *models.GrpcapiHelloRequest) {
	o.Body = body
}

// WriteToRequest writes these params to a swagger request
func (o *HelloParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if o.Body != nil {
		if err := r.SetBodyParam(o.Body); err != nil {
			return err
		}
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
