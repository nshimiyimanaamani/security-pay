package transactions

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/rugwirobaker/paypack-backend/app/transactions"
	"github.com/rugwirobaker/paypack-backend/pkg/log"
)

// Record handles transaction record
func Record(logger log.Entry, svc transactions.Service) http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {

		if err := CheckContentType(r); err != nil {
			EncodeError(w, err)
			return
		}

		var transaction transactions.Transaction
		if err := json.NewDecoder(r.Body).Decode(&transaction); err != nil {
			EncodeError(w, err)
			return
		}

		if err := transaction.Validate(); err != nil {
			EncodeError(w, err)
			return
		}

		transaction, err := svc.Record(r.Context(), transaction)
		if err != nil {
			EncodeError(w, err)
			return
		}

		response := recordTransRes{
			ID: transaction.ID,
		}
		if err = EncodeResponse(w, response); err != nil {
			EncodeError(w, err)
			return
		}
	}
	return http.HandlerFunc(f)
}

// Retrieve handles transaction retrieval
func Retrieve(logger log.Entry, svc transactions.Service) http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		transaction, err := svc.Retrieve(r.Context(), id)
		if err != nil {
			EncodeError(w, err)
			return
		}

		response := viewTransRes{
			ID:           transaction.ID,
			Property:     transaction.MadeFor,
			Owner:        transaction.MadeBy,
			Amount:       transaction.Amount,
			Method:       transaction.Method,
			Address:      transaction.Address,
			DateRecorded: transaction.DateRecorded,
		}

		if err = EncodeResponse(w, response); err != nil {
			EncodeError(w, err)
			return
		}
	}
	return http.HandlerFunc(f)
}

// List handles transaction list
func List(logger log.Entry, svc transactions.Service) http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)

		offset, err := strconv.ParseUint(vars["offset"], 10, 32)
		if err != nil {
			EncodeError(w, err)
			return
		}

		limit, err := strconv.ParseUint(vars["limit"], 10, 32)
		if err != nil {
			EncodeError(w, err)
			return
		}

		var page transactions.TransactionPage

		page, err = svc.List(r.Context(), offset, limit)
		if err != nil {
			EncodeError(w, err)
			return
		}

		response := transPageRes{
			pageRes: pageRes{
				Total:  page.Total,
				Offset: page.Offset,
				Limit:  page.Limit,
			},
			Transactions: []viewTransRes{},
		}

		for _, transaction := range page.Transactions {
			view := viewTransRes{
				ID:           transaction.ID,
				Property:     transaction.MadeFor,
				Owner:        transaction.MadeBy,
				Amount:       transaction.Amount,
				Method:       transaction.Method,
				Address:      transaction.Address,
				DateRecorded: transaction.DateRecorded,
			}
			response.Transactions = append(response.Transactions, view)
		}

		if err = EncodeResponse(w, response); err != nil {
			EncodeError(w, err)
			return
		}
	}
	return http.HandlerFunc(f)
}

// ListByProperty handles transactions retrieval given house code
func ListByProperty(logger log.Entry, svc transactions.Service) http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {}

	return http.HandlerFunc(f)
}

// ListByOwner handles transactions retrieval given house code
func ListByOwner(logger log.Entry, svc transactions.Service) http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {}

	return http.HandlerFunc(f)
}
