package repository

import (
	"github.com/IkezawaYuki/popple/internal/infrastructure"
)

type BaseRepository interface {
	Begin() infrastructure.Transaction
}

type baseRepository struct {
	dbDriver infrastructure.DBDriver
}

func NewBaseRepository(dbDriver infrastructure.DBDriver) BaseRepository {
	return &baseRepository{
		dbDriver: dbDriver,
	}
}

func (b *baseRepository) Begin() infrastructure.Transaction {
	return b.dbDriver.Begin()
}
