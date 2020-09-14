// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// ModifyEventRequest modify event request
//
// swagger:model ModifyEventRequest
type ModifyEventRequest struct {

	// description
	Description string `json:"description,omitempty"`

	// end at
	// Format: date-time
	EndAt strfmt.DateTime `json:"end_at,omitempty"`

	// notify at
	// Format: date-time
	NotifyAt strfmt.DateTime `json:"notify_at,omitempty"`

	// start at
	// Format: date-time
	StartAt strfmt.DateTime `json:"start_at,omitempty"`

	// title
	Title string `json:"title,omitempty"`
}

// Validate validates this modify event request
func (m *ModifyEventRequest) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateEndAt(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateNotifyAt(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateStartAt(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *ModifyEventRequest) validateEndAt(formats strfmt.Registry) error {

	if swag.IsZero(m.EndAt) { // not required
		return nil
	}

	if err := validate.FormatOf("end_at", "body", "date-time", m.EndAt.String(), formats); err != nil {
		return err
	}

	return nil
}

func (m *ModifyEventRequest) validateNotifyAt(formats strfmt.Registry) error {

	if swag.IsZero(m.NotifyAt) { // not required
		return nil
	}

	if err := validate.FormatOf("notify_at", "body", "date-time", m.NotifyAt.String(), formats); err != nil {
		return err
	}

	return nil
}

func (m *ModifyEventRequest) validateStartAt(formats strfmt.Registry) error {

	if swag.IsZero(m.StartAt) { // not required
		return nil
	}

	if err := validate.FormatOf("start_at", "body", "date-time", m.StartAt.String(), formats); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *ModifyEventRequest) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ModifyEventRequest) UnmarshalBinary(b []byte) error {
	var res ModifyEventRequest
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
