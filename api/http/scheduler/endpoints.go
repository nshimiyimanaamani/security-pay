package scheduler

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rugwirobaker/paypack-backend/api/http/encoding"
	"github.com/rugwirobaker/paypack-backend/core/scheduler"
	"github.com/rugwirobaker/paypack-backend/pkg/errors"
	"github.com/rugwirobaker/paypack-backend/pkg/log"
)

// Schedule task
func Schedule(lgger log.Entry, svc scheduler.Service) http.Handler {
	const op errors.Op = "api/http/scheduler/Enqueue"

	f := func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		var task = vars["task"]

		err := svc.Schedule(r.Context(), task)
		if err != nil {
			err = errors.E(op, err)
			lgger.SystemErr(err)
			encoding.EncodeError(w, errors.Kind(err), err)
			return
		}

		res := map[string]string{"message": fmt.Sprintf("%s: is pending", task)}

		if err := encoding.Encode(w, http.StatusOK, res); err != nil {
			err = errors.E(op, err)
			lgger.SystemErr(err)
			encoding.EncodeError(w, errors.Kind(err), err)
			return
		}
	}
	return http.HandlerFunc(f)
}
