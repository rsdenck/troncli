package logger

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
)

// Options for logger configuration
type Options struct {
	Debug    bool
	LogFile  string
	UseColor bool
}

// Init initializes the global logger
func Init(opts Options) error {
	var handler slog.Handler

	// Ensure log directory exists if log file is specified
	if opts.LogFile != "" {
		dir := filepath.Dir(opts.LogFile)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}

	// Create a multi-writer handler
	// 1. Console Handler (UI)
	//    - Info/Warn -> Stdout
	//    - Error -> Stderr
	// 2. File Handler (Debug)
	//    - All levels including Debug -> File

	consoleHandler := NewConsoleHandler(os.Stdout, os.Stderr, opts)
	
	if opts.LogFile != "" {
		f, err := os.OpenFile(opts.LogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		fileHandler := slog.NewJSONHandler(f, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})
		handler = NewMultiHandler(consoleHandler, fileHandler)
	} else {
		handler = consoleHandler
	}

	logger := slog.New(handler)
	slog.SetDefault(logger)

	return nil
}

// ConsoleHandler handles UI output
type ConsoleHandler struct {
	out   io.Writer
	err   io.Writer
	opts  Options
	attrs []slog.Attr
	group string
}

func NewConsoleHandler(out, err io.Writer, opts Options) *ConsoleHandler {
	return &ConsoleHandler{
		out:  out,
		err:  err,
		opts: opts,
	}
}

func (h *ConsoleHandler) Enabled(ctx context.Context, level slog.Level) bool {
	if h.opts.Debug {
		return true
	}
	return level >= slog.LevelInfo
}

func (h *ConsoleHandler) Handle(ctx context.Context, r slog.Record) error {
	// Simple text format for UI
	// Error -> Stderr, others -> Stdout
	w := h.out
	if r.Level >= slog.LevelError {
		w = h.err
	}

	// Format: [LEVEL] Message key=value
	// Or just Message for Info
	msg := r.Message
	
	// Write to writer
	_, err := io.WriteString(w, fmt.Sprintf("[%s] %s\n", r.Level, msg))
	return err
}

func (h *ConsoleHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &ConsoleHandler{
		out:   h.out,
		err:   h.err,
		opts:  h.opts,
		attrs: append(h.attrs, attrs...),
		group: h.group,
	}
}

func (h *ConsoleHandler) WithGroup(name string) slog.Handler {
	return &ConsoleHandler{
		out:   h.out,
		err:   h.err,
		opts:  h.opts,
		attrs: h.attrs,
		group: name,
	}
}

// MultiHandler dispatches to multiple handlers
type MultiHandler struct {
	handlers []slog.Handler
}

func NewMultiHandler(handlers ...slog.Handler) *MultiHandler {
	return &MultiHandler{handlers: handlers}
}

func (h *MultiHandler) Enabled(ctx context.Context, level slog.Level) bool {
	for _, handler := range h.handlers {
		if handler.Enabled(ctx, level) {
			return true
		}
	}
	return false
}

func (h *MultiHandler) Handle(ctx context.Context, r slog.Record) error {
	for _, handler := range h.handlers {
		if handler.Enabled(ctx, r.Level) {
			_ = handler.Handle(ctx, r)
		}
	}
	return nil
}

func (h *MultiHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	handlers := make([]slog.Handler, len(h.handlers))
	for i, handler := range h.handlers {
		handlers[i] = handler.WithAttrs(attrs)
	}
	return NewMultiHandler(handlers...)
}

func (h *MultiHandler) WithGroup(name string) slog.Handler {
	handlers := make([]slog.Handler, len(h.handlers))
	for i, handler := range h.handlers {
		handlers[i] = handler.WithGroup(name)
	}
	return NewMultiHandler(handlers...)
}
