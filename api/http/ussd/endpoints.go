package ussd

import (
	"io/ioutil"
	lg "log"
	"net/http"

	"github.com/rugwirobaker/paypack-backend/api/http/encoding"
	"github.com/rugwirobaker/paypack-backend/core/identity/uuid"
	"github.com/rugwirobaker/paypack-backend/core/ussd"
	"github.com/rugwirobaker/paypack-backend/pkg/errors"
	"github.com/rugwirobaker/paypack-backend/pkg/log"
)

// Process recieves ussd callback requests
func Process(lgger log.Entry, svc ussd.Service) http.Handler {
	const op errors.Op = "api/http/ussd/Process"

	f := func(w http.ResponseWriter, r *http.Request) {

		// debug(r, w)

		req := &ussd.Request{
			GwRef:    uuid.New().ID(),
			TenantID: "paypack",
			// ServiceID: uuid.New().ID(),
		}

		err := encoding.Decode(r, &req)
		if err != nil {
			err := errors.E(op, err)
			lgger.SystemErr(err)
			encoding.EncodeError(w, errors.Kind(err), err)
			return
		}

		res, err := svc.Process(r.Context(), req)
		if err != nil {
			lgger.SystemErr(err)
			encoding.Encode(w, http.StatusOK, res)
			return
		}

		if err := encoding.Encode(w, http.StatusOK, res); err != nil {
			lgger.SystemErr(err)
			encoding.EncodeError(w, errors.Kind(err), err)
			return
		}
	}
	return http.HandlerFunc(f)
}

func debug(r *http.Request, w http.ResponseWriter) {
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		lg.Print("err ", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	lg.Printf("%v", string(buf))
}
