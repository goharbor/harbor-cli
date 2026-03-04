package logger

import (
	"context"
	"fmt"
	"io"
	"log/slog"

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
	timestamp := r.Time.Format("15:04:05")
	level := formatLevel(r.Level)
	msg := r.Message

	_, err := fmt.Fprintf(
		h.out,
		"%s | [ %s ] %s\n",
		timestamp,
		level,
		msg,
	)

	return err
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
		return debugStyle().Render("DEBUG")
	case slog.LevelInfo:
		return infoStyle().Render("INFO ") // space is intentional for symmetry
	case slog.LevelWarn:
		return warnStyle().Render("WARN ") // space is intentional for symmetry
	case slog.LevelError:
		return errorStyle().Render("ERROR")
	default:
		return level.String()
	}
}
