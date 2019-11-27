package feedback

import (
	"net/http"

	"github.com/rugwirobaker/paypack-backend/pkg/log"

	"github.com/gorilla/mux"
	"github.com/rugwirobaker/paypack-backend/app/feedback"
)

// ProtocolHandler adapts the feedback service into an http.handler
type ProtocolHandler func(lgger log.Entry, svc feedback.Service) http.Handler

// HandlerOpts are the generic options
// for a ProtocolHandler
type HandlerOpts struct {
	Service feedback.Service
	Logger  *log.Logger
}

// Recode handlers new feedback message submission
func Recode(lgger log.Entry, svc feedback.Service) http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		if err := CheckContentType(r); err != nil {
			EncodeError(w, err)
			return
		}

		var req feedback.Message

		err := decode(r.Body, &req)
		if err != nil {
			EncodeError(w, err)
			return
		}
		res, err := svc.Record(ctx, &req)
		if err != nil {
			EncodeError(w, err)
			return
		}
		encode(w, RecordRes{ID: res.ID})
	}

	return http.HandlerFunc(f)
}

// Update handles feedback updates
func Update(lgger log.Entry, svc feedback.Service) http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		if err := CheckContentType(r); err != nil {
			EncodeError(w, err)
			return
		}

		var req feedback.Message

		err := decode(r.Body, &req)
		if err != nil {
			EncodeError(w, err)
			return
		}

		vars := mux.Vars(r)

		req.ID = vars["id"]

		if err := svc.Update(ctx, req); err != nil {
			EncodeError(w, err)
			return
		}
		encode(w, map[string]string{"message": "message updated"})
	}

	return http.HandlerFunc(f)
}

// Retrieve handles feedback entry retrieval
func Retrieve(lgger log.Entry, svc feedback.Service) http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		vars := mux.Vars(r)
		id := vars["id"]

		req, err := svc.Retrieve(ctx, id)
		if err != nil {
			EncodeError(w, err)
			return
		}
		res := RetrieveRes{
			ID:        req.ID,
			Title:     req.Title,
			Body:      req.Body,
			Creator:   req.Creator,
			CreatedAt: req.CreatedAt,
			UpdatedAt: req.UpdatedAt,
		}

		encode(w, res)
	}

	return http.HandlerFunc(f)
}

// Delete handles feedback entry delete
func Delete(lgger log.Entry, svc feedback.Service) http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		vars := mux.Vars(r)
		id := vars["id"]

		err := svc.Delete(ctx, id)
		if err != nil {
			EncodeError(w, err)
			return
		}
		encode(w, map[string]string{"message": "message deleted"})
	}

	return http.HandlerFunc(f)
}
