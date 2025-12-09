package controllers

import (
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"gorm.io/gorm"

	"github.com/nabilulilalbab/rental-mobil/models"
	"github.com/nabilulilalbab/rental-mobil/services"
)

type CarController struct {
	service  services.CarService
	template *template.Template
}

func NewCarController(service services.CarService, tmpl *template.Template) *CarController {
	return &CarController{service: service, template: tmpl}
}

func (c *CarController) ListCars(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	cars, err := c.service.GetAllCars()
	if err != nil {
		http.Error(w, "Gagal mengambil data mobil", http.StatusInternalServerError)
		return
	}
	data := NewTemplateData(r)
	data["DataMobil"] = cars

	if err = c.template.ExecuteTemplate(w, "index.html", data); err != nil {
		log.Printf("Error executing template: %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (c *CarController) ShowAddCarForm(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	err := c.template.ExecuteTemplate(w, "add.html", nil)
	if err != nil {
		log.Printf("Error executing template: %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (c *CarController) CreateCar(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "Ukuran file terlalu besar", http.StatusBadRequest)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Gagal mem-parsing form", http.StatusBadRequest)
		return
	}
	price, err := strconv.ParseFloat(r.PostFormValue("price"), 64)
	if err != nil {
		http.Error(w, "Harga tidak valid", http.StatusBadRequest)
		return
	}
	capacity, err := strconv.Atoi(r.PostFormValue("capacity"))
	if err != nil {
		http.Error(w, "Kapasitas tidak valid", http.StatusBadRequest)
		return
	}
	newCar := models.Car{
		Name:         r.PostFormValue("name"),
		Description:  r.PostFormValue("description"),
		Transmission: r.PostFormValue("transmission"),
		Capacity:     capacity,
		PricePerDay:  price,
	}
	images := r.MultipartForm.File["images"]
	_, err = c.service.CreateCar(newCar, images)
	if err != nil {
		log.Printf("Gagal menyimpan mobil: %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/mobil", http.StatusSeeOther)
}

func (c *CarController) ShowEditCarForm(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		http.Error(w, "ID tidak valid", http.StatusBadRequest)
		return
	}
	car, err := c.service.GetCarByID(uint(id))
	if err != nil {
		http.Error(w, "Mobil tidak ditemukan", http.StatusNotFound)
		return
	}
	data := NewTemplateData(r)
	data["Mobil"] = car

	err = c.template.ExecuteTemplate(w, "edit.html", data)
	if err != nil {
		http.Error(w, "gagal parsing", http.StatusNotFound)
	}
}

// File: controllers/car_controller.go

func (c *CarController) UpdateCar(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Wajib panggil ParseForm untuk method POST
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Gagal mem-parsing form", http.StatusBadRequest)
		return
	}

	// Ambil ID dari URL dan lakukan konversi dengan error handling
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		http.Error(w, "ID tidak valid", http.StatusBadRequest)
		return
	}

	// Ambil semua data dari form dan konversi dengan error handling
	price, err := strconv.ParseFloat(r.PostFormValue("price"), 64)
	if err != nil {
		http.Error(w, "Harga tidak valid", http.StatusBadRequest)
		return
	}
	capacity, err := strconv.Atoi(r.PostFormValue("capacity"))
	if err != nil {
		http.Error(w, "Kapasitas tidak valid", http.StatusBadRequest)
		return
	}

	// 2. Ambil status 'IsAvailable' dari form dan konversi ke boolean
	isAvailable := r.PostFormValue("is_available") == "true"

	updatedCar := models.Car{
		Model:        gorm.Model{ID: uint(id)}, // 1. Sertakan ID mobil agar GORM tahu ini adalah UPDATE
		Name:         r.PostFormValue("name"),
		Description:  r.PostFormValue("description"), // Jangan lupa description
		Transmission: r.PostFormValue("transmission"),
		Capacity:     capacity,
		PricePerDay:  price,
		IsAvailable:  isAvailable, // Gunakan nilai baru dari form
	}

	_, err = c.service.UpdateCar(updatedCar)
	if err != nil {
		log.Printf("Gagal update mobil: %v", err)
		http.Error(w, "Gagal memperbarui data mobil", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/mobil", http.StatusSeeOther)
}

func (c *CarController) DeleteCar(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		http.Error(w, "id not valid", http.StatusNotFound)
	}

	err = c.service.DeleteCar(uint(id))
	if err != nil {
		http.Error(w, "gagal deleted cars", http.StatusNotFound)
	}

	http.Redirect(w, r, "/mobil", http.StatusSeeOther)
}

func (c *CarController) ExportCarsExcel(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	buffer, err := c.service.ExportCarsToExcel()
	if err != nil {
		http.Error(w, "Gagal membuat file Excel", http.StatusInternalServerError)
		return
	}

	// Set HTTP Headers untuk file download
	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Disposition", "attachment; filename=daftar_mobil.xlsx")
	w.Header().Set("Content-Length", strconv.Itoa(buffer.Len()))

	// Kirim buffer sebagai response
	_, err = buffer.WriteTo(w)
	if err != nil {
		http.Error(w, "Gagal mengirim file", http.StatusInternalServerError)
	}
}

func MustFloat(s string) float64 {
	i, _ := strconv.ParseFloat(s, 64)
	return i
}
