package ippanel

import (
	"encoding/json"
	"errors"
	"fmt"
)

// ResponseCode api response code error type
type ResponseCode int

const (
	ErrForbidden           ResponseCode = 403
	ErrNotFound            ResponseCode = 404
	ErrUnprocessableEntity ResponseCode = 422
	ErrInternalServer      ResponseCode = 500
)

// Error general service error
type Error struct {
	Code    ResponseCode
	Message interface{}
}

// FieldErrs input field level errors
type FieldErrs map[string][]string

// Error implement error interface
func (e Error) Error() string {
	switch e.Message.(type) {
	case string:
		return e.Message.(string)
	case FieldErrs:
		m, _ := json.Marshal(e.Message)
		return string(m)
	}

	return fmt.Sprint(e.Code)
}

// fieldErrsRes field errors response
type fieldErrsRes struct {
	Errors FieldErrs `json:"error"`
}

// defaultErrsRes default template for errors body
type defaultErrsRes struct {
	Errors string `json:"error"`
}

// ParseErrors ...
func ParseErrors(res *BaseResponse) error {
	var err error
	e := Error{Code: res.Code}

	messageFieldErrs := fieldErrsRes{}
	if err = json.Unmarshal(res.Data, &messageFieldErrs); err == nil {
		e.Message = messageFieldErrs.Errors
	} else {
		messageDefaultErrs := defaultErrsRes{}
		if err = json.Unmarshal(res.Data, &messageDefaultErrs); err == nil {
			e.Message = messageDefaultErrs.Errors
		}
	}

	if err != nil {
		return errors.New("cant marshal errors into standard template")
	}

	return e
}
