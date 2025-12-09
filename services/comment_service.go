package services

import (
	"errors"

	"github.com/nabilulilalbab/rental-mobil/models"
	"github.com/nabilulilalbab/rental-mobil/repositories"
)

type CommentService interface {
	CreateComment(comment models.Comment) (models.Comment, error)
}

type commentServiceImpl struct {
	repo repositories.CommentRepository
}

func NewCommentService(repo repositories.CommentRepository) CommentService {
	return &commentServiceImpl{repo}
}

func (s *commentServiceImpl) CreateComment(comment models.Comment) (models.Comment, error) {
	if comment.Content == "" {
		return models.Comment{}, errors.New("komentar tidak boleh kosong")
	}
	return s.repo.Create(comment)
}
