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
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	harborerr "github.com/goharbor/harbor-cli/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew_SingleFrame(t *testing.T) {
	err := harborerr.New("something went wrong")
	require.NotNil(t, err)
	assert.Equal(t, "something went wrong", err.Error())
	assert.Equal(t, "something went wrong", err.Message())
	assert.Equal(t, []string{"something went wrong"}, err.Errors())
	assert.Empty(t, err.Hints())
	assert.Nil(t, err.Cause())
}

func TestNewf_FormatsMessage(t *testing.T) {
	err := harborerr.Newf("resource %q not found", "project-x")
	require.NotNil(t, err)
	assert.Equal(t, `resource "project-x" not found`, err.Error())
	assert.Equal(t, `resource "project-x" not found`, err.Message())
}

func TestWrap_PushesOutermostFrame(t *testing.T) {
	root := harborerr.New("root problem")
	wrapped := harborerr.Wrap(root, "operation failed")

	assert.Equal(t, "root problem", wrapped.Error())
	assert.Equal(t, "operation failed", wrapped.Message())
	assert.Equal(t, []string{"operation failed", "root problem"}, wrapped.Errors())
}

func TestWrapf_FormatsOutermostFrame(t *testing.T) {
	root := harborerr.New("auth failed")
	wrapped := harborerr.Wrapf(root, "login for user %q", "alice")

	assert.Equal(t, "auth failed", wrapped.Error())
	assert.Equal(t, `login for user "alice"`, wrapped.Message())
	assert.Equal(t, []string{`login for user "alice"`, "auth failed"}, wrapped.Errors())
}

func TestWrap_ThreeLevels(t *testing.T) {
	root := harborerr.New("db timeout")
	mid := harborerr.Wrap(root, "repository unavailable")
	top := harborerr.Wrap(mid, "delete artifact failed")

	assert.Equal(t, "db timeout", top.Error())
	assert.Equal(t, "delete artifact failed", top.Message())
	assert.Equal(t, []string{"delete artifact failed", "repository unavailable", "db timeout"}, top.Errors())
}

func TestWrap_PlainStdlibError(t *testing.T) {
	plain := errors.New("network error")
	wrapped := harborerr.Wrap(plain, "could not reach registry")

	assert.Equal(t, "network error", wrapped.Error())
	assert.Equal(t, "could not reach registry", wrapped.Message())
	assert.Equal(t, []string{"could not reach registry", "network error"}, wrapped.Errors())
}

func TestErrors_SingleFrame(t *testing.T) {
	err := harborerr.New("standalone")
	assert.Equal(t, []string{"standalone"}, err.Errors())
}

func TestErrors_OutermostFirst(t *testing.T) {
	err := harborerr.Wrap(harborerr.Wrap(harborerr.New("level-0"), "level-1"), "level-2")
	assert.Equal(t, []string{"level-2", "level-1", "level-0"}, err.Errors())
}

func TestWithHint_AttachesToOutermostFrame(t *testing.T) {
	err := harborerr.New("error").WithHint("try again later")
	assert.Equal(t, []string{"try again later"}, err.Hints())
}

func TestWithHint_MultipleHintsOnSameFrame(t *testing.T) {
	err := harborerr.New("error").
		WithHint("hint one").
		WithHint("hint two").
		WithHint("hint three")
	assert.Equal(t, []string{"hint one", "hint two", "hint three"}, err.Hints())
}

func TestHints_AcrossFrames_OutermostFirst(t *testing.T) {
	root := harborerr.New("root").WithHint("root-hint")
	top := harborerr.Wrap(root, "top").WithHint("top-hint")

	assert.Equal(t, []string{"top-hint", "root-hint"}, top.Hints())
}

func TestHints_PackageLevel_HarborError(t *testing.T) {
	err := harborerr.New("error").WithHint("check config")
	assert.Equal(t, []string{"check config"}, harborerr.Hints(err))
}

func TestHints_PackageLevel_PlainError(t *testing.T) {
	assert.Nil(t, harborerr.Hints(errors.New("plain")))
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
	assert.Equal(t, "original", result.Error())
}

func TestAsError_FromPlainError_WrapsIntoSingleFrame(t *testing.T) {
	plain := errors.New("plain error")
	result := harborerr.AsError(plain)
	require.NotNil(t, result)
	assert.Equal(t, "plain error", result.Error())
	assert.Equal(t, []string{"plain error"}, result.Errors())
	assert.Equal(t, plain, result.Cause())
}

func TestWithCause_AttachesCauseForUnwrap(t *testing.T) {
	sentinel := errors.New("sentinel")
	err := harborerr.New("wrapper").WithCause(sentinel)

	assert.Equal(t, sentinel, err.Cause())
	assert.True(t, errors.Is(err, sentinel))
}

func TestWithCause_OnlyFirstCauseIsStored(t *testing.T) {
	first := errors.New("first")
	second := errors.New("second")
	err := harborerr.New("e").WithCause(first).WithCause(second)
	assert.Equal(t, first, err.Cause())
}

func TestCause_PackageLevel_HarborError(t *testing.T) {
	sentinel := errors.New("root")
	err := harborerr.New("top").WithCause(sentinel)
	assert.Equal(t, sentinel, harborerr.Cause(err))
}

func TestCause_PackageLevel_PlainError(t *testing.T) {
	assert.Nil(t, harborerr.Cause(errors.New("plain")))
}

func TestUnwrap_ErrorsAs_FindsOuterFrame(t *testing.T) {
	inner := harborerr.New("inner")
	outer := harborerr.New("outer").WithCause(inner)

	var target *harborerr.Error
	assert.True(t, errors.As(outer, &target))
	assert.Equal(t, "outer", target.Error())
}

func TestStatus_PlainError_ReturnsEmpty(t *testing.T) {
	assert.Equal(t, "", harborerr.Status(errors.New("plain")))
}

func TestStatus_NoCause_ReturnsEmpty(t *testing.T) {
	assert.Equal(t, "", harborerr.Status(harborerr.New("no cause")))
}

type harborPayloadError struct {
	Payload *harborPayloadBody
}

type harborPayloadBody struct {
	Errors []harborPayloadEntry `json:"errors"`
}

type harborPayloadEntry struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (h *harborPayloadError) Error() string {
	if h.Payload != nil && len(h.Payload.Errors) > 0 {
		b, _ := json.Marshal(h.Payload)
		return fmt.Sprintf("[%s] %s", h.Payload.Errors[0].Code, string(b))
	}
	return "harbor error"
}

func TestWrap_HarborPayloadCause_ExtractsMessage(t *testing.T) {
	apiErr := &harborPayloadError{
		Payload: &harborPayloadBody{
			Errors: []harborPayloadEntry{
				{Code: "NOT_FOUND", Message: "repository does not exist"},
			},
		},
	}
	err := harborerr.Wrap(apiErr, "delete artifact failed")

	assert.Equal(t, "repository does not exist", err.Error())
	assert.Equal(t, "delete artifact failed", err.Message())
	assert.Equal(t, []string{"delete artifact failed", "repository does not exist"}, err.Errors())
}

func TestFrames_SingleFrame(t *testing.T) {
	err := harborerr.New("root").WithHint("hint-a")
	frames := err.Frames()
	require.Len(t, frames, 1)
	assert.Equal(t, "root", frames[0].Message)
	assert.Equal(t, []string{"hint-a"}, frames[0].Hints)
}

func TestFrames_MultipleFrames_OutermostFirst(t *testing.T) {
	root := harborerr.New("root").WithHint("root-hint")
	mid := harborerr.Wrap(root, "mid").WithHint("mid-hint")
	top := harborerr.Wrap(mid, "top").WithHint("top-hint")

	frames := top.Frames()
	require.Len(t, frames, 3)
	assert.Equal(t, "top", frames[0].Message)
	assert.Equal(t, []string{"top-hint"}, frames[0].Hints)
	assert.Equal(t, "mid", frames[1].Message)
	assert.Equal(t, []string{"mid-hint"}, frames[1].Hints)
	assert.Equal(t, "root", frames[2].Message)
	assert.Equal(t, []string{"root-hint"}, frames[2].Hints)
}

func TestFrames_NoHints(t *testing.T) {
	err := harborerr.Wrap(harborerr.New("root"), "top")
	frames := err.Frames()
	require.Len(t, frames, 2)
	assert.Empty(t, frames[0].Hints)
	assert.Empty(t, frames[1].Hints)
}

func TestChaining_WrapWithHints(t *testing.T) {
	root := harborerr.New("connection refused").WithHint("check firewall rules")
	top := harborerr.Wrap(root, "could not reach registry").WithHint("verify server URL")

	assert.Equal(t, "connection refused", top.Error())
	assert.Equal(t, "could not reach registry", top.Message())
	assert.Equal(t, []string{"could not reach registry", "connection refused"}, top.Errors())
	assert.Equal(t, []string{"verify server URL", "check firewall rules"}, top.Hints())
}
