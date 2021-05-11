package modlog

import (
	"context"
	"io"

	"github.com/rs/zerolog"
)

type ContextHook func(evt *zerolog.Event, ctx context.Context) *zerolog.Event

type Config struct {
	Output          io.Writer
	TimeFieldFormat string
	ContextHooks    []ContextHook
}
