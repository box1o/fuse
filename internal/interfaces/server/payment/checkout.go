package payment

import (
	"encoding/json"
	"io"
	"net/http"

	"fuse/internal/interfaces/server/middleware"
	paymentService "fuse/internal/services/payment"
	"fuse/pkg/errors"
	"fuse/pkg/log"

	"github.com/google/uuid"
)

const maxCheckoutRequestBodySize = 64 * 1024

// @Summary		Create a payment checkout session
// @Description	Creates a Stripe checkout session for the selected credit pack.
// @Tags			payments
// @Accept		json
// @Produce		json
// @Param			request	body	CreateCheckoutRequest	true	"Checkout details"
// @Success		201	{object}	CreateCheckoutResponse
// @Failure		400	{object}	errors.HTTPError
// @Failure		401	{object}	errors.HTTPError
// @Failure		500	{object}	errors.HTTPError
// @Router			/payments/checkout [post]
func (h *Handler) CreateCheckout(
	writer http.ResponseWriter,
	request *http.Request,
) {
	var body CreateCheckoutRequest

	if err := decodeCreateCheckoutRequest(writer, request, &body); err != nil {
		log.Warn("failed to decode checkout request: %v", err)

		errors.WriteError(
			writer,
			errors.ErrBadRequest.WithDetail("invalid request payload"),
		)
		return
	}

	ownerID := middleware.GetUserIDFromContext(request.Context())
	if ownerID == uuid.Nil {
		errors.WriteError(
			writer,
			errors.ErrUnauthorized.WithDetail("user not authenticated"),
		)
		return
	}

	output, err := h.checkoutService.CreateCheckout(
		request.Context(),
		paymentService.CreateCheckoutInput{
			OwnerID:      ownerID,
			CreditPackID: body.CreditPackID,
			SuccessURL:   body.SuccessURL,
			CancelURL:    body.CancelURL,
		},
	)
	if err != nil {
		log.Warn("failed to create payment checkout: %v", err)
		errors.WriteError(writer, errors.ToHTTP(err))
		return
	}

	response := CreateCheckoutResponse{
		PaymentID:   output.PaymentID,
		SessionID:   output.SessionID,
		CheckoutURL: output.CheckoutURL,
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(writer).Encode(response); err != nil {
		log.Error("failed to encode checkout response: %v", err)
	}
}

func decodeCreateCheckoutRequest(
	writer http.ResponseWriter,
	request *http.Request,
	destination *CreateCheckoutRequest,
) error {
	request.Body = http.MaxBytesReader(
		writer,
		request.Body,
		maxCheckoutRequestBodySize,
	)
	defer request.Body.Close()

	decoder := json.NewDecoder(request.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(destination); err != nil {
		return err
	}

	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		return errors.ErrBadRequest.WithDetail(
			"request body must contain one JSON object",
		)
	}

	return nil
}
