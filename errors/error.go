package errors

import "github.com/pkg/errors"

type StackFrame struct {
	File     string
	FileLine int
	Function string
}

type Severity int

const (
	SeverityError   = 0
	SeverityWarning = 1
)

type Origin int

const (
	OriginInput       = 0 // cause of error is in user input
	OriginApplication = 1 // cause of error is in this application
	OriginThirdParty  = 2 // cause of error is in another application or API, or the interface to the aforementioned
)

type Error interface {
	Error() string
	Internal() bool
	Warning() bool
	Is(error) bool
}

type ApplicationError struct {
	Severity   Severity
	Origin     Origin
	Message    string
	Inner      error
	StackTrace []StackFrame
}

func (e *ApplicationError) Error() string {
	return e.Message
}

func (e *ApplicationError) Internal() bool {
	return e.Origin == OriginApplication || e.Origin == OriginThirdParty
}

func (e *ApplicationError) Warning() bool {
	return e.Severity == SeverityWarning
}

func (e *ApplicationError) Is(err error) bool {
	if e.Inner == nil || err == nil {
		return false
	}

	return errors.Is(e.Inner, err)
}

type QueryError struct {
	ApplicationError
	Query string
	Args  []interface{}
}

func OriginString(origin Origin) string {
	switch origin {
	case OriginInput:
		return "input"
	case OriginApplication:
		return "application"
	case OriginThirdParty:
		return "third-party"
	}

	return "unknown"
}
