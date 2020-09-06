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

// NewCalendarServiceGetEventsMonthParams creates a new CalendarServiceGetEventsMonthParams object
// with the default values initialized.
func NewCalendarServiceGetEventsMonthParams() *CalendarServiceGetEventsMonthParams {
	var ()
	return &CalendarServiceGetEventsMonthParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewCalendarServiceGetEventsMonthParamsWithTimeout creates a new CalendarServiceGetEventsMonthParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewCalendarServiceGetEventsMonthParamsWithTimeout(timeout time.Duration) *CalendarServiceGetEventsMonthParams {
	var ()
	return &CalendarServiceGetEventsMonthParams{

		timeout: timeout,
	}
}

// NewCalendarServiceGetEventsMonthParamsWithContext creates a new CalendarServiceGetEventsMonthParams object
// with the default values initialized, and the ability to set a context for a request
func NewCalendarServiceGetEventsMonthParamsWithContext(ctx context.Context) *CalendarServiceGetEventsMonthParams {
	var ()
	return &CalendarServiceGetEventsMonthParams{

		Context: ctx,
	}
}

// NewCalendarServiceGetEventsMonthParamsWithHTTPClient creates a new CalendarServiceGetEventsMonthParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewCalendarServiceGetEventsMonthParamsWithHTTPClient(client *http.Client) *CalendarServiceGetEventsMonthParams {
	var ()
	return &CalendarServiceGetEventsMonthParams{
		HTTPClient: client,
	}
}

/*CalendarServiceGetEventsMonthParams contains all the parameters to send to the API endpoint
for the calendar service get events month operation typically these are written to a http.Request
*/
type CalendarServiceGetEventsMonthParams struct {

	/*Body*/
	Body *models.GetEventsRequest

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the calendar service get events month params
func (o *CalendarServiceGetEventsMonthParams) WithTimeout(timeout time.Duration) *CalendarServiceGetEventsMonthParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the calendar service get events month params
func (o *CalendarServiceGetEventsMonthParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the calendar service get events month params
func (o *CalendarServiceGetEventsMonthParams) WithContext(ctx context.Context) *CalendarServiceGetEventsMonthParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the calendar service get events month params
func (o *CalendarServiceGetEventsMonthParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the calendar service get events month params
func (o *CalendarServiceGetEventsMonthParams) WithHTTPClient(client *http.Client) *CalendarServiceGetEventsMonthParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the calendar service get events month params
func (o *CalendarServiceGetEventsMonthParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithBody adds the body to the calendar service get events month params
func (o *CalendarServiceGetEventsMonthParams) WithBody(body *models.GetEventsRequest) *CalendarServiceGetEventsMonthParams {
	o.SetBody(body)
	return o
}

// SetBody adds the body to the calendar service get events month params
func (o *CalendarServiceGetEventsMonthParams) SetBody(body *models.GetEventsRequest) {
	o.Body = body
}

// WriteToRequest writes these params to a swagger request
func (o *CalendarServiceGetEventsMonthParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

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