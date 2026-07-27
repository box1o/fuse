package payment

import (
	"io"
	"net/http"
	"strings"

	paymentService "fuse/internal/services/payment"
	"fuse/pkg/errors"
	"fuse/pkg/log"
)

const maxWebhookRequestBodySize = 1 << 20 // 1 MiB

// @Summary		Handle Stripe webhook
// @Description	Processes an incoming Stripe webhook event for payment updates.
// @Tags			payments
// @Accept		json
// @Produce		json
// @Param			Stripe-Signature	header	string	true	"Stripe signature"
// @Param			request	body	object	true	"Stripe webhook payload"
// @Success		204
// @Failure		400	{object}	errors.HTTPError
// @Failure		500	{object}	errors.HTTPError
// @Router			/payments/webhook [post]
func (h *Handler) HandleWebhook(writer http.ResponseWriter, request *http.Request) {
	payload, err := readWebhookPayload(writer, request)
	if err != nil {
		errors.WriteError(
			writer,
			errors.ErrBadRequest.WithDetail("invalid webhook payload"),
		)
		return
	}

	event, err := h.webhookParser.ParseWebhook(
		payload,
		request.Header.Get("Stripe-Signature"),
	)
	if err != nil {
		log.Warn("failed to parse payment webhook: %v", err)
		errors.WriteError(writer, errors.ToHTTP(err))
		return
	}

	if err := h.webhookService.HandleWebhook(
		request.Context(),
		event,
	); err != nil {
		log.Error(
			"failed to process payment webhook %q: %v",
			event.ID,
			err,
		)
		errors.WriteError(writer, errors.ToHTTP(err))
		return
	}

	writer.WriteHeader(http.StatusNoContent)
}

func readWebhookPayload(writer http.ResponseWriter, request *http.Request) ([]byte, error) {
	request.Body = http.MaxBytesReader(
		writer,
		request.Body,
		maxWebhookRequestBodySize,
	)
	defer request.Body.Close()

	payload, err := io.ReadAll(request.Body)
	if err != nil {
		return nil, err
	}

	if strings.TrimSpace(string(payload)) == "" {
		return nil, paymentService.ErrInvalidWebhookEvent
	}

	return payload, nil
}
