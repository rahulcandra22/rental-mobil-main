package services

import (
	"bytes"
	"errors"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/nabilulilalbab/rental-mobil/models"
	"github.com/nabilulilalbab/rental-mobil/repositories"
	"github.com/nabilulilalbab/rental-mobil/utils"
)

type CarService interface {
	GetAllCars() ([]models.Car, error)
	GetCarByID(id uint) (models.Car, error)
	CreateCar(car models.Car, images []*multipart.FileHeader) (models.Car, error)
	UpdateCar(car models.Car) (models.Car, error)
	DeleteCar(id uint) error
	ExportCarsToExcel() (*bytes.Buffer, error)
	SearchCarsByName(query string) ([]models.Car, error)
	// GetAllAvailableCars() ([]models.Car, error)
}

type carService struct {
	repo repositories.CarRepository
}

func NewCarService(repo repositories.CarRepository) CarService {
	return &carService{repo: repo}
}

func (s *carService) GetAllCars() ([]models.Car, error) {
	return s.repo.FindAll()
}

func (s *carService) CreateCar(car models.Car, images []*multipart.FileHeader) (models.Car, error) {
	// Logika penyimpanan file dan data dalam satu transaksi
	// Repository kita sudah di-inject dengan DB, kita bisa akses langsung dari service
	// untuk memulai transaksi. Ini adalah pola yang umum.

	// Di dunia nyata, repo akan punya method yang menerima tx,
	// tapi untuk menjaga agar tidak bertele-tele, kita lakukan di service.

	tx := s.repo.(*repositories.CarRepositoryImpl).GetDB().Begin()
	if tx.Error != nil {
		return models.Car{}, tx.Error
	}
	// 1. Simpan data mobil terlebih dahulu untuk mendapatkan ID
	if err := tx.Create(&car).Error; err != nil {
		tx.Rollback() // Batalkan transaksi jika gagal
		return models.Car{}, err
	}

	// 2. Proses dan simpan setiap gambar
	for _, imageHeader := range images {
		// Buat nama file unik

		uniqueFileName := strconv.FormatUint(uint64(car.ID), 10) + "_" + strconv.FormatInt(time.Now().UnixNano(), 10) + filepath.Ext(imageHeader.Filename)
		diskPath := filepath.Join("static", "uploads", "cars", uniqueFileName)

		// Simpan file ke disk
		dst, err := os.Create(diskPath)
		if err != nil {
			tx.Rollback()
			return models.Car{}, err
		}
		defer dst.Close()

		src, err := imageHeader.Open()
		if err != nil {
			tx.Rollback()
			return models.Car{}, err
		}
		defer src.Close()

		_, err = io.Copy(dst, src)
		if err != nil {
			tx.Rollback()
			return models.Car{}, err
		}

		// 2. Path untuk disimpan ke DB sebagai URL (selalu gunakan '/')
		webPath := "/static/uploads/cars/" + uniqueFileName

		// Simpan path URL yang benar ke database
		carImage := models.CarImage{
			CarID: car.ID,
			URL:   webPath, // GUNAKAN INI
		}
		if err := tx.Create(&carImage).Error; err != nil {
			tx.Rollback()
			return models.Car{}, err
		}
	}

	return car, tx.Commit().Error // Commit transaksi jika semua berhasil
}

func (s *carService) GetCarByID(id uint) (models.Car, error) {
	return s.repo.FindByID(id)
}

func (s *carService) UpdateCar(car models.Car) (models.Car, error) {
	if car.Name == "" {
		return models.Car{}, errors.New("nama mobil tidak boleh kosong")
	}
	return s.repo.Update(car)
}

func (s *carService) DeleteCar(id uint) error {
	return s.repo.Delete(id)
}

func (s *carService) ExportCarsToExcel() (*bytes.Buffer, error) {
	// 1. Ambil semua data mobil
	cars, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	// 2. Panggil utility untuk membuat file Excel
	return utils.ExportCarsToExcel(cars)
}

func (s *carService) SearchCarsByName(query string) ([]models.Car, error) {
	return s.repo.SearchByName(query)
}

//
// func (s *carService) GetAllAvailableCars() ([]models.Car, error) {
// 	return s.repo.FindAllAvailable()
// }
