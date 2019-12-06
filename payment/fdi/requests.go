package fdi

type authRequest struct {
	AppID  string `json:"appId"`
	Secret string `json:"secret"`
}

// PullRequest ...
type pullRequest struct {
	TrxRef      string  `json:"trxRef"`
	ChannelID   string  `json:"channelId"`
	AccountID   string  `json:"accountId"`
	Msisdn      string  `json:"msisdn"`
	Amount      float64 `json:"amount"`
	CallbackURL string  `json:"callback"`
}
