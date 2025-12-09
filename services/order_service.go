package services

import (
	"bytes"
	"errors"

	"gorm.io/gorm"

	"github.com/nabilulilalbab/rental-mobil/models"
	"github.com/nabilulilalbab/rental-mobil/repositories"
	"github.com/nabilulilalbab/rental-mobil/utils"
)

type OrderService interface {
	CreateOrder(order models.Order) (models.Order, error)
	GetOrdersByUserID(userID uint) ([]models.Order, error)
	GetAllOrders() ([]models.Order, error)
	UpdateOrderStatus(orderID uint, newStatus string) (models.Order, error)
	ExportOrdersToExcel() (*bytes.Buffer, error)
	GetOrderByID(orderID uint) (models.Order, error)
	ExportOrdersToPDF() (*bytes.Buffer, error)
}

type orderServiceImpl struct {
	orderRepo repositories.OrderRepository
	carRepo   repositories.CarRepository
	db        *gorm.DB
}

func NewOrderService(orderRepo repositories.OrderRepository, carRepo repositories.CarRepository, db *gorm.DB) OrderService {
	return &orderServiceImpl{orderRepo, carRepo, db}
}

func (s *orderServiceImpl) CreateOrder(order models.Order) (models.Order, error) {
	// Mulai transaksi
	tx := s.db.Begin()
	if tx.Error != nil {
		return models.Order{}, tx.Error
	}

	// 1. Ambil data mobil untuk validasi dan kalkulasi harga
	car, err := s.carRepo.FindByID(order.CarID)
	if err != nil {
		tx.Rollback()
		return models.Order{}, errors.New("mobil tidak ditemukan")
	}
	if !car.IsAvailable {
		tx.Rollback()
		return models.Order{}, errors.New("mobil tidak tersedia untuk dibooking")
	}

	// 2. Kalkulasi total harga di server (jangan percaya client)
	duration := order.ReturnDate.Sub(order.PickupDate).Hours() / 24
	order.TotalPrice = float64(duration) * car.PricePerDay
	order.Status = "pending" // Set status awal

	// 3. Buat record pesanan di dalam transaksi
	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		return models.Order{}, err
	}

	// 4. Update status mobil menjadi tidak tersedia di dalam transaksi
	if err := tx.Model(&car).Update("is_available", false).Error; err != nil {
		tx.Rollback()
		return models.Order{}, err
	}

	// Jika semua berhasil, commit transaksi
	return order, tx.Commit().Error
}

func (s *orderServiceImpl) GetOrdersByUserID(userID uint) ([]models.Order, error) {
	return s.orderRepo.FindByUserID(userID)
}

// Implementasi GetAllOrders
func (s *orderServiceImpl) GetAllOrders() ([]models.Order, error) {
	return s.orderRepo.FindAll()
}

// Implementasi UpdateOrderStatus
func (s *orderServiceImpl) UpdateOrderStatus(orderID uint, newStatus string) (models.Order, error) {
	// Di dunia nyata, di sini akan ada validasi status
	var order models.Order
	// Ambil order yang ada
	err := s.db.First(&order, orderID).Error
	if err != nil {
		return models.Order{}, err
	}
	// Ubah statusnya
	order.Status = newStatus
	// Simpan perubahan
	return s.orderRepo.Update(order)
}

func (s *orderServiceImpl) ExportOrdersToExcel() (*bytes.Buffer, error) {
	orders, err := s.orderRepo.FindAll() // Metode ini sudah mengambil semua data yang kita butuhkan
	if err != nil {
		return nil, err
	}
	return utils.ExportOrdersToExcel(orders)
}

func (s *orderServiceImpl) GetOrderByID(orderID uint) (models.Order, error) {
	return s.orderRepo.FindByID(orderID)
}

func (s *orderServiceImpl) ExportOrdersToPDF() (*bytes.Buffer, error) {
	orders, err := s.orderRepo.FindAll()
	if err != nil {
		return nil, err
	}
	return utils.ExportOrdersToPDF(orders)
}
