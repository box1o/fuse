package credit

import (
	"encoding/json"
	"net/http"

	"fuse/internal/interfaces/server/middleware"
)

func (handler *Handler) GetBalance(
	writer http.ResponseWriter,
	request *http.Request,
) {
	ownerID := middleware.GetUserIDFromContext(
		request.Context(),
	)

	balance, err := handler.balanceService.GetBalance(
		request.Context(),
		ownerID,
	)
	if err != nil {
		http.Error(
			writer,
			"failed to load credit balance",
			http.StatusInternalServerError,
		)
		return
	}

	writer.Header().Set(
		"Content-Type",
		"application/json",
	)
	writer.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(writer).Encode(
		CreditBalanceResponse{
			Balance: balance.Value(),
		},
	); err != nil {
		return
	}
}
