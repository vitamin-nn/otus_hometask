// Code generated by go-swagger; DO NOT EDIT.

package calendar_service

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"github.com/vitamin-nn/otus_hometask/hw12_13_14_15_calendar/tests/httpclient/models"
)

// CalendarServiceGetEventsDayReader is a Reader for the CalendarServiceGetEventsDay structure.
type CalendarServiceGetEventsDayReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *CalendarServiceGetEventsDayReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewCalendarServiceGetEventsDayOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewCalendarServiceGetEventsDayDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewCalendarServiceGetEventsDayOK creates a CalendarServiceGetEventsDayOK with default headers values
func NewCalendarServiceGetEventsDayOK() *CalendarServiceGetEventsDayOK {
	return &CalendarServiceGetEventsDayOK{}
}

/*CalendarServiceGetEventsDayOK handles this case with default header values.

A successful response.
*/
type CalendarServiceGetEventsDayOK struct {
	Payload *models.GetEventsResponse
}

func (o *CalendarServiceGetEventsDayOK) Error() string {
	return fmt.Sprintf("[POST /get_events_day][%d] calendarServiceGetEventsDayOK  %+v", 200, o.Payload)
}

func (o *CalendarServiceGetEventsDayOK) GetPayload() *models.GetEventsResponse {
	return o.Payload
}

func (o *CalendarServiceGetEventsDayOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.GetEventsResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewCalendarServiceGetEventsDayDefault creates a CalendarServiceGetEventsDayDefault with default headers values
func NewCalendarServiceGetEventsDayDefault(code int) *CalendarServiceGetEventsDayDefault {
	return &CalendarServiceGetEventsDayDefault{
		_statusCode: code,
	}
}

/*CalendarServiceGetEventsDayDefault handles this case with default header values.

An unexpected error response
*/
type CalendarServiceGetEventsDayDefault struct {
	_statusCode int

	Payload *models.RuntimeError
}

// Code gets the status code for the calendar service get events day default response
func (o *CalendarServiceGetEventsDayDefault) Code() int {
	return o._statusCode
}

func (o *CalendarServiceGetEventsDayDefault) Error() string {
	return fmt.Sprintf("[POST /get_events_day][%d] CalendarService_GetEventsDay default  %+v", o._statusCode, o.Payload)
}

func (o *CalendarServiceGetEventsDayDefault) GetPayload() *models.RuntimeError {
	return o.Payload
}

func (o *CalendarServiceGetEventsDayDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.RuntimeError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
