package usecase

import (
	"context"
	"github.com/IkezawaYuki/popple/internal/service"
	"sync"
)

type BatchUsecase struct {
	customerService service.CustomerService
	customerUsecase CustomerUsecase
	slack           service.SlackService
}

func NewBatchUsecase(customerUsecase CustomerUsecase, slackService service.SlackService) *BatchUsecase {
	return &BatchUsecase{
		customerUsecase: customerUsecase,
		slack:           slackService,
	}
}

func (b *BatchUsecase) Execute(ctx context.Context) error {
	customers, err := b.customerService.FindAuthCustomers(ctx)
	if err != nil {
		return err
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

			// Fetch Instagram Media
			if err := b.customerUsecase.FetchAndPost(ctx, customerID); err != nil {
				return
			}

		}(customer.ID)
	}

	// Wait for all goroutines to finish
	wg.Wait()

	return nil
}
