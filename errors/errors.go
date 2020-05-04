package errors

import (
	"errors"
	"fmt"
	"github.com/uhppoted/uhppoted-api/uhppoted"
	"net/http"
)

type IError struct {
	Err      error  `json:"-"`
	DeviceID uint32 `json:"-"`
	Tag      string `json:"-"`
	Code     int    `json:"error-code"`
	Message  string `json:"message"`
}

// var (
// 	InvalidDeviceID    = ferror(fmt.Errorf("%w: Missing device ID", uhppoted.BadRequest), "Missing device ID")
// 	InvalidCardNumber  = ferror(fmt.Errorf("%w: Missing/invalid card number", uhppoted.BadRequest), "Missing/invalid card number")
// 	InvalidDoorID      = ferror(fmt.Errorf("%w: Missing/invalid door ID", uhppoted.BadRequest), "Missing/invalid door ID")
// 	InvalidDoorDelay   = ferror(fmt.Errorf("%w: Missing/invalid door delay", uhppoted.BadRequest), "Missing/invalid door delay")
// 	InvalidDoorControl = ferror(fmt.Errorf("%w: Missing/invalid door control", uhppoted.BadRequest), "Missing/invalid door control")
// 	InvalidEventID     = ferror(fmt.Errorf("%w: Missing/invalid event ID", uhppoted.BadRequest), "Missing/invalid event ID")
// 	InvalidDateTime    = ferror(fmt.Errorf("%w: Missing/invalid date/time", uhppoted.BadRequest), "Missing/invalid date/time")
// )

func (e *IError) Error() string {
	return fmt.Sprintf("%v", e.Err)
}

func Errorf(err error, deviceID uint32, tag string, msg string) *IError {
	status := http.StatusInternalServerError

	if errors.Is(err, uhppoted.InternalServerError) {
		status = http.StatusInternalServerError
	} else if errors.Is(err, uhppoted.NotFound) {
		status = http.StatusNotFound
	} else if errors.Is(err, uhppoted.BadRequest) {
		status = http.StatusBadRequest
	}

	return &IError{
		Err:      err,
		DeviceID: deviceID,
		Tag:      tag,
		Code:     status,
		Message:  msg,
	}
}
