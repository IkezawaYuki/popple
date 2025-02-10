package service

import (
	"context"
	"github.com/IkezawaYuki/popple/internal/domain/model"
	"github.com/IkezawaYuki/popple/internal/repository"
)

type postService struct {
	postRepo repository.PostRepository
}

type PostService interface {
	IsLinked(ctx context.Context, instagramMediaID string) (bool, error)
	Create(ctx context.Context, post *model.Post) error
	Update(ctx context.Context, post *model.Post) error
	FindByCustomerID(ctx context.Context, customerID int) ([]*model.Post, error)
}

func NewPostService(postRepo repository.PostRepository) PostService {
	return &postService{
		postRepo: postRepo,
	}
}

func (s *postService) IsLinked(ctx context.Context, instagramMediaID string) (bool, error) {
	posts, err := s.postRepo.Get(ctx, &repository.PostFilter{
		InstagramMediaID: &instagramMediaID,
	})
	if err != nil {
		return false, err
	}
	return len(posts) > 0, nil
}

func (s *postService) Create(ctx context.Context, post *model.Post) error {
	return s.postRepo.Save(ctx, post)
}

func (s *postService) Update(ctx context.Context, post *model.Post) error {
	return s.postRepo.Save(ctx, post)
}

func (s *postService) FindByCustomerID(ctx context.Context, customerID int) ([]*model.Post, error) {
	return s.postRepo.Get(ctx, &repository.PostFilter{
		CustomerID: &customerID,
	})
}
