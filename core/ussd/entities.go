package ussd

import (
	"encoding/json"
	"io"
	"time"

	validate "github.com/go-playground/validator/v10"
)

// Request ...
type Request struct {
	//ussd session as passed by the telecom.
	Session string `validate:"required" json:"sessionId"`

	// ussd code and where relevant, with your service extension
	ServiceCode string `validate:"required" json:"serviceCode"`

	// identifies of the mobile subscriber interacting with your ussd application
	NetworkCode string `validate:"required" json:"networkCode"`

	// unique USSD Gateway reference generated when the session starts
	GatewayRef string `validate:"required" json:"gwRef"`

	// is the number of the mobile subscriber interacting with your ussd application.
	MSISDN string `validate:"required" json:"msisdn"`

	// ussd Gateway Timestamp
	Timestamp string `validate:"required" json:"gwTstamp"`
	// user input.

	UserInput string `validate:"required" json:"userInput"`

	// your application/service unique identifier
	ServiceID string `validate:"required" json:"serviceId"`

	// uniquely identifies your business account with the relevant KYC info
	TenantID string `validate:"required" json:"tenantId"`
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
	//ussd session as passed by the telecom.
	Session string `validate:"required" json:"sessionId"`

	// unique USSD Gateway reference generated when the session starts
	GatewayRef string `validate:"required" json:"gwRef"`

	// Your unique app/service session reference as captured in your app or service.
	AppRef string `validate:"required" json:"appRef"`

	//Your unique app /service timestamp as captured in your app or service
	Timestamp time.Time

	//The text to be passed back to the telco formatted with “\n” for new  lines
	Text string `validate:"required" json:"text"`

	//Indicats whether we are the end of our session
	End int `validate:"required" json:"continueSession"`
}
