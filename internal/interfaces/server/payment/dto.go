package payment

import "github.com/google/uuid"

type CreateCheckoutRequest struct {
	CreditPackID uuid.UUID `json:"credit_pack_id"`
	SuccessURL   string    `json:"success_url"`
	CancelURL    string    `json:"cancel_url"`
}

type CreateCheckoutResponse struct {
	PaymentID   uuid.UUID `json:"payment_id"`
	SessionID   string    `json:"session_id"`
	CheckoutURL string    `json:"checkout_url"`
}
