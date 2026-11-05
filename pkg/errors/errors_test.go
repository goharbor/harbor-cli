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
package errors_test

import (
	"errors"
	"fmt"
	"testing"

	harborerr "github.com/goharbor/harbor-cli/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew_MessageOnly(t *testing.T) {
	err := harborerr.New("something went wrong")
	require.NotNil(t, err)

	output := err.Error()
	assert.Contains(t, output, "something went wrong")
	assert.Equal(t, "something went wrong", err.Message())
	assert.Equal(t, []string{"something went wrong"}, err.Errors())
	assert.Empty(t, err.Hints())
	assert.Nil(t, err.Cause())
}

func TestNew_MessageWithHints(t *testing.T) {
	err := harborerr.New("auth failed", "check your credentials", "ensure token is not expired")
	require.NotNil(t, err)

	output := err.Error()
	assert.Contains(t, output, "auth failed")
	assert.Contains(t, output, "check your credentials")
	assert.Contains(t, output, "ensure token is not expired")

	assert.Equal(t, "auth failed", err.Message())
	assert.Equal(t, []string{"check your credentials", "ensure token is not expired"}, err.Hints())
	assert.Nil(t, err.Cause())
}

func TestNewf_FormatsMessage(t *testing.T) {
	err := harborerr.Newf("resource %q not found", "project-x")
	require.NotNil(t, err)

	output := err.Error()
	assert.Contains(t, output, `resource "project-x" not found`)
	assert.Equal(t, `resource "project-x" not found`, err.Message())
	assert.Empty(t, err.Hints())
}

func TestNewWithCause_Basic(t *testing.T) {
	cause := errors.New("connection refused")
	err := harborerr.NewWithCause(cause, "could not reach registry")

	output := err.Error()
	assert.Contains(t, output, "could not reach registry")
	assert.Contains(t, output, "Cause:")
	assert.Contains(t, output, "connection refused")
	assert.NotContains(t, output, "Messages:")

	assert.Equal(t, cause, err.Cause())
	assert.Equal(t, "could not reach registry", err.Message())
}

func TestNewWithCause_WithHints(t *testing.T) {
	cause := errors.New("401 Unauthorized")
	err := harborerr.NewWithCause(cause, "authentication failed",
		"ensure your credentials are correct",
		"run `harbor login` to re-authenticate",
	)

	output := err.Error()
	assert.Contains(t, output, "authentication failed")
	assert.Contains(t, output, "ensure your credentials are correct")
	assert.Contains(t, output, "run `harbor login` to re-authenticate")
	assert.Contains(t, output, "Cause:")
	assert.Contains(t, output, "401 Unauthorized")
	assert.NotContains(t, output, "Messages:")

	assert.Equal(t, cause, err.Cause())
	assert.Equal(t, []string{
		"ensure your credentials are correct",
		"run `harbor login` to re-authenticate",
	}, err.Hints())
}

func TestNewWithCause_CauseIsHarborError_ShowsCauseHintsInHeader(t *testing.T) {
	inner := harborerr.New("db connection lost", "check that the database is running")
	outer := harborerr.NewWithCause(inner, "repository unavailable")

	output := outer.Error()
	assert.Contains(t, output, "repository unavailable")
	assert.Contains(t, output, "check that the database is running")
	assert.NotContains(t, output, "Messages:")
}

func TestWithMessage_AppendsFrame(t *testing.T) {
	err := harborerr.NewWithCause(
		errors.New("dial tcp 127.0.0.1:5432: connect: connection refused"),
		"repository unavailable", "check that the database is running",
	).WithMessage("failed to delete artifact",
		"retry after resolving the underlying issue",
		"use --force to skip confirmation prompts",
	)

	output := err.Error()
	assert.Contains(t, output, "repository unavailable")
	assert.Contains(t, output, "check that the database is running")
	assert.Contains(t, output, "Cause:")
	assert.Contains(t, output, "connection refused")
	assert.Contains(t, output, "Messages:")
	assert.Contains(t, output, "failed to delete artifact")
	assert.Contains(t, output, "retry after resolving the underlying issue")
	assert.Contains(t, output, "use --force to skip confirmation prompts")

	assert.Equal(t, "repository unavailable", err.Message())
	assert.Len(t, err.Frames(), 2)
}

func TestWithMessage_MultipleFrames(t *testing.T) {
	err := harborerr.New("step 1").
		WithMessage("step 2", "hint A").
		WithMessage("step 3", "hint B", "hint C")

	assert.Len(t, err.Frames(), 3)
	assert.Equal(t, "step 1", err.Frames()[0].Message)
	assert.Equal(t, "step 2", err.Frames()[1].Message)
	assert.Equal(t, []string{"hint A"}, err.Frames()[1].Hints)
	assert.Equal(t, "step 3", err.Frames()[2].Message)
	assert.Equal(t, []string{"hint B", "hint C"}, err.Frames()[2].Hints)
}

func TestWithMessage_NoCause_NoRootHeader(t *testing.T) {
	err := harborerr.New("first").WithMessage("second", "a hint")

	output := err.Error()
	assert.NotContains(t, output, "Root:")
	assert.NotContains(t, output, "Code:")
	assert.Contains(t, output, "first")
	assert.Contains(t, output, "second")
	assert.Contains(t, output, "a hint")
}

func TestError_EmptyFrames(t *testing.T) {
	err := &harborerr.Error{}
	assert.Equal(t, "", err.Error())
}

func TestError_NoCause_SingleFrameNoHints(t *testing.T) {
	err := harborerr.New("operation not supported")
	output := err.Error()

	assert.NotContains(t, output, "Root:")
	assert.NotContains(t, output, "Code:")
	assert.Contains(t, output, "operation not supported")
}

func TestError_NoCause_SingleFrameWithHints(t *testing.T) {
	err := harborerr.New("plugin system not implemented",
		"this command is a placeholder for future plugin management",
	)

	output := err.Error()
	assert.NotContains(t, output, "Root:")
	assert.Contains(t, output, "plugin system not implemented")
	assert.Contains(t, output, "this command is a placeholder for future plugin management")
}

func TestMessage_ReturnsFirstFrame(t *testing.T) {
	err := harborerr.New("first").WithMessage("second")
	assert.Equal(t, "first", err.Message())
}

func TestErrors_AllMessages(t *testing.T) {
	err := harborerr.New("a").WithMessage("b").WithMessage("c")
	assert.Equal(t, []string{"a", "b", "c"}, err.Errors())
}

func TestHints_AcrossAllFrames(t *testing.T) {
	err := harborerr.New("m1", "h1", "h2").WithMessage("m2", "h3")
	assert.Equal(t, []string{"h1", "h2", "h3"}, err.Hints())
}

func TestFrames_ReturnsAll(t *testing.T) {
	err := harborerr.New("root", "hint-a")
	frames := err.Frames()
	require.Len(t, frames, 1)
	assert.Equal(t, "root", frames[0].Message)
	assert.Equal(t, []string{"hint-a"}, frames[0].Hints)
}

func TestCause_ReturnsUnderlyingError(t *testing.T) {
	sentinel := errors.New("sentinel")
	err := harborerr.NewWithCause(sentinel, "wrapper")
	assert.Equal(t, sentinel, err.Cause())
}

func TestCause_NilWhenNoCause(t *testing.T) {
	err := harborerr.New("standalone")
	assert.Nil(t, err.Cause())
}

func TestStatus_NoCause_ReturnsEmpty(t *testing.T) {
	assert.Equal(t, "", harborerr.Status(harborerr.New("no cause")))
}

func TestStatus_PlainError_ReturnsEmpty(t *testing.T) {
	assert.Equal(t, "", harborerr.Status(errors.New("plain")))
}

func TestUnwrap_ReturnsTheCause(t *testing.T) {
	sentinel := errors.New("sentinel")
	err := harborerr.NewWithCause(sentinel, "wrapper")

	assert.Equal(t, sentinel, err.Unwrap())
	assert.True(t, errors.Is(err, sentinel))
}

func TestErrorsAs_FindsHarborError(t *testing.T) {
	inner := harborerr.New("inner")
	outer := harborerr.NewWithCause(inner, "outer")

	var target *harborerr.Error
	assert.True(t, errors.As(outer, &target))
	assert.Contains(t, target.Error(), "inner")
	assert.Contains(t, target.Error(), "outer")
}

func TestIsError_True_DirectHarborError(t *testing.T) {
	assert.True(t, harborerr.IsError(harborerr.New("err")))
}

func TestIsError_True_WrappedWithFmtErrorf(t *testing.T) {
	wrapped := fmt.Errorf("outer: %w", harborerr.New("inner"))
	assert.True(t, harborerr.IsError(wrapped))
}

func TestIsError_False_PlainError(t *testing.T) {
	assert.False(t, harborerr.IsError(errors.New("plain")))
}

func TestAsError_FromHarborError_ReturnsSame(t *testing.T) {
	original := harborerr.New("original")
	wrapped := fmt.Errorf("wrapped: %w", original)

	result := harborerr.AsError(wrapped)
	require.NotNil(t, result)
	assert.Contains(t, result.Error(), "original")
}

func TestAsError_FromPlainError_WrapsIntoSingleFrame(t *testing.T) {
	plain := errors.New("plain error")
	result := harborerr.AsError(plain)
	require.NotNil(t, result)

	assert.Contains(t, result.Error(), "plain error")
	assert.Equal(t, []string{"plain error"}, result.Errors())
	assert.Equal(t, plain, result.Cause())
}

func TestCause_PackageLevel_HarborError(t *testing.T) {
	sentinel := errors.New("root")
	err := harborerr.NewWithCause(sentinel, "top")
	assert.Equal(t, sentinel, harborerr.Cause(err))
}

func TestCause_PackageLevel_PlainError(t *testing.T) {
	assert.Nil(t, harborerr.Cause(errors.New("plain")))
}

func TestHints_PackageLevel_HarborError(t *testing.T) {
	err := harborerr.New("error", "check config")
	assert.Equal(t, []string{"check config"}, harborerr.Hints(err))
}

func TestHints_PackageLevel_PlainError(t *testing.T) {
	assert.Nil(t, harborerr.Hints(errors.New("plain")))
}
