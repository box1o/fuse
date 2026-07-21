package payment

import "github.com/google/uuid"

// CreateCheckoutRequest contains the details required to create a payment checkout session.
type CreateCheckoutRequest struct {
	CreditPackID uuid.UUID `json:"credit_pack_id"`
	SuccessURL   string    `json:"success_url"`
	CancelURL    string    `json:"cancel_url"`
}

// CreateCheckoutResponse contains the checkout session information returned after a successful payment initialization.
type CreateCheckoutResponse struct {
	PaymentID   uuid.UUID `json:"payment_id"`
	SessionID   string    `json:"session_id"`
	CheckoutURL string    `json:"checkout_url"`
}
