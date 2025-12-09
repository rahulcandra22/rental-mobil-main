package routes

import (
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/nabilulilalbab/rental-mobil/controllers"
	"github.com/nabilulilalbab/rental-mobil/middlewares"
)

func NewRouter(
	carController *controllers.CarController,
	authController *controllers.AuthController,
	homeController *controllers.HomeController,
	orderController *controllers.OrderController,
	commentController *controllers.CommentController,
) *httprouter.Router {
	router := httprouter.New()

	// --- Grup 1: Aset Statis ---
	router.ServeFiles("/static/*filepath", http.Dir("static"))

	// --- Grup 2: Rute Publik & Otentikasi ---
	// Halaman utama (Publik dengan Auth Opsional)
	router.GET("/", middlewares.OptionalAuthMiddleware(homeController.Index))

	// Rute hanya untuk tamu (Guest-Only)
	router.GET("/register", middlewares.GuestMiddleware(authController.ShowRegisterForm))
	router.POST("/register", authController.Register)
	router.GET("/login", middlewares.GuestMiddleware(authController.ShowLoginForm))
	router.POST("/login", authController.Login)

	// --- Grup 3: Rute Pengguna (Wajib Login) ---
	router.GET("/logout", middlewares.AuthMiddleware(authController.Logout))
	router.GET("/mobil/detail/:id", middlewares.AuthMiddleware(homeController.Detail))
	router.GET("/pesan/:id", middlewares.AuthMiddleware(homeController.ShowBookingForm))
	router.POST("/pesanan/buat", middlewares.AuthMiddleware(orderController.Create))
	router.GET("/pesanan-saya", middlewares.AuthMiddleware(orderController.ListUserOrders))
	router.GET("/pesanan/bayar/:id", middlewares.AuthMiddleware(orderController.ShowPaymentPage))
	router.POST("/pesanan/konfirmasi-bayar/:id", middlewares.AuthMiddleware(orderController.ConfirmPayment))
	router.POST("/komentar/buat", middlewares.AuthMiddleware(commentController.Create))

	// --- Grup 4: Rute Admin (Wajib Login sebagai Admin) ---
	// Manajemen Mobil
	router.GET("/mobil", middlewares.AuthMiddleware(middlewares.AdminMiddleware(carController.ListCars)))
	router.GET("/mobil/tambah", middlewares.AuthMiddleware(middlewares.AdminMiddleware(carController.ShowAddCarForm)))
	router.POST("/mobil/tambah", middlewares.AuthMiddleware(middlewares.AdminMiddleware(carController.CreateCar)))
	router.GET("/mobil/edit/:id", middlewares.AuthMiddleware(middlewares.AdminMiddleware(carController.ShowEditCarForm)))
	router.POST("/mobil/update/:id", middlewares.AuthMiddleware(middlewares.AdminMiddleware(carController.UpdateCar)))
	router.POST("/mobil/delete/:id", middlewares.AuthMiddleware(middlewares.AdminMiddleware(carController.DeleteCar)))
	router.GET("/mobil/export/excel", middlewares.AuthMiddleware(middlewares.AdminMiddleware(carController.ExportCarsExcel)))

	// Manajemen Pesanan
	router.GET("/admin/pesanan", middlewares.AuthMiddleware(middlewares.AdminMiddleware(orderController.ListAllOrders)))
	router.POST("/admin/pesanan/update/:id", middlewares.AuthMiddleware(middlewares.AdminMiddleware(orderController.UpdateStatus)))
	router.GET("/admin/pesanan/export/excel", middlewares.AuthMiddleware(middlewares.AdminMiddleware(orderController.ExportOrdersExcel)))
	router.GET("/admin/pesanan/export/pdf", middlewares.AuthMiddleware(middlewares.AdminMiddleware(orderController.ExportOrdersPDF)))

	return router
}

