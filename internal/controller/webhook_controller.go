package controller

import (
	"fmt"
	"github.com/IkezawaYuki/popple/config"
	"github.com/labstack/echo/v4"
	"github.com/stripe/stripe-go/webhook"
	"io"
	"net/http"
	"os"
)

type WebhookController struct {
}

func NewWebhookController() WebhookController {
	return WebhookController{}
}

func (w *WebhookController) StripeWebhook(c echo.Context) error {
	payload, err := io.ReadAll(c.Request().Body)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error reading request body: %v\n", err)
		return c.NoContent(http.StatusServiceUnavailable)
	}

	event, err := webhook.ConstructEvent(payload, c.Request().Header.Get("Stripe-Signature"), config.Env.StripeEndpointSecret)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error verifying webhook signature: %v\n", err)
		return c.NoContent(http.StatusBadRequest)
	}

	// Unmarshal the event data into an appropriate struct depending on its Type
	switch event.Type {
	case "payment_intent.succeeded":
		// Then define and call a function to handle the event payment_intent.succeeded
		// ... handle other event types
	default:
		_, _ = fmt.Fprintf(os.Stderr, "Unhandled event type: %s\n", event.Type)
	}

	return c.NoContent(http.StatusOK)
}
