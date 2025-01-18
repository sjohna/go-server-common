package handler

import (
	"context"
	"github.com/sjohna/go-server-common/errors"
	"github.com/sjohna/go-server-common/log"
	"net/http"
)

type HandlerFunc func(ctx context.Context, r *http.Request) (interface{}, errors.Error)

func Handler(handler HandlerFunc) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		ret, err := handler(ctx, r)

		if err != nil {
			if err.Warning() {
				log.Ctx(ctx).WithError(err).Warn("Error returned from handler func")
			} else {
				log.Ctx(ctx).WithError(err).Error("Error returned from handler func")
			}

			httpResponseCode := http.StatusBadRequest
			if err.Internal() {
				httpResponseCode = http.StatusInternalServerError
			}

			http.Error(w, err.Error(), httpResponseCode)
			return
		}

		if ret != nil {
			err := RespondJSON(ctx, w, ret)
			if err != nil {
				log.Ctx(ctx).WithError(err).Error("Error writing JSON response to handler!!!!")
			}

			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
