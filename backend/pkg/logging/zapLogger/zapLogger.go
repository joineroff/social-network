package zapLogger

import (
	"github.com/joineroff/social-network/backend/pkg/logging"
	"go.uber.org/zap"
)

func New(cfg *Config, appName string) (logging.Logger, error) {
	config := cfg.withDefaults()

	builder, err := config.newBuilder()
	if err != nil {
		return nil, err
	}

	logger, err := builder.Build()
	if err != nil {
		return nil, err
	}

	return wrap(logger.With(zap.String("app_name", appName))), nil
}

func wrap(logger *zap.Logger) logging.Logger {
	return &zapLogger{logger: logger.Sugar()}
}

// Logger is a wrapper over zap logger. It is claimed to be thread safe.
type zapLogger struct {
	logger *zap.SugaredLogger
}

func (log *zapLogger) Debug(msg string, keysAndValues ...interface{}) {
	log.logger.Debugw(msg, keysAndValues...)
}

func (log *zapLogger) Info(msg string, keysAndValues ...interface{}) {
	log.logger.Infow(msg, keysAndValues...)
}

func (log *zapLogger) Warn(msg string, keysAndValues ...interface{}) {
	log.logger.Warnw(msg, keysAndValues...)
}

func (log *zapLogger) Error(msg string, keysAndValues ...interface{}) {
	log.logger.Errorw(msg, keysAndValues...)
}

func (log *zapLogger) With(keysAndValues ...interface{}) logging.Logger {
	l := log.logger.With(keysAndValues...)

	return wrap(l.Desugar())
}
