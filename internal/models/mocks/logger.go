package mocks

import (
	"context"
	"log/slog"
)

type SpyLogger struct {
	Called  bool
	Entries []slog.Record
}

func (s *SpyLogger) Enabled(ctx context.Context, level slog.Level) bool {
	return true
}

func (s *SpyLogger) Handle(ctx context.Context, r slog.Record) error {
	s.Called = true
	s.Entries = append(s.Entries, r)
	return nil
}

func (s *SpyLogger) WithAttrs(attrs []slog.Attr) slog.Handler {
	return s
}

func (s *SpyLogger) WithGroup(name string) slog.Handler {
	return s
}
