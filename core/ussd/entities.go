package ussd

import (
	"encoding/json"
	"io"

	validate "github.com/go-playground/validator/v10"
)

// SessionRequest ...
type SessionRequest struct {
	//ussd session as passed by the telecom.
	ID string `validate:"required" json:"sessionId"`

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

func (ses *SessionRequest) FromJSON(r io.Reader) error {
	decoder := json.NewDecoder(r)
	return decoder.Decode(ses)
}

func (ses *SessionRequest) Validate() error {
	validator := validate.New()
	return validator.Struct(ses)
}

func (ses *SessionRequest) ToJSON(w io.Writer) error {
	encoder := json.NewEncoder(w)
	return encoder.Encode(ses)
}

// type Command struct {
// 	screen   string
// 	level    string
// 	children []*Command
// }

// func ParseCmd(input string) error {
// 	if input == "" {
// 	}
// 	return nil
// }
