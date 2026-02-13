package logging

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func Init(env string) zerolog.Logger {
	if env == "development" {
		return zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().Logger()
	}
	return log.Logger.Output(os.Stderr).With().Timestamp().Logger()
}

func WithRequestID(l zerolog.Logger, requestID string) zerolog.Logger {
	return l.With().Str("request_id", requestID).Logger()
}
