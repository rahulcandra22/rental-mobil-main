package main

import (
	"log"
	"net/http"

	"github.com/nabilulilalbab/rental-mobil/config"
	"github.com/nabilulilalbab/rental-mobil/controllers"
	"github.com/nabilulilalbab/rental-mobil/repositories"
	"github.com/nabilulilalbab/rental-mobil/routes"
	"github.com/nabilulilalbab/rental-mobil/services"
	"github.com/nabilulilalbab/rental-mobil/utils"
)

func main() {
	// === Inisialisasi Database ===
	db := config.ConnectDatabase()

	cachedTemplates := utils.ParseTemplates()

	// === Dependency Injection dengan DB ===
	// Berikan instance 'db' ke repository
	carRepo := repositories.NewCarRepository(db)

	// Inisialisasi untuk Car
	carSvc := services.NewCarService(carRepo)
	carCtrl := controllers.NewCarController(carSvc, cachedTemplates)
	// Inisialisasi untuk User/Auth
	userRepo := repositories.NewUserRepository(db)
	userSvc := services.NewUserService(userRepo)
	authCtrl := controllers.NewAuthController(userSvc, cachedTemplates)
	// Inisialisasi untuk Halaman Utama
	homeCtrl := controllers.NewHomeController(carSvc, cachedTemplates)
	// Inisialisasi untuk Order
	orderRepo := repositories.NewOrderRepository(db)
	orderSvc := services.NewOrderService(orderRepo, carRepo, db) // Butuh carRepo & db
	orderCtrl := controllers.NewOrderController(orderSvc, cachedTemplates)
	// Inisialisasi untuk Comment
	commentRepo := repositories.NewCommentRepository(db)
	commentSvc := services.NewCommentService(commentRepo)
	commentCtrl := controllers.NewCommentController(commentSvc)
	router := routes.NewRouter(carCtrl, authCtrl, homeCtrl, orderCtrl, commentCtrl)

	port := ":8080"
	log.Printf("Server berjalan di http://localhost%s\n", port)

	err := http.ListenAndServe(port, router)
	if err != nil {
		log.Fatal(err)
	}
}
