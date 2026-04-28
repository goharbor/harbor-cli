// Copyright Project Harbor Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0

package label

import (
	"testing"
)

// Regression for https://github.com/goharbor/harbor-cli/issues/814.
//
// The shared label-lookup helper must default to `Scope: "p"` so a
// caller with no project name still produces a valid request shape
// (Harbor's ListLabels API rejects an empty `scope` with
// `invalid scope:`). The version that resolves a project name into
// an ID needs the live API, so we only exercise the empty-name
// branch here — that is the path that previously crashed every
// invocation when the artifact reference parser produced an empty
// projectName.

func TestBuildProjectListFlags_EmptyProjectName(t *testing.T) {
	got, err := buildProjectListFlags("")
	if err != nil {
		t.Fatalf("buildProjectListFlags(\"\") returned error: %v", err)
	}
	if got.Scope != "p" {
		t.Errorf("Scope: got %q, want %q", got.Scope, "p")
	}
	if got.ProjectID != 0 {
		t.Errorf("ProjectID: got %d, want 0 (no project to resolve)", got.ProjectID)
	}
}
