package sms

import (
	"fmt"
	"time"

	validate "github.com/go-playground/validator/v10"
)

var errInvalidToken = fmt.Errorf("invalid authorization token")
var errMissingAuth = fmt.Errorf("missing authorization header")

type smsResponse struct {
	Success          bool   `validate:"required" json:"success"`
	Message          string `validate:"required" json:"message,omitempty"`
	Cost             int    `validate:"required" json:"cost,omitempty"`
	MessageReference string `validate:"required" json:"msgRef,omitempty"`
	GatewayReference string `validate:"required" json:"gatewayRef,omitempty"`
}

func (res *smsResponse) Validate() error {
	validator := validate.New()
	return validator.Struct(res)
}

func (res *smsResponse) RetrieveError() error {
	if res.Message != "" {
		switch res.Message {
		case "JWT Authorization Token incorrect or malformed":
			return errInvalidToken
		case "Missing Authorization Header":
			return errMissingAuth
		default:
			return fmt.Errorf("%s", res.Message)
		}
	}
	return nil
}

type authResponse struct {
	Success      bool      `json:"success" validate:"required"`
	Message      string    `json:"message,omitempty"`
	ExpiresAt    time.Time `json:"expires_at,omitempty"`
	AccessToken  string    `json:"access_token,omitempty"`
	RefreshToken string    `json:"refresh_token,omitempty"`
}

func (res *authResponse) Validate() error {
	validator := validate.New()
	return validator.Struct(res)
}

func (res *authResponse) Error() error {
	return nil
}

type authRequest struct {
	AppID     string `validate:"required" json:"api_username"`
	AppSecret string `validate:"required" json:"api_password"`
}

type refreshRequest struct {
	Token string `json:"refresh_token"`
}

type singleMSISDNRequest struct {
	MSISDN    string `json:"msisdn"`
	Message   string `json:"message"`
	SenderID  string `json:"sender_id"`
	Reference string `json:"msgRef"`
}

type bulkMSISDNRequest struct {
	MSISDN    []string `json:"msisdn"`    // a list of numbers to send message to
	Message   string   `json:"message"`   // the sms message content
	SenderID  string   `json:"sender_id"` // the id to display in the message
	Reference string   `json:"msgRef"`
}
