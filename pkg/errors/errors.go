package errors

import (
	"errors"
	"fmt"
	"net/http"
	"runtime"

	"github.com/sirupsen/logrus"
)

// Kind enums
const (
	KindNotFound       = http.StatusNotFound
	KindBadRequest     = http.StatusBadRequest
	KindUnexpected     = http.StatusInternalServerError
	KindAlreadyExists  = http.StatusConflict
	KindRateLimit      = http.StatusTooManyRequests
	KindNotImplemented = http.StatusNotImplemented
	KindRedirect       = http.StatusMovedPermanently
)

var _ (error) = (*Error)(nil)

// Error is a Larissa system error.
// It carries information and behavior
// as to what caused this error so that
// callers can implement logic around it.
type Error struct {
	Kind     int
	Op       Op
	Entity   Ent
	Err      error
	Severity logrus.Level
}

func (e Error) Error() string {
	return e.Err.Error()
}

// Is is a shorthand for checking an error against a kind.
func Is(err error, kind int) bool {
	if err == nil {
		return false
	}
	return Kind(err) == kind
}

// Op describes any independent function or
// method in Larissa. A series of operations
// forms a more readable stack trace.
type Op string

func (o Op) String() string {
	return string(o)
}

// Ent represents the entity(model) in an Error
type Ent string

// E is a helper function to construct an Error type
// Operation always comes first, module path and version
// come second, they are optional. Args must have at least
// an error or a string to describe what exactly went wrong.
// You can optionally pass a Logrus severity to indicate
// the log level of an error based on the context it was constructed in.
func E(op Op, args ...interface{}) error {
	e := Error{Op: op}
	if len(args) == 0 {
		msg := "errors.E called with 0 args"
		_, file, line, ok := runtime.Caller(1)
		if ok {
			msg = fmt.Sprintf("%v - %v:%v", msg, file, line)
		}
		e.Err = errors.New(msg)
	}
	for _, a := range args {
		switch a := a.(type) {
		case error:
			e.Err = a
		case string:
			e.Err = errors.New(a)
		case Ent:
			e.Entity = a
		case logrus.Level:
			e.Severity = a
		case int:
			e.Kind = a
		}
	}
	if e.Err == nil {
		e.Err = errors.New(KindText(e))
	}
	return e
}

// Severity returns the log level of an error
// if none exists, then the level is Error because
// it is an unexpected.
func Severity(err error) logrus.Level {
	e, ok := err.(Error)
	if !ok {
		return logrus.ErrorLevel
	}

	// if there's no severity (0 is Panic level in logrus
	// which we should not use since cloud providers only have
	// debug, info, warn, and error) then look for the
	// child's severity.
	if e.Severity < logrus.ErrorLevel {
		return Severity(e.Err)
	}

	return e.Severity
}

// Expect is a helper that returns an Info level
// if the error has the expected kind, otherwise
// it returns an Error level.
func Expect(err error, kinds ...int) logrus.Level {
	for _, kind := range kinds {
		if Kind(err) == kind {
			return logrus.InfoLevel
		}
	}
	return logrus.ErrorLevel
}

// Kind recursively searches for the
// first error kind it finds.
func Kind(err error) int {
	e, ok := err.(Error)
	if !ok {
		return KindUnexpected
	}

	if e.Kind != 0 {
		return e.Kind
	}

	return Kind(e.Err)
}

// KindText returns a friendly string
// of the Kind type. Since we use http
// status codes to represent error kinds,
// this method just defers to the net/http
// text representations of statuses.
func KindText(err error) string {
	return http.StatusText(Kind(err))
}

// Ops aggregates the error's operation
// with all the embedded errors' operations.
// This way you can construct a queryable
// stack trace.
func Ops(err Error) []Op {
	ops := []Op{err.Op}
	for {
		embeddedErr, ok := err.Err.(Error)
		if !ok {
			break
		}

		ops = append(ops, embeddedErr.Op)
		err = embeddedErr
	}

	return ops
}
