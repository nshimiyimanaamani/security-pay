package payment

// routes
const (
	DebitRoute        = "/payment/initialize"
	ProcessDebitRoute = "/payment/confirm" //used to receive the callback from the payment gateway

	CreditRoute         = "/payment/credit/initialize"
	ProcessCreditRoute  = "/payment/credit/confirm"
	PaymentReportsRoute = "/payment/reports"
	PaymentMetrics      = "/payment/metrics"
)
