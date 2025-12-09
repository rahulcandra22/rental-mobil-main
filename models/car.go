package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name     string `gorm:"not null"`
	Email    string `gorm:"uniqueIndex;not null"`
	Password string `gorm:"not null"`
	Role     string `gorm:"default:'user'"` // 'user' or 'admin'

	// Relasi: Satu user bisa punya banyak pesanan
	Orders   []Order   `gorm:"foreignKey:UserID"`
	Comments []Comment `gorm:"foreignKey:UserID"`
}
type Car struct {
	gorm.Model
	Name         string `gorm:"not null"`
	Description  string // Field untuk deskripsi mobil
	Capacity     int
	Transmission string
	PricePerDay  float64 `gorm:"not null"`
	IsAvailable  bool    `gorm:"default:true"`

	// Relasi: Satu mobil punya banyak gambar
	Images []CarImage `gorm:"foreignKey:CarID"`

	// Relasi: Satu mobil bisa punya banyak pesanan
	Orders   []Order   `gorm:"foreignKey:CarID"`
	Comments []Comment `gorm:"foreignKey:CarID"`
}

type Order struct {
	gorm.Model
	UserID         uint
	CarID          uint
	PickupDate     time.Time
	ReturnDate     time.Time
	PickupLocation string
	ReturnLocation string
	TotalPrice     float64
	PaymentNote    string `gorm:"type:text"`
	Status         string `gorm:"default:'pending'"`

	// TAMBAHKAN DUA BARIS RELASI INI
	User User `gorm:"foreignKey:UserID"`
	Car  Car  `gorm:"foreignKey:CarID"`
}

type Comment struct {
	gorm.Model
	UserID  uint
	CarID   uint
	Content string `gorm:"type:text"`

	// Relasi: satu comment milik satu user
	User User `gorm:"foreignKey:UserID"`
}
