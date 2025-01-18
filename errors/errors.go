package errors

import (
	"context"
	"errors"
	"runtime"
)

// generally assume errors are application and error severity, and use other functions for when that's not the case

func Wrap(err error, message string) *ApplicationError {
	return &ApplicationError{
		SeverityError,
		OriginApplication,
		message,
		err,
		stackTrace(),
	}
}

func New(message string) *ApplicationError {
	return &ApplicationError{
		SeverityError,
		OriginApplication,
		message,
		nil,
		stackTrace(),
	}
}

func NewInput(message string) *ApplicationError {
	return &ApplicationError{
		SeverityWarning,
		OriginInput,
		message,
		nil,
		stackTrace(),
	}
}

func WrapDBError(err error, message string) *ApplicationError {
	severity := SeverityError
	if errors.Is(err, context.Canceled) {
		severity = SeverityWarning
	}

	return &ApplicationError{
		Severity(severity),
		OriginThirdParty,
		message,
		err,
		stackTrace(),
	}
}

func WrapQueryError(err error, message string, query string, args ...interface{}) *QueryError {
	severity := SeverityError
	if errors.Is(err, context.Canceled) {
		severity = SeverityWarning
	}

	return &QueryError{
		ApplicationError{
			Severity(severity),
			OriginThirdParty,
			message,
			err,
			stackTrace(),
		},
		query,
		args,
	}
}

func WrapInputError(err error, message string) *ApplicationError {
	return &ApplicationError{
		SeverityError,
		OriginInput,
		message,
		err,
		stackTrace(),
	}
}

func stackTrace() []StackFrame {
	pcs := make([]uintptr, 100)

	n := runtime.Callers(3, pcs) // skip Callers itself, stackTrace, and the caller in this file

	if n == 0 {
		return nil
	}

	frames := make([]StackFrame, 0, n)
	callersFrames := runtime.CallersFrames(pcs)

	if callersFrames == nil {
		return nil
	}

	for {
		frame, more := callersFrames.Next()

		frames = append(frames, StackFrame{
			frame.File,
			frame.Line,
			frame.Function,
		})

		if !more {
			break
		}
	}

	return frames
}
