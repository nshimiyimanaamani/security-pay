package ussd

import "context"

type Service interface {
	Process(ctx context.Context, r SessionRequest) error
}

func New() Service {
	return &service{}
}

type service struct{}

func (svc *service) Process(ctx context.Context, r SessionRequest) error {
	return nil
}
