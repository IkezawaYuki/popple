package service

import (
	"context"
	"github.com/IkezawaYuki/popple/config"
	"github.com/IkezawaYuki/popple/internal/infrastructure"
	"github.com/IkezawaYuki/popple/internal/usecase/dto/exreq"
)

type SlackService interface {
	SendAlert(ctx context.Context, msg string) error
	SendNotification(ctx context.Context, msg string) error
}

type slackService struct {
	webhookURL string
	httpClient infrastructure.HttpClient
}

func NewSlackService(httpClient infrastructure.HttpClient) SlackService {
	return &slackService{
		webhookURL: config.Env.SlackWebhookURL,
		httpClient: httpClient,
	}
}

func (s *slackService) SendAlert(ctx context.Context, msg string) error {
	payload := exreq.SlackPayload{
		IconEmoji: ":wink",
		Text:      msg,
		Username:  "popple",
	}
	_, err := s.httpClient.PostRequest(ctx, s.webhookURL, payload, "")
	if err != nil {
		return err
	}
	return nil
}

func (s *slackService) SendNotification(ctx context.Context, msg string) error {
	payload := exreq.SlackPayload{
		IconEmoji: ":wink",
		Text:      msg,
		Username:  "popple",
	}
	_, err := s.httpClient.PostRequest(ctx, s.webhookURL, payload, "")
	if err != nil {
		return err
	}
	return nil
}
