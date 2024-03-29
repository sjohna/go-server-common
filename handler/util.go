package handler

import (
	"context"
	"encoding/json"
	"github.com/sjohna/go-server-common/log"
	"io/ioutil"
	"net/http"
)

func HandlerContext(r *http.Request, handler string) (context.Context, log.Logger) {
	logger := r.Context().Value("logger").(log.Logger).WithField("handler", handler)
	logger.Trace("Handler called")
	ctx := context.WithValue(r.Context(), "logger", logger)
	return ctx, logger
}

func LogHandlerReturn(logger log.Logger) {
	logger.Trace("Handler returned")
}

func UnmarshalRequestBody(logger log.Logger, r *http.Request, value interface{}) error {
	body, err := ioutil.ReadAll(r.Body)
	defer func() {
		err := r.Body.Close()
		if err != nil {
			logger.WithError(err).Error("UnmarshalRequestBody: Failed to close request body")
		}
	}()
	if err != nil {
		return err
	}

	// Unmarshal
	err = json.Unmarshal(body, value)
	if err != nil {
		return err
	}

	return nil
}

func RespondError(logger log.Logger, w http.ResponseWriter, err error, httpResponseCode int) {
	logger.WithError(err).WithField("httpResponseCode", httpResponseCode).Error("Handler RespondError")
	http.Error(w, err.Error(), httpResponseCode)
}

func RespondInternalServerError(logger log.Logger, w http.ResponseWriter, err error) {
	logger.WithError(err).WithField("httpResponseCode", http.StatusInternalServerError).Error("Handler RespondInternalServerError")
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

func RespondClientError(logger log.Logger, w http.ResponseWriter, err error) {
	logger.WithError(err).WithField("httpResponseCode", http.StatusBadRequest).Warn("Handler RespondClientError")
	http.Error(w, err.Error(), http.StatusBadRequest)
}

func RespondJSON(logger log.Logger, w http.ResponseWriter, value interface{}) {
	w.Header().Set("Content-Type", "application/json")

	jsonResp, err := json.Marshal(value)
	if err != nil {
		logger.WithError(err).Error("Error marshalling response")
		RespondInternalServerError(logger, w, err)
		return
	}

	written, err := w.Write(jsonResp)
	if err != nil {
		logger.WithError(err).Error("Error writing response")
	} else {
		logger.WithField("responseBytes", written).Debug("Respond success")
	}
}
