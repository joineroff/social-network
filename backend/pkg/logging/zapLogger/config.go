package zapLogger

import (
	"errors"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	timeEncoderEpoch = "epoch"

	defaultLevel       = infoLevel
	defaultTimeEncoder = timeEncoderEpoch

	debugLevel = "debug"
	errorLevel = "error"
	infoLevel  = "info"
	warnLevel  = "warn"
)

type Config struct {
	Debug       bool   `envconfig:"debug"`
	Level       string `envconfig:"level"`
	TimeEncoder string `envconfig:"time_encoder"`
}

func (c *Config) newBuilder() (*zap.Config, error) {
	var config zap.Config
	if c.Debug {
		config = zap.NewDevelopmentConfig()
	} else {
		config = zap.NewProductionConfig()
	}

	// Disable useless output
	config.DisableCaller = true

	var err error
	config.Level, err = c.getLevel()

	return &config, err
}

func (c *Config) getLevel() (lvl zap.AtomicLevel, err error) {
	switch c.Level {
	case debugLevel:
		lvl = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	case errorLevel:
		lvl = zap.NewAtomicLevelAt(zapcore.ErrorLevel)
	case infoLevel:
		lvl = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	case warnLevel:
		lvl = zap.NewAtomicLevelAt(zapcore.WarnLevel)
	default:
		err = errors.New("unknown log level")
	}
	return
}

func (c *Config) withDefaults() (config Config) {
	if c != nil {
		config = *c
	}

	if config.Level == "" {
		config.Level = defaultLevel
	}

	if config.TimeEncoder == "" {
		config.TimeEncoder = defaultTimeEncoder
	}

	return
}
