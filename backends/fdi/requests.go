package fdi

type authorization struct {
	ID     string `json:"appId"`
	Secret string `json:"secret"`
}

// Request ...
type Request struct {
	TrxRef      string  `json:"trxRef"`
	ChannelID   string  `json:"channelId"`
	AccountID   string  `json:"accountId"`
	Msisdn      string  `json:"msisdn"`
	Amount      float64 `json:"amount"`
	CallbackURL string  `json:"callback"`
}
