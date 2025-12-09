package repositories

import (
	"gorm.io/gorm"

	"github.com/nabilulilalbab/rental-mobil/models"
)

type CommentRepository interface {
	Create(comment models.Comment) (models.Comment, error)
}

type commentRepositoryImpl struct {
	db *gorm.DB
}

func NewCommentRepository(db *gorm.DB) CommentRepository {
	return &commentRepositoryImpl{db}
}

func (r *commentRepositoryImpl) Create(comment models.Comment) (models.Comment, error) {
	err := r.db.Create(&comment).Error
	return comment, err
}
