package di

import (
	"github.com/IkezawaYuki/popple/internal/infrastructure"
	"github.com/IkezawaYuki/popple/internal/repository"
)

func NewBaseRepository() repository.BaseRepository {
	return repository.NewBaseRepository(
		NewDbDriver(),
	)
}

func NewAdminRepository() repository.AdminRepository {
	return repository.NewAdminRepository(
		NewDbDriver(),
	)
}

func NewPostRepository() repository.PostRepository {
	return repository.NewPostRepository(
		NewDbDriver(),
	)
}

func NewCustomerRepository() repository.CustomerRepository {
	return repository.NewCustomerRepository(
		NewDbDriver(),
	)
}

func NewRodutRepository() repository.RodutRepository {
	return repository.NewRodutRepository(
		infrastructure.NewHttpClient(),
	)
}

func NewRedisRepository() repository.RedisRepository {
	return repository.NewRedisRepository(
		redisClient,
	)
}
