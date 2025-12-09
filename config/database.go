package config

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/nabilulilalbab/rental-mobil/models"
)

func ConnectDatabase() *gorm.DB {
	// Kita akan menggunakan file "rental.db" sebagai database kita
	database, err := gorm.Open(sqlite.Open("rental.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Gagal terhubung ke database:", err)
	}
	log.Println("Menjalankan migrasi database...")
	err = database.AutoMigrate(
		&models.User{},
		&models.Car{},
		&models.CarImage{},
		&models.Order{},
		&models.Comment{},
	)
	if err != nil {
		log.Fatal("Gagal melakukan migrasi database:", err)
	}

	log.Println("Koneksi database dan migrasi berhasil.")
	return database
}
