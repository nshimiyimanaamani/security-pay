package ussd

import (
	"encoding/json"
	"io"
	"time"

	validate "github.com/go-playground/validator/v10"
	"github.com/rugwirobaker/paypack-backend/pkg/errors"
)

// Request ...
type Request struct {
	SessionID   string `json:"sessionId,omitempty" validate:"required"`
	GwRef       string `json:"gwRef,omitempty" validate:"required"`
	GwTstamp    string `json:"gwTstamp,omitempty"`
	Msisdn      string `json:"msisdn,omitempty" validate:"required"`
	NetworkCode string `json:"networkCode,omitempty" validate:"required"`
	ServiceCode string `json:"serviceCode,omitempty" validate:"required"`
	UserInput   string `json:"userInput,omitempty" validate:"required"`
	ServiceID   string `json:"serviceId,omitempty" validate:"required"`
	TenantID    string `json:"tenantId,omitempty" validate:"required"`
}

// FromJSON ...
func (req *Request) FromJSON(r io.Reader) error {
	decoder := json.NewDecoder(r)
	return decoder.Decode(req)
}

// Validate ...
func (req *Request) Validate() error {
	const op errors.Op = "core/ussd/Request.Validate"

	var validator = validate.New()

	if err := validator.Struct(req); err != nil {
		return errors.E(op, err, errors.KindBadRequest)
	}
	return nil
}

// ToJSON ...
func (req *Request) ToJSON(w io.Writer) error {
	encoder := json.NewEncoder(w)
	return encoder.Encode(req)
}

// Response ...
type Response struct {
	SessionID string    `json:"sessionId" validate:"required"`
	GwRef     string    `json:"gwRef" validate:"required"`
	AppRef    string    `json:"appRef" validate:"required"`
	GwTstamp  time.Time `json:"gwTstamp,omitempty" validate:"required"`
	Text      string    `json:"text" validate:"required"`
	End       int       `json:"continueSession" validate:"required"`
}

func (res Response) String() string {
	return res.Text
}

// Tail ...
func (res Response) Tail() bool {
	if res.End == 1 {
		return true
	}
	if res.End == 0 {
		return false
	}
	return false
}

// ToJSON ...
func (res *Response) ToJSON(w io.Writer) error {
	encoder := json.NewEncoder(w)
	return encoder.Encode(res)
}
