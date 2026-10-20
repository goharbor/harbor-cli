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
	"strings"
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
	mu       sync.Mutex
	out      io.Writer
	level    slog.Leveler
	preAttrs []slog.Attr // retained from WithAttrs calls
	groups   []string    // retained from WithGroup calls
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

	// Merge pre-attached attrs with record attrs.
	var attrs []slog.Attr
	attrs = append(attrs, h.preAttrs...)
	r.Attrs(func(a slog.Attr) bool {
		attrs = append(attrs, a)
		return true
	})

	if len(attrs) == 0 {
		return nil
	}

	// Compute column width across all attrs.
	maxKey := 0
	for _, a := range attrs {
		if k := len(h.qualifiedKey(a.Key)); k > maxKey {
			maxKey = k
		}
	}

	for _, a := range attrs {
		_, err = fmt.Fprintf(h.out, "  %s : %v\n",
			fmt.Sprintf("%-*s", maxKey, h.qualifiedKey(a.Key)),
			a.Value.String(),
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func (h *PrettyHandler) qualifiedKey(key string) string {
	if len(h.groups) == 0 {
		return key
	}
	return strings.Join(h.groups, ".") + "." + key
}

func (h *PrettyHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newPreAttrs := make([]slog.Attr, len(h.preAttrs)+len(attrs))
	copy(newPreAttrs, h.preAttrs)
	copy(newPreAttrs[len(h.preAttrs):], attrs)

	return &PrettyHandler{
		preAttrs: newPreAttrs,
		groups:   h.groups,
	}
}

func (h *PrettyHandler) WithGroup(name string) slog.Handler {
	if name == "" {
		return h
	}
	newGroups := make([]string, len(h.groups)+1)
	copy(newGroups, h.groups)
	newGroups[len(h.groups)] = name

	return &PrettyHandler{
		preAttrs: h.preAttrs,
		groups:   newGroups,
	}
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
