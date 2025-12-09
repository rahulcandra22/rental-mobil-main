package models

import "gorm.io/gorm"

// CarImage adalah model untuk menyimpan URL gambar dari sebuah mobil.
type CarImage struct {
	gorm.Model
	CarID uint   // Foreign Key yang menghubungkan ke tabel 'cars'
	URL   string `gorm:"not null"`
}
