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
package views_test

import (
	"go/parser"
	"go/token"
	"path/filepath"
	"testing"
)

// TestViewPackageNames guards against package clauses copied from other
// view packages. Select views use the feature name; list views use "list".
func TestViewPackageNames(t *testing.T) {
	tests := []struct {
		file string
		want string
	}{
		{file: "artifact/select/view.go", want: "artifact"},
		{file: "artifact/tags/select/view.go", want: "tags"},
		{file: "repository/select/view.go", want: "repository"},
		{file: "webhook/select/view.go", want: "webhook"},
		{file: "label/select/view.go", want: "label"},
		{file: "replication/task/list/view.go", want: "list"},
	}

	for _, tt := range tests {
		t.Run(tt.file, func(t *testing.T) {
			f, err := parser.ParseFile(token.NewFileSet(), filepath.FromSlash(tt.file), nil, parser.PackageClauseOnly)
			if err != nil {
				t.Fatalf("parse %s: %v", tt.file, err)
			}
			if got := f.Name.Name; got != tt.want {
				t.Errorf("package in %s = %q, want %q", tt.file, got, tt.want)
			}
		})
	}
}
