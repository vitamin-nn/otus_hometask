// Code generated by go-swagger; DO NOT EDIT.

package calendar_service

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"net/http"
	"time"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"

	"github.com/vitamin-nn/otus_hometask/hw12_13_14_15_calendar/tests/httpclient/models"
)

// NewCalendarServiceCreateEventParams creates a new CalendarServiceCreateEventParams object
// with the default values initialized.
func NewCalendarServiceCreateEventParams() *CalendarServiceCreateEventParams {
	var ()
	return &CalendarServiceCreateEventParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewCalendarServiceCreateEventParamsWithTimeout creates a new CalendarServiceCreateEventParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewCalendarServiceCreateEventParamsWithTimeout(timeout time.Duration) *CalendarServiceCreateEventParams {
	var ()
	return &CalendarServiceCreateEventParams{

		timeout: timeout,
	}
}

// NewCalendarServiceCreateEventParamsWithContext creates a new CalendarServiceCreateEventParams object
// with the default values initialized, and the ability to set a context for a request
func NewCalendarServiceCreateEventParamsWithContext(ctx context.Context) *CalendarServiceCreateEventParams {
	var ()
	return &CalendarServiceCreateEventParams{

		Context: ctx,
	}
}

// NewCalendarServiceCreateEventParamsWithHTTPClient creates a new CalendarServiceCreateEventParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewCalendarServiceCreateEventParamsWithHTTPClient(client *http.Client) *CalendarServiceCreateEventParams {
	var ()
	return &CalendarServiceCreateEventParams{
		HTTPClient: client,
	}
}

/*CalendarServiceCreateEventParams contains all the parameters to send to the API endpoint
for the calendar service create event operation typically these are written to a http.Request
*/
type CalendarServiceCreateEventParams struct {

	/*Body*/
	Body *models.CreateEventRequest

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the calendar service create event params
func (o *CalendarServiceCreateEventParams) WithTimeout(timeout time.Duration) *CalendarServiceCreateEventParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the calendar service create event params
func (o *CalendarServiceCreateEventParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the calendar service create event params
func (o *CalendarServiceCreateEventParams) WithContext(ctx context.Context) *CalendarServiceCreateEventParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the calendar service create event params
func (o *CalendarServiceCreateEventParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the calendar service create event params
func (o *CalendarServiceCreateEventParams) WithHTTPClient(client *http.Client) *CalendarServiceCreateEventParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the calendar service create event params
func (o *CalendarServiceCreateEventParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithBody adds the body to the calendar service create event params
func (o *CalendarServiceCreateEventParams) WithBody(body *models.CreateEventRequest) *CalendarServiceCreateEventParams {
	o.SetBody(body)
	return o
}

// SetBody adds the body to the calendar service create event params
func (o *CalendarServiceCreateEventParams) SetBody(body *models.CreateEventRequest) {
	o.Body = body
}

// WriteToRequest writes these params to a swagger request
func (o *CalendarServiceCreateEventParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

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
