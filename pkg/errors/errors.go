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

type Frame struct {
	Message string
	Hints   []string
}

type Error struct {
	frames []Frame
	cause  error
}

func New(message string) *Error {
	return &Error{
		frames: []Frame{{Message: message}},
	}
}

func Newf(format string, args ...any) *Error {
	return &Error{
		frames: []Frame{{Message: fmt.Sprintf(format, args...)}},
	}
}

func Wrap(err error, message string) *Error {
	e := AsError(err)
	e.frames = append([]Frame{{Message: message}}, e.frames...)
	return e
}

func Wrapf(err error, format string, args ...any) *Error {
	return Wrap(err, fmt.Sprintf(format, args...))
}

func (e *Error) WithHint(hint string) *Error {
	if len(e.frames) > 0 {
		e.frames[0].Hints = append(e.frames[0].Hints, hint)
	}
	return e
}

func (e *Error) WithCause(cause error) *Error {
	if e.cause == nil {
		e.cause = cause
	}
	return e
}

func (e *Error) Error() string {
	if len(e.frames) == 0 {
		return ""
	}
	return e.frames[len(e.frames)-1].Message
}

func (e *Error) Message() string {
	if len(e.frames) == 0 {
		return ""
	}
	return e.frames[0].Message
}

func (e *Error) Errors() []string {
	msgs := make([]string, len(e.frames))
	for i, f := range e.frames {
		msgs[i] = f.Message
	}
	return msgs
}

func (e *Error) Hints() []string {
	var all []string
	for _, f := range e.frames {
		all = append(all, f.Hints...)
	}
	return all
}

func (e *Error) Frames() []Frame {
	return e.frames
}

func (e *Error) Cause() error { return e.cause }

func (e *Error) Status() string {
	if e.cause == nil {
		return ""
	}
	return parseHarborErrorCode(e.cause)
}

func (e *Error) Unwrap() error { return e.cause }

func AsError(err error) *Error {
	var e *Error
	if errors.As(err, &e) {
		return e
	}
	return &Error{
		frames: []Frame{{Message: parseHarborErrorMsg(err)}},
		cause:  err,
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
