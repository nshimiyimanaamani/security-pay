package transactions

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/rugwirobaker/paypack-backend/app/transactions"
)

// MakeAdapter takes a transaction service instance and returns a http handler
func MakeAdapter(router *mux.Router) func(svc transactions.Service) {
	handler := func(svc transactions.Service) {

		router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			handleRecordTransaction(svc, w, r)
		}).Methods("POST")

		router.HandleFunc("/{id}", func(w http.ResponseWriter, r *http.Request) {
			handleViewTransaction(svc, w, r)
		}).Methods("GET")

		router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			handleListTransaction(svc, w, r)
		}).Queries("offset", "{offset}", "limit", "{limit}").Methods("GET")

		// router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		// }).Queries("property", "{property}", "offset", "{offset}", "limit", "{limit}").Methods("GET")

		// router.HandleFunc("/{method}", func(w http.ResponseWriter, r *http.Request) {}).Methods("GET")

		// router.HandleFunc("/{month}", func(w http.ResponseWriter, r *http.Request) {}).Methods("GET")

		// router.HandleFunc("/{year}", func(w http.ResponseWriter, r *http.Request) {}).Methods("GET")
	}
	return handler
}

func handleRecordTransaction(svc transactions.Service, w http.ResponseWriter, r *http.Request) {
	var err error

	if err = CheckContentType(r); err != nil {
		EncodeError(w, err)
		return
	}

	var transaction transactions.Transaction
	if err = json.NewDecoder(r.Body).Decode(&transaction); err != nil {
		EncodeError(w, err)
		return
	}

	if err = transaction.Validate(); err != nil {
		EncodeError(w, err)
		return
	}

	token := r.Header.Get("Authorization")

	transaction, err = svc.RecordTransaction(token, transaction)
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

func handleViewTransaction(svc transactions.Service, w http.ResponseWriter, r *http.Request) {
	var err error

	vars := mux.Vars(r)

	id := vars["id"]
	token := r.Header.Get("Authorization")

	transaction, err := svc.ViewTransaction(token, id)
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

func handleListTransaction(svc transactions.Service, w http.ResponseWriter, r *http.Request) {
	var err error

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

	token := r.Header.Get("Authorization")

	page, err = svc.ListTransactions(token, offset, limit)
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
