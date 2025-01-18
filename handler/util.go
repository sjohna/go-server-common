package handler

import (
	"context"
	"encoding/json"
	"github.com/sjohna/go-server-common/errors"
	"github.com/sjohna/go-server-common/log"
	"io"
	"net/http"
)

func UnmarshalRequestBody(ctx context.Context, r *http.Request, value interface{}) errors.Error {
	body, err := io.ReadAll(r.Body)
	defer func() {
		err := r.Body.Close()
		if err != nil {
			myErr := errors.WrapInputError(err, "Error closing request body")
			log.Ctx(ctx).WithError(myErr).Error("UnmarshalRequestBody: Failed to close request body")
		}
	}()
	if err != nil {
		return errors.Wrap(err, "Failed to read request body")
	}

	// Unmarshal
	err = json.Unmarshal(body, value)
	if err != nil {
		return errors.WrapInputError(err, "Failed to unmarshal request body")
	}

	return nil
}

func RespondJSON(ctx context.Context, w http.ResponseWriter, value interface{}) errors.Error {
	w.Header().Set("Content-Type", "application/json")

	jsonResp, err := json.Marshal(value)
	if err != nil {
		myErr := errors.Wrap(err, "Error marshalling JSON response")
		return myErr
	}

	_, err = w.Write(jsonResp)
	if err != nil {
		myErr := errors.Wrap(err, "Error writing JSON response")
		return myErr
	}

	return nil
}
