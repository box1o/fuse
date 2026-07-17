package payments

import "time"

type CreditGrant struct {
	ID                string    `json:"id"`
	UserID            string    `json:"user_id"`
	CreditPackID      string    `json:"credit_pack_id"`
	Credits           int64     `json:"credits"`
	Source            string    `json:"stripe_checkout"`
	ProviderReference string    `json:"provider_reference"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type CreditBalance struct {
	GrantedCredits   int64 `json:"granted_credits"`
	UsedCredits      int64 `json:"used_credits"`
	AvailableCredits int64 `json:"available_credits"`
}
