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

//go:build e2e

// Package e2e holds tests that talk to a real Harbor instance.
//
// They are behind the `e2e` build tag so `go test ./...` stays offline. Start a
// throwaway Harbor with `test/harbor/harbor.sh up`, then run:
//
//	go test -tags e2e ./test/e2e/...
package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"
	"time"
)

// Defaults match the stack in test/harbor/docker-compose.yaml.
const (
	defaultServer   = "http://localhost:8080"
	defaultUsername = "admin"
	defaultPassword = "Harbor12345"
)

// startupTimeout bounds how long TestMain waits for Harbor to answer. Harbor
// needs roughly a minute from a cold start before core serves the API.
const startupTimeout = 3 * time.Minute

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// Server is the Harbor instance under test, e.g. http://localhost:8080.
func Server() string { return envOr("HARBOR_URL", defaultServer) }

// Username is an account with system admin rights on Server.
func Username() string { return envOr("HARBOR_USERNAME", defaultUsername) }

// Password belongs to Username.
func Password() string { return envOr("HARBOR_PASSWORD", defaultPassword) }

func TestMain(m *testing.M) {
	if err := waitForHarbor(Server(), startupTimeout); err != nil {
		fmt.Fprintf(os.Stderr, `
Harbor is not reachable at %s: %v

Start a throwaway instance and try again:

    ./test/harbor/harbor.sh up
    go test -tags e2e ./test/e2e/...

Point the tests somewhere else with HARBOR_URL, HARBOR_USERNAME and HARBOR_PASSWORD.
`, Server(), err)
		os.Exit(1)
	}
	os.Exit(m.Run())
}

// waitForHarbor polls the unauthenticated ping endpoint until Harbor answers.
func waitForHarbor(server string, timeout time.Duration) error {
	client := &http.Client{Timeout: 5 * time.Second}
	deadline := time.Now().Add(timeout)

	var lastErr error
	for {
		resp, err := client.Get(server + "/api/v2.0/ping")
		if err == nil {
			resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				return nil
			}
			lastErr = fmt.Errorf("unexpected status %s", resp.Status)
		} else {
			lastErr = err
		}

		if time.Now().After(deadline) {
			return fmt.Errorf("gave up after %s: %w", timeout, lastErr)
		}
		time.Sleep(2 * time.Second)
	}
}

// robot is the subset of Harbor's robot account response the tests need.
type robot struct {
	Name   string `json:"name"`
	Secret string `json:"secret"`
	ID     int64  `json:"id"`
}

// createRobot creates a system level robot account and registers its removal
// with t.Cleanup.
func createRobot(t *testing.T, name string) robot {
	t.Helper()

	body := map[string]any{
		"name":        name,
		"description": "created by harbor-cli e2e tests",
		"duration":    1,
		"level":       "system",
		"permissions": []map[string]any{
			{
				"kind":      "system",
				"namespace": "/",
				"access":    []map[string]string{{"resource": "project", "action": "list"}},
			},
		},
	}

	var created robot
	doJSON(t, http.MethodPost, "/api/v2.0/robots", body, &created)

	t.Cleanup(func() {
		doJSON(t, http.MethodDelete, fmt.Sprintf("/api/v2.0/robots/%d", created.ID), nil, nil)
	})

	return created
}

// doJSON performs an authenticated API call against Server and decodes the
// response into out when out is non-nil.
func doJSON(t *testing.T, method, path string, body, out any) {
	t.Helper()

	var payload io.Reader
	if body != nil {
		encoded, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("marshal %s %s: %v", method, path, err)
		}
		payload = bytes.NewReader(encoded)
	}

	req, err := http.NewRequest(method, Server()+path, payload)
	if err != nil {
		t.Fatalf("build request %s %s: %v", method, path, err)
	}
	req.SetBasicAuth(Username(), Password())
	req.Header.Set("Content-Type", "application/json")

	resp, err := (&http.Client{Timeout: 30 * time.Second}).Do(req)
	if err != nil {
		t.Fatalf("%s %s: %v", method, path, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		t.Fatalf("%s %s: unexpected status %s", method, path, resp.Status)
	}
	if out == nil {
		return
	}
	if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
		t.Fatalf("decode %s %s: %v", method, path, err)
	}
}
