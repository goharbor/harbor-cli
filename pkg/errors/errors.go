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

type Error struct {
	cause   error
	message string
	hint    string
}

func New(message string) *Error {
	return &Error{
		message: message,
	}
}

func (e *Error) WithCause(cause error) *Error {
	e.cause = cause
	return e
}

func (e *Error) WithHint(hint string) *Error {
	e.hint = hint
	return e
}

func (e *Error) Error() string {
	if e.cause != nil {
		return e.message + ": " + e.cause.Error()
	}
	return e.message
}

func (e *Error) Hint() string {
	return e.hint
}

func (e *Error) Unwrap() error {
	return e.cause
}
