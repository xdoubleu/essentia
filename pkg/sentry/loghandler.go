package sentry

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/getsentry/sentry-go"
	"github.com/xdoubleu/essentia/pkg/config"
)

// LogHandler is used for capturing logs and sending these to Sentry.
type LogHandler struct {
	level   slog.Level
	handler slog.Handler
	goas    []groupOrAttrs
}

type groupOrAttrs struct {
	group string      // group name if non-empty
	attrs []slog.Attr // attrs if non-empty
}

// NewLogHandler returns a new [SentryLogHandler].
func NewLogHandler(env string, handler slog.Handler) slog.Handler {
	level := slog.LevelInfo

	if env == config.DevEnv {
		level = slog.LevelDebug
	}

	return &LogHandler{
		handler: handler,
		level:   level,
		goas:    []groupOrAttrs{},
	}
}

// Enabled checks if logs are enabled in
// a [LogHandler] for a certain [slog.Level].
func (l *LogHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= l.level
}

// WithAttrs adds [[]slog.Attr] to a [SentryLogHandler].
func (l *LogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	if len(attrs) == 0 {
		return l
	}

	l.handler = l.handler.WithAttrs(attrs)
	return l.withGroupOrAttrs(groupOrAttrs{group: "", attrs: attrs})
}

// WithGroup adds a group to a [SentryLogHandler].
func (l *LogHandler) WithGroup(name string) slog.Handler {
	if name == "" {
		return l
	}

	l.handler = l.handler.WithGroup(name)
	return l.withGroupOrAttrs(groupOrAttrs{group: name, attrs: []slog.Attr{}})
}

func (l *LogHandler) withGroupOrAttrs(goa groupOrAttrs) slog.Handler {
	l2 := *l
	l2.goas = make([]groupOrAttrs, len(l.goas)+1)
	copy(l2.goas, l.goas)
	l2.goas[len(l2.goas)-1] = goa
	return &l2
}

// Handle handles a [slog.Record] by a [SentryLogHandler].
func (l *LogHandler) Handle(ctx context.Context, record slog.Record) error {
	if record.Level == slog.LevelError {
		l.sendErrorToSentry(ctx, recordToError(record))
	}

	return l.handler.Handle(ctx, record)
}

func (l *LogHandler) sendErrorToSentry(ctx context.Context, err error) {
	if hub := sentry.GetHubFromContext(ctx); hub != nil {
		hub.WithScope(func(scope *sentry.Scope) {
			prefix := ""

			for _, goa := range l.goas {
				temporaryPrefix := prefix
				if goa.group != "" {
					temporaryPrefix = fmt.Sprintf("%s.", goa.group)
				}

				if len(goa.attrs) == 0 {
					prefix = temporaryPrefix
					continue
				}

				for _, attr := range goa.attrs {
					scope.SetTag(
						fmt.Sprintf("%s%s", temporaryPrefix, attr.Key),
						attr.Value.String(),
					)
				}
			}

			scope.SetLevel(sentry.LevelError)
			hub.CaptureException(err)
		})
	}
}

func recordToError(record slog.Record) error {
	return errors.New(record.Message)
}
