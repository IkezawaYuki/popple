package usecase

import (
	"context"
	"github.com/IkezawaYuki/popple/internal/service"
	"github.com/IkezawaYuki/popple/internal/usecase/dto/res"
	"sync"
)

type batchUsecase struct {
	customerService service.CustomerService
	customerUsecase CustomerUsecase
	slack           service.SlackService
}

type BatchUsecase interface {
	Execute(ctx context.Context) (*res.Message, error)
}

func NewBatchUsecase(customerService service.CustomerService, customerUsecase CustomerUsecase, slackService service.SlackService) BatchUsecase {
	return &batchUsecase{
		customerService: customerService,
		customerUsecase: customerUsecase,
		slack:           slackService,
	}
}

func (b *batchUsecase) Execute(ctx context.Context) (*res.Message, error) {
	customers, err := b.customerService.FindAuthCustomers(ctx)
	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup
	sem := make(chan struct{}, 20) // 最大20件の並列処理を制限するためのセマフォ

	for _, customer := range customers {
		wg.Add(1)

		// セマフォを取得
		sem <- struct{}{}

		// Start a goroutine for each customer
		go func(customerID int) {
			defer wg.Done()
			defer func() { <-sem }() // 処理が完了したらセマフォを解放

			if _, err := b.customerUsecase.FetchAndPost(ctx, customerID); err != nil {
				_ = b.slack.SendAlert(ctx, err.Error())
				return
			}

		}(customer.ID)
	}
	
	wg.Wait()

	return &res.Message{Message: "ok"}, nil
}
