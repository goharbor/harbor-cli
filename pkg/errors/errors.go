// Copyright Project Harbor Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package errors

import (
	"errors"
	"fmt"
)

var (
	as = errors.As
)

type Error struct {
	cause   error
	message string
	hints   []string
}

func New(message string) *Error {
	return &Error{
		message: message,
	}
}

func Newf(format string, args ...any) *Error {
	return &Error{
		message: fmt.Sprintf(format, args...),
	}
}

func AsError(err error) *Error {
	var e *Error
	if errors.As(err, &e) {
		return e
	}
	return &Error{
		message: err.Error(),
		cause:   err,
	}
}

func IsError(err error) bool {
	var e *Error
	return as(err, &e)
}

func Cause(err error) error {
	if e := isHarborError(err); e != nil {
		return e.Cause()
	}
	return nil
}

func Hints(err error) []string {
	if e := isHarborError(err); e != nil {
		return e.Hints()
	}
	return nil
}

func Status(err error) string {
	if e := isHarborError(err); e != nil {
		return e.Status()
	}
	return ""
}

func (e *Error) WithCause(cause error) *Error {
	if e.cause == nil {
		e.cause = cause
	}
	return e
}

func (e *Error) WithHint(hint string) *Error {
	e.hints = append(e.hints, hint)
	return e
}

func (e *Error) Error() string {
	if e.cause != nil {
		causeMsg := parseHarborErrorMsg(e.cause)
		return e.message + ": " + causeMsg
	}
	return e.message
}

func (e *Error) Cause() error { return e.cause }

func (e *Error) Status() string { return parseHarborErrorCode(e.cause) }

func (e *Error) Hints() []string { return e.hints }

func (e *Error) Message() string { return e.message }

func (e *Error) Unwrap() error { return e.cause }
