package nova

// import (
// 	"errors"
// 	"fmt"
// 	"log"
// 	"net/http"
// 	"os"
// 	"time"

// 	"github.com/rugwirobaker/paypack-backend/app/payment"
// )

// // Timeout sets the default client timeout
// var Timeout = 30 * time.Second

// // ErrBackendError is returned for all unexpected errors from 3-third party api
// var ErrBackendError = errors.New("backend returned unexpected error")

// type gateway struct {
// 	endpoint string
// 	token    string
// 	client   http.Client
// }

// // New instantiates the novapay api gateway
// func New(cfg *Config) *gateway {
// 	if cfg.Endpoint == "" || cfg.Token == "" {
// 		log.Println(fmt.Sprintf("could not create gateway: unacceptable nova configuration"))
// 		os.Exit(1)
// 	}
// 	var timeout = Timeout

// 	if cfg.TimeOut != 0 {
// 		timeout = cfg.TimeOut
// 	}
// 	client := http.Client{Timeout: timeout}
// 	return &gateway{
// 		client:   client,
// 		endpoint: cfg.Endpoint,
// 		token:    cfg.Token,
// 	}
// }

// func (gw *gateway) Initiate(py payment.Transaction) (string, error) {
// 	//empty := payment.Message{}

// 	// body := PaymentToRequest(gw.token, py)

// 	// bb, err := json.Marshal(body)
// 	// if err != nil {
// 	// 	return empty, err
// 	// }
// 	// request, err := http.NewRequest("POST", gw.endpoint, bytes.NewReader(bb))

// 	// request.Header.Set("Content-type", "application/json")

// 	// if err != nil {
// 	// 	return empty, err
// 	// }
// 	// res, err := gw.client.Do(request)
// 	// if err != nil {
// 	// 	return empty, err
// 	// }
// 	// if res.StatusCode != http.StatusOK {
// 	// 	return empty, ErrBackendError
// 	// }
// 	// message := payment.Message{}
// 	// if err := json.NewDecoder(res.Body).Decode(&message); err != nil {
// 	// 	return empty, err
// 	// }
// 	// return message, nil
// }

// func (gw *gateway) Validate(r payment.Validation) payment.Validation {
// 	// res := payment.Validation{Token: gw.token, Ref: r.ExternalTransactionsID}
// 	// return res
// }
