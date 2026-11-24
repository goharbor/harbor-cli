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
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"testing"
)

// TestViewPackageNames guards against package clauses copied from other
// view packages. It walks every select/view.go and list/view.go under
// pkg/views: select views use their parent directory name, list views
// use "list".
func TestViewPackageNames(t *testing.T) {
	var checked int
	err := fs.WalkDir(os.DirFS("."), ".", func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || path.Base(p) != "view.go" {
			return nil
		}
		dir := path.Dir(p)
		var want string
		switch path.Base(dir) {
		case "select":
			want = path.Base(path.Dir(dir))
		case "list":
			want = "list"
		default:
			return nil
		}
		checked++
		f, err := parser.ParseFile(token.NewFileSet(), filepath.FromSlash(p), nil, parser.PackageClauseOnly)
		if err != nil {
			t.Errorf("parse %s: %v", p, err)
			return nil
		}
		if got := f.Name.Name; got != want {
			t.Errorf("package in %s = %q, want %q", p, got, want)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("walk pkg/views: %v", err)
	}
	if checked == 0 {
		t.Fatal("no select/list view.go files found")
	}
}
