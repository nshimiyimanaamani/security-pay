package transactions

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	transport "github.com/rugwirobaker/paypack-backend/api/http"
	"github.com/rugwirobaker/paypack-backend/app/transactions"
	"github.com/rugwirobaker/paypack-backend/models"
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

		//router.HandleFunc("/{property}", func(w http.ResponseWriter, r *http.Request) {}).Methods("GET")

		// router.HandleFunc("/{method}", func(w http.ResponseWriter, r *http.Request) {}).Methods("GET")

		// router.HandleFunc("/{month}", func(w http.ResponseWriter, r *http.Request) {}).Methods("GET")

		// router.HandleFunc("/{year}", func(w http.ResponseWriter, r *http.Request) {}).Methods("GET")
	}
	return handler
}

func handleRecordTransaction(svc transactions.Service, w http.ResponseWriter, r *http.Request) {
	var err error

	if err = transport.CheckContentType(r); err != nil {
		transport.EncodeError(w, err)
		return
	}

	var transaction models.Transaction
	if err = json.NewDecoder(r.Body).Decode(&transaction); err != nil {
		transport.EncodeError(w, err)
		return
	}

	if err = transaction.Validate(); err != nil {
		transport.EncodeError(w, err)
		return
	}

	transaction, err = svc.RecordTransaction(transaction)
	if err != nil {
		transport.EncodeError(w, err)
		return
	}

	response := recordTransRes{
		ID: transaction.ID,
	}
	if err = transport.EncodeResponse(w, response); err != nil {
		transport.EncodeError(w, err)
		return
	}
}

func handleViewTransaction(svc transactions.Service, w http.ResponseWriter, r *http.Request) {
	var err error

	vars := mux.Vars(r)

	id := vars["id"]

	transaction, err := svc.ViewTransaction(id)
	if err != nil {
		transport.EncodeError(w, err)
		return
	}

	response := viewTransRes{
		ID:       transaction.ID,
		Property: transaction.Property,
		Amount:   transaction.Amount,
		Method:   transaction.Method,
	}

	if err = transport.EncodeResponse(w, response); err != nil {
		transport.EncodeError(w, err)
		return
	}
}

func handleListTransaction(svc transactions.Service, w http.ResponseWriter, r *http.Request) {
	var err error

	vars := mux.Vars(r)

	offset, err := strconv.ParseUint(vars["offset"], 10, 32)
	if err != nil {
		transport.EncodeError(w, err)
		return
	}

	limit, err := strconv.ParseUint(vars["limit"], 10, 32)
	if err != nil {
		transport.EncodeError(w, err)
		return
	}

	// if offset == 0 || limit == 0 {
	// 	transport.EncodeError(w, models.ErrInvalidEntity)
	// 	return
	// }

	var page models.TransactionPage

	page, err = svc.ListTransactions(offset, limit)
	if err != nil {
		transport.EncodeError(w, err)
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
			ID:       transaction.ID,
			Property: transaction.Property,
			Amount:   transaction.Amount,
			Method:   transaction.Method,
		}
		response.Transactions = append(response.Transactions, view)
	}

	if err = transport.EncodeResponse(w, response); err != nil {
		transport.EncodeError(w, err)
		return
	}

}
