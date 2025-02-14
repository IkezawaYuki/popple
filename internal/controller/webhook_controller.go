package controller

import (
	"fmt"
	"github.com/IkezawaYuki/popple/config"
	"github.com/IkezawaYuki/popple/internal/usecase"
	"github.com/labstack/echo/v4"
	"github.com/stripe/stripe-go/webhook"
	"io"
	"net/http"
	"os"
)

type WebhookController struct {
	adminUsecase usecase.AdminUsecase
}

func NewWebhookController(adminUsecase usecase.AdminUsecase) WebhookController {
	return WebhookController{
		adminUsecase: adminUsecase,
	}
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
		
	default:
		_, _ = fmt.Fprintf(os.Stderr, "Unhandled event type: %s\n", event.Type)
	}

	return c.NoContent(http.StatusOK)
}
