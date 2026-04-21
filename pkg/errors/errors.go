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
	"strings"

	"github.com/charmbracelet/lipgloss/tree"
	"github.com/goharbor/harbor-cli/pkg/views"
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

func New(message string, hints ...string) *Error {
	return &Error{
		frames: []Frame{{Message: message, Hints: hints}},
	}
}

func Newf(format string, args ...any) *Error {
	return &Error{
		frames: []Frame{{Message: fmt.Sprintf(format, args...)}},
	}
}

func NewWithCause(cause error, message string, hints ...string) *Error {
	return &Error{
		cause:  cause,
		frames: []Frame{{Message: message, Hints: hints}},
	}
}

func (e *Error) WithMessage(message string, hints ...string) *Error {
	e.frames = append(e.frames, Frame{Message: message, Hints: hints})
	return e
}

func (e *Error) Error() string {
	if len(e.frames) == 0 {
		return ""
	}

	var parts []string

	rootFrame := e.frames[0]
	parts = append(parts, views.ErrCauseStyle.Render(rootFrame.Message))

	if e.cause != nil {
		if code := parseHarborErrorCode(e.cause); code != "" {
			parts = append(parts,
				views.ErrTitleStyle.Render("Code: ")+views.ErrCauseStyle.Render(code),
			)
		}
	}

	if e.cause != nil {
		causeText := e.cause.Error()
		if he := isHarborError(e.cause); he != nil {
			causeText = he.Message()
		}

		cause := views.ErrTitleStyle.Render("Cause: ")
		causeText = views.ErrCauseStyle.Render(causeText)
		causeTree := tree.Root(cause + causeText).
			Enumerator(tree.RoundedEnumerator).
			EnumeratorStyle(views.ErrEnumeratorStyle).
			ItemStyle(views.ErrHintStyle)

		if he := isHarborError(e.cause); he != nil {
			for _, h := range he.Hints() {
				causeTree.Child(h)
			}
		}
		parts = append(parts, causeTree.String())
	}

	if len(rootFrame.Hints) > 0 {
		hintsTree := tree.New().
			Root("Hints:").
			RootStyle(views.ErrTitleStyle).
			Enumerator(tree.RoundedEnumerator).
			EnumeratorStyle(views.ErrEnumeratorStyle).
			ItemStyle(views.ErrHintStyle)

		for _, h := range rootFrame.Hints {
			hintsTree.Child(h)
		}
		parts = append(parts, hintsTree.String())
	}

	if len(e.frames) > 1 {
		msgTree := tree.New().
			Root("Messages:").
			RootStyle(views.ErrTitleStyle).
			Enumerator(tree.RoundedEnumerator).
			EnumeratorStyle(views.ErrEnumeratorStyle).
			ItemStyle(views.ErrTitleStyle)

		for _, f := range e.frames[1:] {
			msgWithHints := tree.Root(f.Message).
				RootStyle(views.ErrTitleStyle).
				Enumerator(tree.RoundedEnumerator).
				EnumeratorStyle(views.ErrEnumeratorStyle).
				ItemStyle(views.ErrHintStyle)
			for _, h := range f.Hints {
				msgWithHints.Child(h)
			}
			msgTree.Child(msgWithHints)
		}
		parts = append(parts, msgTree.String())
	}

	return strings.Join(parts, "\n")
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
