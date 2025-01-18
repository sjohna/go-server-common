package main

import (
	"context"
	"fmt"
	"github.com/sjohna/go-server-common/errors"
	"github.com/sjohna/go-server-common/log"
	"net/http"
	"runtime"
	"sync"
)

func main() {
	logger, configLogger := log.GetApplicationLoggers("testOutput/serverCommonTest", "go-server-common-test")

	logger.Info("Hello world")
	configLogger.Info("Hello world")

	logger.WithField("test", "test").Panic("Test primitive field")
	logger.WithField("test", struct{ Test string }{"test"}).Panic("Test struct field")
	logger.WithField("test", map[string]string{"test": "test"}).Panic("Test map field")
	logger.WithField("test", []string{"test", "test"}).Panic("Test slice field")

	testContext := context.WithValue(context.Background(), "logger", logger)

	http.NewRequestWithContext(testContext, "GET", "http://localhost:8080", nil)

	testStackTrace()

	err := testError()
	logger.WithError(err).Error("something bad happened")

	err2 := testErrorInGoRoutine()
	logger.WithError(err2).Error("something else bad happened")

	err3 := testQueryError()
	logger.WithError(err3).Error("something bad happened")
}

func testStackTrace() {
	pcs := make([]uintptr, 100)

	callers := runtime.Callers(1, pcs)

	fmt.Printf("caller count: %d\n", callers)

	for i := 0; i < callers; i++ {
		fmt.Println(pcs[i])
	}

	frames := runtime.CallersFrames(pcs)

	for {
		frame, more := frames.Next()
		fmt.Printf("\t%s:%d %s\n", frame.File, frame.Line, frame.Function)

		if !more {
			break
		}
	}
}

func testError() errors.Error {
	err := fmt.Errorf("test error")

	return errors.Wrap(err, "wrapped error")
}

func testErrorInGoRoutine() errors.Error {
	wg := sync.WaitGroup{}

	wg.Add(1)

	var err errors.Error

	go func() {
		defer wg.Done()
		err = errors.New("error from goroutine")
	}()

	wg.Wait()
	return err
}

func testQueryError() errors.Error {
	return errors.WrapQueryError(nil, "query error message", "SELECT * FROM some_table WHERE col = $1", "fred")
}
