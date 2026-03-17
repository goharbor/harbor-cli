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
	"context"
	"fmt"
	"io"
	"log/slog"
	"sync"

	"github.com/charmbracelet/lipgloss"
)

// Needs to be functions, otherwise doesn't work
var (
	debugStyle = func() lipgloss.Style { return lipgloss.NewStyle().Foreground(lipgloss.Color("8")) }
	infoStyle  = func() lipgloss.Style { return lipgloss.NewStyle().Foreground(lipgloss.Color("10")) }
	warnStyle  = func() lipgloss.Style { return lipgloss.NewStyle().Foreground(lipgloss.Color("11")) }
	errorStyle = func() lipgloss.Style { return lipgloss.NewStyle().Foreground(lipgloss.Color("9")) }
)

type PrettyHandler struct {
	mu    sync.Mutex
	out   io.Writer
	level slog.Leveler
}

func NewPrettyHandler(out io.Writer, level slog.Leveler) *PrettyHandler {
	return &PrettyHandler{
		out:   out,
		level: level,
	}
}

func (h *PrettyHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.level.Level()
}

func (h *PrettyHandler) Handle(_ context.Context, r slog.Record) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	timestamp := r.Time.Format("15:04:05")
	level := formatLevel(r.Level)

	// Print main log line
	_, err := fmt.Fprintf(
		h.out,
		"%s %s %s\n",
		fmt.Sprintf("%s |", timestamp),
		level,
		r.Message,
	)
	if err != nil {
		return err
	}

	var attrs []slog.Attr
	maxKey := 0

	r.Attrs(func(a slog.Attr) bool {
		attrs = append(attrs, a)

		if len(a.Key) > maxKey {
			maxKey = len(a.Key)
		}

		return true
	})

	if len(attrs) == 0 {
		return nil
	}

	// Print attributes
	for _, a := range attrs {
		_, err = fmt.Fprintf(h.out, "  %s : %v\n", fmt.Sprintf("%-*s", maxKey, a.Key),
			a.Value.String())
		if err != nil {
			return err
		}
	}

	return nil
}

func (h *PrettyHandler) WithAttrs(_ []slog.Attr) slog.Handler {
	return h
}

func (h *PrettyHandler) WithGroup(_ string) slog.Handler {
	return h
}

func formatLevel(level slog.Level) string {
	switch level {
	case slog.LevelDebug:
		return debugStyle().Render("[ DEBUG ]")
	case slog.LevelInfo:
		return infoStyle().Render("[ INFO  ]") // space is intentional for symmetry
	case slog.LevelWarn:
		return warnStyle().Render("[ WARN  ]") // space is intentional for symmetry
	case slog.LevelError:
		return errorStyle().Render("[ ERROR ]")
	default:
		return level.String()
	}
}
