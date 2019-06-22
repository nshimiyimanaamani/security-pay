package sms

import "github.com/rugwirobaker/paypack-backend/models"

//Service defines the send sms API
type Service interface {
	//Send sends a message and returns nil if the operation is successful
	Send(models.Message) error
}

var _ (Service) = (*smsService)(nil)

type smsService struct{}

//New instanciates a new Service instance
func New() Service {
	return &smsService{}
}

func (svc *smsService) Send(message models.Message) error {
	return message.Validate()
}
