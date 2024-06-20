package slwrap

import "log/slog"

func Wrap(err error) slog.Attr {
	return slog.Attr{
		Key:   "Error",
		Value: slog.StringValue(err.Error()),
	}
}
