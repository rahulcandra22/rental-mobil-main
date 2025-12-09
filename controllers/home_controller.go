package controllers

import (
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"

	"github.com/nabilulilalbab/rental-mobil/models"
	"github.com/nabilulilalbab/rental-mobil/services"
)

type HomeController struct {
	carService services.CarService
	template   *template.Template
}

func NewHomeController(carService services.CarService, tmpl *template.Template) *HomeController {
	return &HomeController{carService, tmpl}
}

// Index akan menangani halaman utama dan hasil pencarian
func (c *HomeController) Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Ambil query pencarian dari URL, contoh: /?q=avanza
	searchQuery := r.URL.Query().Get("q")

	var cars []models.Car
	var err error

	if searchQuery != "" {
		// Jika ada query, lakukan pencarian
		cars, err = c.carService.SearchCarsByName(searchQuery)
	} else {
		// Jika tidak, tampilkan semua mobil tersedia
		cars, err = c.carService.GetAllCars()
	}

	if err != nil {
		log.Printf("Gagal mengambil data mobil: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := NewTemplateData(r) // Gunakan helper dari base_controller
	data["DataMobil"] = cars
	data["SearchQuery"] = searchQuery
	err = c.template.ExecuteTemplate(w, "home.html", data)
	if err != nil {
		log.Printf("Gagal render template: %v", err)
	}
}

func (c *HomeController) Detail(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		http.Error(w, "ID mobil tidak valid", http.StatusBadRequest)
		return
	}
	car, err := c.carService.GetCarByID(uint(id))
	if err != nil {
		// Jika mobil tidak ditemukan, tampilkan halaman 404
		http.NotFound(w, r)
		return
	}
	data := NewTemplateData(r)
	data["Car"] = car
	err = c.template.ExecuteTemplate(w, "detail.html", data)
	if err != nil {
		log.Printf("Gagal render template detail: %v", err)
	}
}

func (c *HomeController) ShowBookingForm(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		http.Error(w, "ID mobil tidak valid", http.StatusBadRequest)
		return
	}

	car, err := c.carService.GetCarByID(uint(id))
	if err != nil {
		http.NotFound(w, r)
		return
	}

	// Jangan izinkan booking jika mobil tidak tersedia
	if !car.IsAvailable {
		// Di dunia nyata, kita bisa menampilkan pesan error yang lebih baik
		http.Redirect(w, r, "/mobil/detail/"+strconv.Itoa(id), http.StatusSeeOther)
		return
	}

	data := NewTemplateData(r)
	data["Car"] = car

	err = c.template.ExecuteTemplate(w, "booking.html", data)
	if err != nil {
		log.Printf("Gagal render template booking: %v", err)
	}
}
