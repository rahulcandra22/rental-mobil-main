package repositories

import (
	"gorm.io/gorm"

	"github.com/nabilulilalbab/rental-mobil/models"
)

type OrderRepository interface {
	Create(order models.Order) (models.Order, error)
	FindByUserID(userID uint) ([]models.Order, error)
	FindAll() ([]models.Order, error)
	Update(order models.Order) (models.Order, error)
	FindByID(orderID uint) (models.Order, error)
}

type orderRepositoryImpl struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepositoryImpl{db: db}
}

func (r *orderRepositoryImpl) Create(order models.Order) (models.Order, error) {
	err := r.db.Create(&order).Error
	return order, err
}

func (r *orderRepositoryImpl) FindByUserID(userID uint) ([]models.Order, error) {
	var orders []models.Order
	// Preload("Car") akan mengambil data mobil terkait
	// Preload("Car.Images") akan mengambil gambar dari mobil tersebut
	// Order("created_at DESC") akan mengurutkan dari yang terbaru
	err := r.db.Preload("Car.Images").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&orders).Error
	return orders, err
}

func (r *orderRepositoryImpl) FindAll() ([]models.Order, error) {
	var orders []models.Order
	// Preload User dan Car untuk menampilkan info relevan
	err := r.db.Preload("User").Preload("Car").Order("created_at DESC").Find(&orders).Error
	return orders, err
}

// Implementasi Update
func (r *orderRepositoryImpl) Update(order models.Order) (models.Order, error) {
	// Save akan update record jika Primary Key ada
	err := r.db.Save(&order).Error
	return order, err
}

func (r *orderRepositoryImpl) FindByID(orderID uint) (models.Order, error) {
	var order models.Order
	// Preload Car agar data mobil bisa ditampilkan di halaman pembayaran
	err := r.db.Preload("Car").First(&order, orderID).Error
	return order, err
}
