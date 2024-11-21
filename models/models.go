package models

type PaymentUpdateWsMessage struct {
	CartID  string `json:"cartID"`
	Payload string `json:"payload"` // 3 statuses: start, success, fail
}
