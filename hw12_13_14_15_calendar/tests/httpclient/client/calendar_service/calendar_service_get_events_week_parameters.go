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

// NewCalendarServiceGetEventsWeekParams creates a new CalendarServiceGetEventsWeekParams object
// with the default values initialized.
func NewCalendarServiceGetEventsWeekParams() *CalendarServiceGetEventsWeekParams {
	var ()
	return &CalendarServiceGetEventsWeekParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewCalendarServiceGetEventsWeekParamsWithTimeout creates a new CalendarServiceGetEventsWeekParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewCalendarServiceGetEventsWeekParamsWithTimeout(timeout time.Duration) *CalendarServiceGetEventsWeekParams {
	var ()
	return &CalendarServiceGetEventsWeekParams{

		timeout: timeout,
	}
}

// NewCalendarServiceGetEventsWeekParamsWithContext creates a new CalendarServiceGetEventsWeekParams object
// with the default values initialized, and the ability to set a context for a request
func NewCalendarServiceGetEventsWeekParamsWithContext(ctx context.Context) *CalendarServiceGetEventsWeekParams {
	var ()
	return &CalendarServiceGetEventsWeekParams{

		Context: ctx,
	}
}

// NewCalendarServiceGetEventsWeekParamsWithHTTPClient creates a new CalendarServiceGetEventsWeekParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewCalendarServiceGetEventsWeekParamsWithHTTPClient(client *http.Client) *CalendarServiceGetEventsWeekParams {
	var ()
	return &CalendarServiceGetEventsWeekParams{
		HTTPClient: client,
	}
}

/*CalendarServiceGetEventsWeekParams contains all the parameters to send to the API endpoint
for the calendar service get events week operation typically these are written to a http.Request
*/
type CalendarServiceGetEventsWeekParams struct {

	/*Body*/
	Body *models.GetEventsRequest

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the calendar service get events week params
func (o *CalendarServiceGetEventsWeekParams) WithTimeout(timeout time.Duration) *CalendarServiceGetEventsWeekParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the calendar service get events week params
func (o *CalendarServiceGetEventsWeekParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the calendar service get events week params
func (o *CalendarServiceGetEventsWeekParams) WithContext(ctx context.Context) *CalendarServiceGetEventsWeekParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the calendar service get events week params
func (o *CalendarServiceGetEventsWeekParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the calendar service get events week params
func (o *CalendarServiceGetEventsWeekParams) WithHTTPClient(client *http.Client) *CalendarServiceGetEventsWeekParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the calendar service get events week params
func (o *CalendarServiceGetEventsWeekParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithBody adds the body to the calendar service get events week params
func (o *CalendarServiceGetEventsWeekParams) WithBody(body *models.GetEventsRequest) *CalendarServiceGetEventsWeekParams {
	o.SetBody(body)
	return o
}

// SetBody adds the body to the calendar service get events week params
func (o *CalendarServiceGetEventsWeekParams) SetBody(body *models.GetEventsRequest) {
	o.Body = body
}

// WriteToRequest writes these params to a swagger request
func (o *CalendarServiceGetEventsWeekParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

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
