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
package log

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"sync"
	"time"
)

// Fields is a type alias for structured log fields, matching logrus.Fields.
type Fields map[string]any

// Level represents a log severity level.
type Level = slog.Level

const (
	LevelDebug = slog.LevelDebug
	LevelInfo  = slog.LevelInfo
	LevelWarn  = slog.LevelWarn
	LevelError = slog.LevelError

	// Logrus-compatible aliases.
	InfoLevel  = slog.LevelInfo
	DebugLevel = slog.LevelDebug
	WarnLevel  = slog.LevelWarn
	ErrorLevel = slog.LevelError
)

// TextFormatter mirrors logrus.TextFormatter for backward-compatible configuration.
type TextFormatter struct {
	FullTimestamp   bool
	TimestampFormat string
	DisableColors   bool
}

// Logger wraps slog.Logger to provide a logrus-compatible API surface.
type Logger struct {
	logger    *slog.Logger
	Out       io.Writer // exported field matching logrus.Logger.Out
	level     Level
	formatter *TextFormatter
}

var (
	mu            sync.Mutex
	defaultLogger *Logger
	currentLevel  Level     = LevelInfo
	currentOutput io.Writer = os.Stderr
)

func init() {
	rebuildDefaultLogger()
}

func rebuildDefaultLogger() {
	opts := &slog.HandlerOptions{Level: currentLevel}
	defaultLogger = &Logger{
		logger: slog.New(slog.NewTextHandler(currentOutput, opts)),
		Out:    currentOutput,
		level:  currentLevel,
	}
}

// SetLevel sets the minimum log level on the default logger.
func SetLevel(level Level) {
	mu.Lock()
	defer mu.Unlock()
	currentLevel = level
	rebuildDefaultLogger()
}

// SetOutput sets the output destination on the default logger.
func SetOutput(w io.Writer) {
	mu.Lock()
	defer mu.Unlock()
	currentOutput = w
	rebuildDefaultLogger()
}

// New creates a new Logger instance.
func New() *Logger {
	return &Logger{
		logger: slog.New(slog.NewTextHandler(os.Stderr, nil)),
		Out:    os.Stderr,
		level:  LevelInfo,
	}
}

// StandardLogger returns the default package-level logger.
func StandardLogger() *Logger {
	return defaultLogger
}

// Out returns the underlying io.Writer.

// SetLevel sets the minimum log level on this Logger instance.
func (l *Logger) SetLevel(level Level) {
	l.level = level
	l.rebuildHandler()
}

// SetFormatter configures the formatter for this Logger instance.
func (l *Logger) SetFormatter(f *TextFormatter) {
	l.formatter = f
	l.rebuildHandler()
}

// SetOutput sets the output destination on this Logger instance.
func (l *Logger) SetOutput(w io.Writer) {
	l.Out = w
	l.rebuildHandler()
}

func (l *Logger) rebuildHandler() {
	opts := &slog.HandlerOptions{Level: l.level}
	if l.formatter != nil && l.formatter.TimestampFormat != "" {
		opts.ReplaceAttr = func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey && len(groups) == 0 {
				a.Value = slog.StringValue(time.Now().Format(l.formatter.TimestampFormat))
			}
			return a
		}
	}
	l.logger = slog.New(slog.NewTextHandler(l.Out, opts))
}

// WithTime returns a Logger with the given time added as the "time" attribute.
func (l *Logger) WithTime(t time.Time) *Logger {
	return l.WithField("time", t)
}

// WithField returns a Logger with the given key-value pair added to its context.
func WithField(key string, value any) *Logger {
	return defaultLogger.WithField(key, value)
}

// WithField returns a Logger with the given key-value pair added to its context.
func (l *Logger) WithField(key string, value any) *Logger {
	return &Logger{
		logger:    l.logger.With(key, value),
		Out:       l.Out,
		level:     l.level,
		formatter: l.formatter,
	}
}

// WithFields returns a Logger with the given key-value pairs added to its context.
func WithFields(fields Fields) *Logger {
	return defaultLogger.WithFields(fields)
}

// WithFields returns a Logger with the given key-value pairs added to its context.
func (l *Logger) WithFields(fields Fields) *Logger {
	attrs := make([]any, 0, len(fields)*2)
	for k, v := range fields {
		attrs = append(attrs, k, v)
	}
	return &Logger{
		logger:    l.logger.With(attrs...),
		Out:       l.Out,
		level:     l.level,
		formatter: l.formatter,
	}
}

// --- Package-level convenience methods ---

// Debug logs a message at Debug level.
func Debug(args ...any) {
	defaultLogger.logger.Debug(fmt.Sprint(args...))
}

func (l *Logger) Debug(args ...any) {
	l.logger.Debug(fmt.Sprint(args...))
}

// Debugf logs a formatted message at Debug level.
func Debugf(format string, args ...any) {
	defaultLogger.logger.Debug(fmt.Sprintf(format, args...))
}

func (l *Logger) Debugf(format string, args ...any) {
	l.logger.Debug(fmt.Sprintf(format, args...))
}

// Info logs a message at Info level.
func Info(args ...any) {
	defaultLogger.logger.Info(fmt.Sprint(args...))
}

func (l *Logger) Info(args ...any) {
	l.logger.Info(fmt.Sprint(args...))
}

// Infof logs a formatted message at Info level.
func Infof(format string, args ...any) {
	defaultLogger.logger.Info(fmt.Sprintf(format, args...))
}

func (l *Logger) Infof(format string, args ...any) {
	l.logger.Info(fmt.Sprintf(format, args...))
}

// Warn logs a message at Warn level.
func Warn(args ...any) {
	defaultLogger.logger.Warn(fmt.Sprint(args...))
}

func (l *Logger) Warn(args ...any) {
	l.logger.Warn(fmt.Sprint(args...))
}

// Warnf logs a formatted message at Warn level.
func Warnf(format string, args ...any) {
	defaultLogger.logger.Warn(fmt.Sprintf(format, args...))
}

func (l *Logger) Warnf(format string, args ...any) {
	l.logger.Warn(fmt.Sprintf(format, args...))
}

// Warningf is an alias for Warnf.
func Warningf(format string, args ...any) {
	Warnf(format, args...)
}

func (l *Logger) Warningf(format string, args ...any) {
	l.Warnf(format, args...)
}

// Error logs a message at Error level.
func Error(args ...any) {
	defaultLogger.logger.Error(fmt.Sprint(args...))
}

func (l *Logger) Error(args ...any) {
	l.logger.Error(fmt.Sprint(args...))
}

// Errorf logs a formatted message at Error level.
func Errorf(format string, args ...any) {
	defaultLogger.logger.Error(fmt.Sprintf(format, args...))
}

func (l *Logger) Errorf(format string, args ...any) {
	l.logger.Error(fmt.Sprintf(format, args...))
}

// Fatal logs a message at Error level and exits with status 1.
func Fatal(args ...any) {
	defaultLogger.logger.Error(fmt.Sprint(args...))
	os.Exit(1)
}

func (l *Logger) Fatal(args ...any) {
	l.logger.Error(fmt.Sprint(args...))
	os.Exit(1)
}

// Fatalf logs a formatted message at Error level and exits with status 1.
func Fatalf(format string, args ...any) {
	defaultLogger.logger.Error(fmt.Sprintf(format, args...))
	os.Exit(1)
}

func (l *Logger) Fatalf(format string, args ...any) {
	l.logger.Error(fmt.Sprintf(format, args...))
	os.Exit(1)
}

// Printf logs a formatted message at Info level.
func Printf(format string, args ...any) {
	defaultLogger.logger.Info(fmt.Sprintf(format, args...))
}

func (l *Logger) Printf(format string, args ...any) {
	l.logger.Info(fmt.Sprintf(format, args...))
}

// Println logs a message at Info level.
func Println(args ...any) {
	defaultLogger.logger.Info(fmt.Sprintln(args...))
}

func (l *Logger) Println(args ...any) {
	l.logger.Info(fmt.Sprintln(args...))
}
