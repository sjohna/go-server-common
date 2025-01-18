package log

import "context"

var General Logger
var Config Logger

// TODO: init these here
func SetGlobalLoggers(general, config Logger) {
	General = general
	Config = config
}

func Ctx(ctx context.Context) Logger {
	if ctx.Value("logger").(Logger) != nil {
		return ctx.Value("logger").(Logger)
	}

	return General
}
