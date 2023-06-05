package logger

import (
	"io"
	"os"
	"sync"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

var once sync.Once

var log zerolog.Logger

func Get() *zerolog.Logger {
	once.Do(func() {
		zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

		var output io.Writer = zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		}

		log = zerolog.New(output).
			Level(zerolog.TraceLevel).
			With().
			Timestamp().
			Caller().
			Logger()
	})

	return &log
}

func LogErr(err error, logLevel ...zerolog.Level) {
	if err == nil {
		return
	}

	lvl := zerolog.WarnLevel

	if err != nil {
		lvl = zerolog.ErrorLevel
	}

	if len(logLevel) > 0 {
		lvl = logLevel[0]
	}

	Get().WithLevel(lvl).Caller(1).Msg(err.Error())
}
