package repositories

import (
	"gorm.io/gorm"

	"github.com/nabilulilalbab/rental-mobil/models"
)

type CarRepository interface {
	FindAll() ([]models.Car, error)
	FindByID(id uint) (models.Car, error)
	Save(car models.Car) (models.Car, error)
	Update(car models.Car) (models.Car, error)
	Delete(id uint) error
	GetDB() *gorm.DB
	// SearchAvailable(query string) ([]models.Car, error)
	// FindAllAvailable() ([]models.Car, error)
	SearchByName(query string) ([]models.Car, error)
}

// UBAH INI: dari 'carRepository' menjadi 'CarRepository'
type CarRepositoryImpl struct {
	db *gorm.DB
}

// UBAH INI: return &CarRepositoryImpl
func NewCarRepository(db *gorm.DB) CarRepository {
	return &CarRepositoryImpl{db: db}
}

// UBAH SEMUA RECEIVER: dari (r *carRepository) menjadi (r *CarRepositoryImpl)
func (r *CarRepositoryImpl) FindAll() ([]models.Car, error) {
	var cars []models.Car
	err := r.db.Preload("Images").Find(&cars).Error
	return cars, err
}

func (r *CarRepositoryImpl) FindByID(id uint) (models.Car, error) {
	var car models.Car
	err := r.db.Preload("Images").Preload("Comments.User").First(&car, id).Error
	return car, err
}

func (r *CarRepositoryImpl) Save(car models.Car) (models.Car, error) {
	err := r.db.Create(&car).Error
	return car, err
}

func (r *CarRepositoryImpl) Update(car models.Car) (models.Car, error) {
	err := r.db.Save(&car).Error
	return car, err
}

func (r *CarRepositoryImpl) Delete(id uint) error {
	err := r.db.Delete(&models.Car{}, id).Error
	return err
}

func (r *CarRepositoryImpl) GetDB() *gorm.DB {
	return r.db
}

// Tambahkan implementasi metode baru di CarRepositoryImpl
func (r *CarRepositoryImpl) SearchByName(query string) ([]models.Car, error) {
	var cars []models.Car
	searchQuery := "%" + query + "%"
	// Hapus filter "is_available"
	err := r.db.Preload("Images").
		Where("name LIKE ?", searchQuery).
		Find(&cars).Error
	return cars, err
}

//
// func (r *CarRepositoryImpl) FindAllAvailable() ([]models.Car, error) {
// 	var cars []models.Car
// 	err := r.db.Preload("Images").
// 		Where("is_available = ?", true).
// 		Find(&cars).Error
// 	return cars, err
// }
