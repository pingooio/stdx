package slogx

import (
	"log/slog"

	"github.com/pingooio/stdx/timex"
)

func Err(err error) slog.Attr {
	return slog.String("error", err.Error())
}

func Time(key string, t timex.Time) slog.Attr {
	return slog.String(key, t.String())
}
