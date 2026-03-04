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
package logger

import (
	"bytes"
	"log/slog"
	"os"
	"strings"
	"testing"
)

func TestSetup(t *testing.T) {
	tests := []struct {
		name         string
		verbose      bool
		format       string
		logFunc      func()
		expectLog    bool
		expectJSON   bool
		expectPretty bool
	}{
		{
			name:      "pretty_warn_filters_info",
			verbose:   false,
			format:    "",
			logFunc:   func() { slog.Info("info message") },
			expectLog: false,
		},
		{
			name:         "pretty_warn_allows_warn",
			verbose:      false,
			format:       "",
			logFunc:      func() { slog.Warn("warn message") },
			expectLog:    true,
			expectPretty: true,
		},
		{
			name:         "verbose_pretty_allows_debug",
			verbose:      true,
			format:       "",
			logFunc:      func() { slog.Debug("debug message") },
			expectLog:    true,
			expectPretty: true,
		},
		{
			name:       "json_format_output",
			verbose:    true,
			format:     "json",
			logFunc:    func() { slog.Info("json message") },
			expectLog:  true,
			expectJSON: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer

			// capture stderr
			old := os.Stderr
			r, w, _ := os.Pipe()
			os.Stderr = w

			Setup(tt.verbose, tt.format)

			tt.logFunc()

			w.Close()
			os.Stderr = old

			_, err := buf.ReadFrom(r)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			out := buf.String()

			if tt.expectLog && out == "" {
				t.Fatalf("expected log output but got none")
			}

			if !tt.expectLog && out != "" {
				t.Fatalf("expected no log output, got: %s", out)
			}

			if tt.expectJSON && !strings.Contains(out, "{") {
				t.Fatalf("expected JSON formatted output: %s", out)
			}

			if tt.expectPretty && !strings.Contains(out, "|") {
				t.Fatalf("expected pretty log format with '|': %s", out)
			}
		})
	}
}
