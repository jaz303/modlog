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
	lock            sync.Mutex
	allLoggers      []*Logger
	configured      bool
	output          io.Writer
	contextHooks    []ContextHook
	defaultLogLevel zl.Level
}{
	defaultLogLevel: zl.InfoLevel,
}

func ModuleLogger(name string) *Logger {
	global.lock.Lock()
	defer global.lock.Unlock()

	for _, l := range global.allLoggers {
		if l.module == name {
			return l
		}
	}

	innerLogger := zll.With().Str(keyModule, name).Logger()

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

func GetAllLogLevels() map[string]zl.Level {
	out := make(map[string]zl.Level)
	for _, l := range global.allLoggers {
		out[l.module] = l.logger.GetLevel()
	}
	return out
}

func SetAllLogLevels(levels map[string]zl.Level) {
	for _, l := range global.allLoggers {
		newLevel, ok := levels[l.module]
		if ok {
			l.SetLevel(newLevel)
		}
	}
}

func SetDefaultLogLevel(level zl.Level) {
	global.defaultLogLevel = level
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
	return logger.Level(global.defaultLogLevel).Output(global.output)
}
