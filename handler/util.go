package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/sirupsen/logrus"
)

func HandlerLogger(r *http.Request, handler string) *logrus.Entry {
	log := r.Context().Value("logger").(*logrus.Entry).WithField("handler", handler)
	log.Info("Handler called")
	return log
}

func LogHandlerReturn(log *logrus.Entry) {
	log.Info("Handler returned")
}

func UnmarshalRequestBody(log *logrus.Entry, r *http.Request, value interface{}) error {
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
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

func RespondError(log *logrus.Entry, w http.ResponseWriter, err error, httpResponseCode int) {
	log.WithError(err).WithField("httpResponseCode", httpResponseCode).Error()
	http.Error(w, err.Error(), httpResponseCode)
}

func RespondInternalServerError(log *logrus.Entry, w http.ResponseWriter, err error) {
	RespondError(log, w, err, 500)
}

func RespondClientError(log *logrus.Entry, w http.ResponseWriter, err error) {
	RespondError(log, w, err, 400)
}

func RespondJSON(log *logrus.Entry, w http.ResponseWriter, value interface{}) {
	w.Header().Set("Content-Type", "application/json")

	jsonResp, err := json.Marshal(value)
	if err != nil {
		RespondError(log, w, err, 500)
		return
	}

	written, err := w.Write(jsonResp)
	if err != nil {
		log.WithError(err).Error("Error writing response")
	} else {
		log.WithField("responseBytes", written).Info("Respond success")
	}
}
