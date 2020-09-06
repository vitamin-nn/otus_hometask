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

// CalendarServiceDeleteEventReader is a Reader for the CalendarServiceDeleteEvent structure.
type CalendarServiceDeleteEventReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *CalendarServiceDeleteEventReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewCalendarServiceDeleteEventOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewCalendarServiceDeleteEventDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewCalendarServiceDeleteEventOK creates a CalendarServiceDeleteEventOK with default headers values
func NewCalendarServiceDeleteEventOK() *CalendarServiceDeleteEventOK {
	return &CalendarServiceDeleteEventOK{}
}

/*CalendarServiceDeleteEventOK handles this case with default header values.

A successful response.
*/
type CalendarServiceDeleteEventOK struct {
	Payload *models.DeleteResponse
}

func (o *CalendarServiceDeleteEventOK) Error() string {
	return fmt.Sprintf("[DELETE /delete/{event_id}][%d] calendarServiceDeleteEventOK  %+v", 200, o.Payload)
}

func (o *CalendarServiceDeleteEventOK) GetPayload() *models.DeleteResponse {
	return o.Payload
}

func (o *CalendarServiceDeleteEventOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.DeleteResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewCalendarServiceDeleteEventDefault creates a CalendarServiceDeleteEventDefault with default headers values
func NewCalendarServiceDeleteEventDefault(code int) *CalendarServiceDeleteEventDefault {
	return &CalendarServiceDeleteEventDefault{
		_statusCode: code,
	}
}

/*CalendarServiceDeleteEventDefault handles this case with default header values.

An unexpected error response
*/
type CalendarServiceDeleteEventDefault struct {
	_statusCode int

	Payload *models.RuntimeError
}

// Code gets the status code for the calendar service delete event default response
func (o *CalendarServiceDeleteEventDefault) Code() int {
	return o._statusCode
}

func (o *CalendarServiceDeleteEventDefault) Error() string {
	return fmt.Sprintf("[DELETE /delete/{event_id}][%d] CalendarService_DeleteEvent default  %+v", o._statusCode, o.Payload)
}

func (o *CalendarServiceDeleteEventDefault) GetPayload() *models.RuntimeError {
	return o.Payload
}

func (o *CalendarServiceDeleteEventDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.RuntimeError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}