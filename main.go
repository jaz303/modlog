package modlog

import (
	"errors"
	"io"
	"os"
	"sync"

	zl "github.com/rs/zerolog"
	zll "github.com/rs/zerolog/log"
)

var global = struct {
	lock         sync.Mutex
	allLoggers   []*Logger
	configured   bool
	output       io.Writer
	contextHooks []ContextHook
}{}

func ModuleLogger(name string) *Logger {
	innerLogger := zll.With().Str(keyModule, name).Logger()

	global.lock.Lock()
	defer global.lock.Unlock()

	if global.configured {
		innerLogger = configureLogger(innerLogger)
	}

	logger := &Logger{
		module: name,
		logger: innerLogger,
	}
	global.allLoggers = append(global.allLoggers, logger)

	return logger
}

func GetAllLogLevels() map[string]int {
	out := make(map[string]int)
	for _, l := range global.allLoggers {
		out[l.module] = int(l.logger.GetLevel())
	}
	return out
}

func SetAllLogLevels(levels map[string]int) {
	for _, l := range global.allLoggers {
		newLevel, ok := levels[l.module]
		if ok {
			l.SetLevel(zl.Level(newLevel))
		}
	}
}

func Configure(cfg *Config) {
	global.lock.Lock()
	defer global.lock.Unlock()

	if global.configured {
		panic(errors.New("modlog.Configure() must only be called once"))
	}

	if cfg.Output != nil {
		global.output = cfg.Output
	} else {
		global.output = zl.ConsoleWriter{Out: os.Stderr}
	}

	if len(cfg.TimeFieldFormat) > 0 {
		zl.TimeFieldFormat = cfg.TimeFieldFormat
	}

	for _, logger := range global.allLoggers {
		logger.logger = configureLogger(logger.logger)
	}

	global.contextHooks = cfg.ContextHooks
	global.configured = true
}

func configureLogger(logger zl.Logger) zl.Logger {
	return logger.Output(global.output)
}
