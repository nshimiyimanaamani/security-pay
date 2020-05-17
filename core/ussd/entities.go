package ussd

import (
	"encoding/json"
	"io"
	"time"

	validate "github.com/go-playground/validator/v10"
)

// Request ...
type Request struct {
	SessionID   string `json:"sessionId,omitempty"`
	GwRef       string `json:"gwRef,omitempty"`
	GwTstamp    string `json:"gwTstamp,omitempty"`
	Msisdn      string `json:"msisdn,omitempty"`
	NetworkCode string `json:"networkCode,omitempty"`
	ServiceCode string `json:"serviceCode,omitempty"`
	UserInput   string `json:"userInput,omitempty"`
	ServiceID   string `json:"serviceId,omitempty"`
	TenantID    string `json:"tenantId,omitempty"`
}

// FromJSON ...
func (req *Request) FromJSON(r io.Reader) error {
	decoder := json.NewDecoder(r)
	return decoder.Decode(req)
}

// Validate ...
func (req *Request) Validate() error {
	validator := validate.New()
	return validator.Struct(req)
}

// ToJSON ...
func (req *Request) ToJSON(w io.Writer) error {
	encoder := json.NewEncoder(w)
	return encoder.Encode(req)
}

// Response ...
type Response struct {
	SessionID  string    `json:"sessionId" validate:"required"`
	GatewayRef string    `json:"gwRef" validate:"required"`
	AppRef     string    `json:"appRef" validate:"required"`
	GwTstamp   time.Time `json:"gwTstamp,omitempty" validate:"required"`
	Text       string    `json:"text" validate:"required"`
	End        int       `json:"continueSession" validate:"required"`
}

// ToJSON ...
func (res *Response) ToJSON(w io.Writer) error {
	encoder := json.NewEncoder(w)
	return encoder.Encode(res)
}
