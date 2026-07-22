package credit

import domainCredit "fuse/internal/domain/credit"

type CreditPackResponse struct {
	ID          string `json:"id"`
	Code        string `json:"code"`
	Name        string `json:"name"`
	Credits     int64  `json:"credits"`
	PriceAmount int64  `json:"price_amount"`
	Currency    string `json:"currency"`
}

type CreditBalanceResponse struct {
	Balance int64 `json:"balance"`
}

func newCreditPackResponse(pack *domainCredit.Pack) CreditPackResponse {
	return CreditPackResponse{
		ID:          pack.ID.String(),
		Code:        pack.Code,
		Name:        pack.Name,
		Credits:     pack.Credits.Value(),
		PriceAmount: pack.PriceAmount,
		Currency:    pack.Currency,
	}
}
