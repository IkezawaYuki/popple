package repository

import (
	"context"
	"github.com/IkezawaYuki/popple/internal/domain/model"
	"github.com/IkezawaYuki/popple/internal/infrastructure"
)

type PostRepository interface {
	Get(ctx context.Context, f *PostFilter) ([]*model.Post, error)
	GetTx(ctx context.Context, f *PostFilter, tx infrastructure.Transaction) ([]*model.Post, error)
	Save(ctx context.Context, post *model.Post) error
	SaveTx(ctx context.Context, post *model.Post, tx infrastructure.Transaction) error
}

func NewPostRepository(dbDriver infrastructure.DBDriver) PostRepository {
	return &postRepository{
		dbDriver: dbDriver,
	}
}

type postRepository struct {
	dbDriver infrastructure.DBDriver
}

func (p *postRepository) Get(ctx context.Context, f *PostFilter) ([]*model.Post, error) {
	var posts []*model.Post
	err := p.dbDriver.Find(ctx, &posts, f)
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func (p *postRepository) GetTx(ctx context.Context, f *PostFilter, tx infrastructure.Transaction) ([]*model.Post, error) {
	var posts []*model.Post
	err := p.dbDriver.FindTx(ctx, &posts, f, tx)
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func (p *postRepository) Save(ctx context.Context, post *model.Post) error {
	return p.dbDriver.Save(ctx, post)
}

func (p *postRepository) SaveTx(ctx context.Context, post *model.Post, tx infrastructure.Transaction) error {
	return p.dbDriver.SaveTx(ctx, post, tx)
}
