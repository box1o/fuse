package credit

import (
	"encoding/json"
	"net/http"
)

func (handler *Handler) ListActivePacks(
	writer http.ResponseWriter,
	request *http.Request,
) {
	packs, err := handler.packService.ListActivePacks(
		request.Context(),
	)
	if err != nil {
		http.Error(
			writer,
			"failed to load credit packs",
			http.StatusInternalServerError,
		)
		return
	}

	response := make(
		[]CreditPackResponse,
		0,
		len(packs),
	)

	for _, pack := range packs {
		response = append(
			response,
			newCreditPackResponse(pack),
		)
	}

	writer.Header().Set(
		"Content-Type",
		"application/json",
	)

	writer.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(writer).Encode(response); err != nil {
		return
	}
}
