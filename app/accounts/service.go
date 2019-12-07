package accounts

type Service struct{
	Create(ctx conttext.Context)(Account,error)
	Retrieve(ctx context.Context id int)(Account,error)
	Delete(ctx context.Context id int)(Account,error)
}