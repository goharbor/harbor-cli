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
package utils_test

import (
	"testing"

	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestBuildQueryParam_Empty(t *testing.T) {
	result, err := utils.BuildQueryParam(nil, nil, nil, []string{"name"})
	assert.NoError(t, err)
	assert.Equal(t, "", result)
}

func TestBuildQueryParam_Fuzzy(t *testing.T) {
	result, err := utils.BuildQueryParam(
		[]string{"name=alice"},
		nil, nil,
		[]string{"name"},
	)
	assert.NoError(t, err)
	assert.Equal(t, "name=~alice", result)
}

func TestBuildQueryParam_FuzzyMultiple(t *testing.T) {
	result, err := utils.BuildQueryParam(
		[]string{"name=alice", "description=test"},
		nil, nil,
		[]string{"name", "description"},
	)
	assert.NoError(t, err)
	assert.Equal(t, "name=~alice,description=~test", result)
}

func TestBuildQueryParam_Match(t *testing.T) {
	result, err := utils.BuildQueryParam(
		nil,
		[]string{"name=alice"},
		nil,
		[]string{"name"},
	)
	assert.NoError(t, err)
	assert.Equal(t, "name=alice", result)
}

func TestBuildQueryParam_MatchMultiple(t *testing.T) {
	result, err := utils.BuildQueryParam(
		nil,
		[]string{"name=alice", "project_id=42"},
		nil,
		[]string{"name", "project_id"},
	)
	assert.NoError(t, err)
	assert.Equal(t, "name=alice,project_id=42", result)
}

func TestBuildQueryParam_Range(t *testing.T) {
	result, err := utils.BuildQueryParam(
		nil, nil,
		[]string{"size=1~100"},
		[]string{"size"},
	)
	assert.NoError(t, err)
	assert.Equal(t, "size=[1~100]", result)
}

func TestBuildQueryParam_RangeMultiple(t *testing.T) {
	result, err := utils.BuildQueryParam(
		nil, nil,
		[]string{"size=1~100", "pull_count=10~50"},
		[]string{"size", "pull_count"},
	)
	assert.NoError(t, err)
	assert.Equal(t, "size=[1~100],pull_count=[10~50]", result)
}

func TestBuildQueryParam_Combined(t *testing.T) {
	result, err := utils.BuildQueryParam(
		[]string{"name=alice"},
		[]string{"project_id=42"},
		[]string{"size=1~100"},
		[]string{"name", "project_id", "size"},
	)
	assert.NoError(t, err)
	assert.Equal(t, "name=~alice,project_id=42,size=[1~100]", result)
}

func TestBuildQueryParam_InvalidFuzzyFormat(t *testing.T) {
	_, err := utils.BuildQueryParam(
		[]string{"invalid"},
		nil, nil,
		[]string{"name"},
	)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid fuzzy arg")
}

func TestBuildQueryParam_InvalidMatchFormat(t *testing.T) {
	_, err := utils.BuildQueryParam(
		nil,
		[]string{"invalid"},
		nil,
		[]string{"name"},
	)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid match arg")
}

func TestBuildQueryParam_InvalidRangeFormat(t *testing.T) {
	_, err := utils.BuildQueryParam(
		nil, nil,
		[]string{"invalid"},
		[]string{"size"},
	)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid range arg")
}

func TestBuildQueryParam_InvalidRangeValue(t *testing.T) {
	_, err := utils.BuildQueryParam(
		nil, nil,
		[]string{"size=1"},
		[]string{"size"},
	)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid range arg")
}

func TestBuildQueryParam_InvalidKey(t *testing.T) {
	_, err := utils.BuildQueryParam(
		[]string{"badkey=value"},
		nil, nil,
		[]string{"name"},
	)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid key")
}

func TestBuildQueryParam_InvalidKeyMatch(t *testing.T) {
	_, err := utils.BuildQueryParam(
		nil,
		[]string{"badkey=value"},
		nil,
		[]string{"name"},
	)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid key")
}

func TestBuildQueryParam_InvalidKeyRange(t *testing.T) {
	_, err := utils.BuildQueryParam(
		nil, nil,
		[]string{"badkey=1~100"},
		[]string{"name"},
	)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid key")
}
