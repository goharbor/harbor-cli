package logger

import (
	"bytes"
	"context"
	"log/slog"
	"strings"
	"testing"
)

func TestEnabled(t *testing.T) {
	ctx := context.Background()

	handlers := []struct {
		name  string
		level slog.Level
	}{
		{"DebugHandler", slog.LevelDebug},
		{"InfoHandler", slog.LevelInfo},
		{"WarnHandler", slog.LevelWarn},
		{"ErrorHandler", slog.LevelError},
	}

	levels := []struct {
		name  string
		level slog.Level
	}{
		{"Debug", slog.LevelDebug},
		{"Info", slog.LevelInfo},
		{"Warn", slog.LevelWarn},
		{"Error", slog.LevelError},
	}

	for _, handler := range handlers {
		h := NewPrettyHandler(&bytes.Buffer{}, handler.level)

		t.Run(handler.name, func(t *testing.T) {
			for _, lvl := range levels {

				expected := lvl.level >= handler.level

				got := h.Enabled(ctx, lvl.level)

				if got != expected {
					t.Fatalf(
						"handler level=%s Enabled(%s) = %v, expected %v",
						handler.level.String(),
						lvl.level.String(),
						got,
						expected,
					)
				}
			}
		})
	}
}

func TestLogFormatting(t *testing.T) {
	tests := []struct {
		name    string
		logFunc func(*slog.Logger)
		level   string
		message string
	}{
		{
			name: "debug",
			logFunc: func(l *slog.Logger) {
				l.Debug("debug message")
			},
			level:   "DEBUG",
			message: "debug message",
		},
		{
			name: "info",
			logFunc: func(l *slog.Logger) {
				l.Info("info message")
			},
			level:   "INFO",
			message: "info message",
		},
		{
			name: "warn",
			logFunc: func(l *slog.Logger) {
				l.Warn("warn message")
			},
			level:   "WARN",
			message: "warn message",
		},
		{
			name: "error",
			logFunc: func(l *slog.Logger) {
				l.Error("error message")
			},
			level:   "ERROR",
			message: "error message",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer

			h := NewPrettyHandler(&buf, slog.LevelDebug)
			logger := slog.New(h)

			tt.logFunc(logger)

			out := buf.String()

			if !strings.Contains(out, tt.message) {
				t.Fatalf("expected message %q in output: %s", tt.message, out)
			}

			if !strings.Contains(out, tt.level) {
				t.Fatalf("expected level %q in output: %s", tt.level, out)
			}

			if !strings.Contains(out, "[") || !strings.Contains(out, "]") {
				t.Fatalf("expected level brackets in output: %s", out)
			}

			if !strings.Contains(out, " | ") {
				t.Fatalf("expected timestamp separator ' | ': %s", out)
			}

			ts := strings.Split(out, " | ")[0]

			if len(ts) != 8 || ts[2] != ':' || ts[5] != ':' {
				t.Fatalf("expected timestamp format HH:MM:SS, got: %s", ts)
			}
		})
	}
}

func TestHandleLevelFiltering(t *testing.T) {
	tests := []struct {
		name         string
		handlerLevel slog.Level
		logLevel     slog.Level
		shouldAppear bool
	}{
		{"warn_handler_info_log", slog.LevelWarn, slog.LevelInfo, false},
		{"warn_handler_warn_log", slog.LevelWarn, slog.LevelWarn, true},
		{"warn_handler_error_log", slog.LevelWarn, slog.LevelError, true},

		{"info_handler_debug_log", slog.LevelInfo, slog.LevelDebug, false},
		{"info_handler_info_log", slog.LevelInfo, slog.LevelInfo, true},
		{"info_handler_warn_log", slog.LevelInfo, slog.LevelWarn, true},

		{"error_handler_warn_log", slog.LevelError, slog.LevelWarn, false},
		{"error_handler_error_log", slog.LevelError, slog.LevelError, true},

		{"debug_handler_debug_log", slog.LevelDebug, slog.LevelDebug, true},
		{"debug_handler_info_log", slog.LevelDebug, slog.LevelInfo, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer

			h := NewPrettyHandler(&buf, tt.handlerLevel)
			logger := slog.New(h)

			msg := "test message"

			switch tt.logLevel {
			case slog.LevelDebug:
				logger.Debug(msg)
			case slog.LevelInfo:
				logger.Info(msg)
			case slog.LevelWarn:
				logger.Warn(msg)
			case slog.LevelError:
				logger.Error(msg)
			}

			out := buf.String()

			if tt.shouldAppear && !strings.Contains(out, msg) {
				t.Fatalf("expected log to appear for handler level %s and log level %s",
					tt.handlerLevel.String(), tt.logLevel.String())
			}

			if !tt.shouldAppear && strings.Contains(out, msg) {
				t.Fatalf("expected log NOT to appear for handler level %s and log level %s",
					tt.handlerLevel.String(), tt.logLevel.String())
			}
		})
	}
}

func TestHandlerAttrAndGroupBehavior(t *testing.T) {
	tests := []struct {
		name string
		run  func(h *PrettyHandler) slog.Handler
	}{
		{
			name: "WithAttrs",
			run: func(h *PrettyHandler) slog.Handler {
				return h.WithAttrs([]slog.Attr{
					slog.String("key", "value"),
				})
			},
		},
		{
			name: "WithGroup",
			run: func(h *PrettyHandler) slog.Handler {
				return h.WithGroup("test")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewPrettyHandler(&bytes.Buffer{}, slog.LevelInfo)

			result := tt.run(h)

			if result != h {
				t.Fatalf("%s should return the same handler instance", tt.name)
			}
		})
	}
}

func TestFormatLevel(t *testing.T) {
	tests := []struct {
		name  string
		level slog.Level
		want  string
	}{
		{"debug", slog.LevelDebug, "DEBUG"},
		{"info", slog.LevelInfo, "INFO"},
		{"warn", slog.LevelWarn, "WARN"},
		{"error", slog.LevelError, "ERROR"},
		{"unknown", slog.Level(999), slog.Level(999).String()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := formatLevel(tt.level)

			if !strings.Contains(out, tt.want) {
				t.Fatalf("expected formatted level to contain %q, got: %s", tt.want, out)
			}
		})
	}
}
