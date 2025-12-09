package controllers

import (
	"encoding/base64"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"

	"github.com/nabilulilalbab/rental-mobil/config"
	"github.com/nabilulilalbab/rental-mobil/models"
	"github.com/nabilulilalbab/rental-mobil/services"
	"github.com/nabilulilalbab/rental-mobil/utils"
)

type OrderController struct {
	service  services.OrderService
	template *template.Template
}

func NewOrderController(service services.OrderService, tmpl *template.Template) *OrderController {
	return &OrderController{service, tmpl}
}

// (Ini harusnya disimpan di database atau config, tapi untuk sementara di sini dulu)
const StaticQRISPayload = "00020101021126570011ID.DANA.WWW011893600915302259148102090225914810303UMI51440014ID.CO.QRIS.WWW0215ID10200176114730303UMI5204581253033605802ID5922Warung Sayur Bu Sugeng6010Kab. Demak610559567630458C7"

func (c *OrderController) Create(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Ambil ID user yang sedang login dari context
	claims, ok := r.Context().Value(config.ClaimsKey).(*config.JWTClaims)
	if !ok {
		http.Error(w, "Gagal mengambil data user", http.StatusInternalServerError)
		return
	}

	r.ParseForm()

	carID, _ := strconv.Atoi(r.PostFormValue("car_id"))
	pickupDate, _ := time.Parse("2006-01-02", r.PostFormValue("pickup_date"))
	returnDate, _ := time.Parse("2006-01-02", r.PostFormValue("return_date"))

	order := models.Order{
		UserID:         claims.UserID,
		CarID:          uint(carID),
		PickupDate:     pickupDate,
		ReturnDate:     returnDate,
		PickupLocation: r.PostFormValue("pickup_location"),
		ReturnLocation: r.PostFormValue("return_location"),
	}

	createdOrder, err := c.service.CreateOrder(order)
	if err != nil {
		log.Printf("Gagal membuat pesanan: %v", err)
		http.Error(w, "Gagal membuat pesanan", http.StatusInternalServerError)
		return
	}

	// Redirect ke halaman utama (atau halaman 'pesanan saya' nanti)
	http.Redirect(w, r, "/pesanan/bayar/"+strconv.Itoa(int(createdOrder.ID)), http.StatusSeeOther)
}

func (c *OrderController) ListUserOrders(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	claims, _ := r.Context().Value(config.ClaimsKey).(*config.JWTClaims)

	orders, err := c.service.GetOrdersByUserID(claims.UserID)
	if err != nil {
		log.Printf("Gagal mengambil pesanan user: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := NewTemplateData(r)
	data["Orders"] = orders

	err = c.template.ExecuteTemplate(w, "my_orders.html", data)
	if err != nil {
		log.Printf("Gagal render template my_orders: %v", err)
	}
}

func (c *OrderController) ListAllOrders(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	orders, err := c.service.GetAllOrders()
	if err != nil {
		log.Printf("Gagal mengambil semua pesanan: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := NewTemplateData(r)
	data["Orders"] = orders

	err = c.template.ExecuteTemplate(w, "admin_orders.html", data)
	if err != nil {
		log.Printf("Gagal render template admin_orders: %v", err)
	}
}

// UpdateStatus memproses perubahan status pesanan dari form admin
func (c *OrderController) UpdateStatus(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	r.ParseForm()
	orderID, _ := strconv.Atoi(ps.ByName("id"))
	newStatus := r.PostFormValue("status")

	_, err := c.service.UpdateOrderStatus(uint(orderID), newStatus)
	if err != nil {
		log.Printf("Gagal update status pesanan: %v", err)
		// Tambahkan logic untuk menampilkan pesan error ke admin jika perlu
	}

	http.Redirect(w, r, "/admin/pesanan", http.StatusSeeOther)
}

func (c *OrderController) ExportOrdersExcel(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	buffer, err := c.service.ExportOrdersToExcel()
	if err != nil {
		http.Error(w, "Gagal membuat file Excel", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Disposition", "attachment; filename=laporan_pesanan.xlsx")
	w.Header().Set("Content-Length", strconv.Itoa(buffer.Len()))

	_, err = buffer.WriteTo(w)
	if err != nil {
		http.Error(w, "Gagal mengirim file", http.StatusInternalServerError)
	}
}

func (c *OrderController) ShowPaymentPage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	orderID, _ := strconv.Atoi(ps.ByName("id"))

	// Ambil data pesanan (kita perlu service baru untuk ini)
	order, err := c.service.GetOrderByID(uint(orderID))
	if err != nil {
		http.NotFound(w, r)
		return
	}

	// Generate QRIS Dinamis
	dynamicPayload, _ := utils.GenerateDynamicQRIS(StaticQRISPayload, uint(order.TotalPrice))
	// Generate Gambar QR Code
	qrCodeBytes, _ := utils.GenerateQRCodeImage(dynamicPayload)
	// Encode ke base64 untuk ditampilkan di HTML
	qrCodeBase64 := base64.StdEncoding.EncodeToString(qrCodeBytes)

	data := NewTemplateData(r)
	data["Order"] = order
	data["QRCode"] = qrCodeBase64

	c.template.ExecuteTemplate(w, "payment.html", data)
}

// Tambahkan handler baru untuk konfirmasi pembayaran (dummy)
func (c *OrderController) ConfirmPayment(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	orderID, _ := strconv.Atoi(ps.ByName("id"))
	// Di dunia nyata, ini akan dipicu oleh webhook dari payment gateway.
	// Untuk sekarang, kita update statusnya langsung.
	_, err := c.service.UpdateOrderStatus(uint(orderID), "confirmed")
	if err != nil {
		log.Printf("Gagal konfirmasi pembayaran: %v", err)
	}
	// Redirect ke halaman riwayat pesanan
	http.Redirect(w, r, "/pesanan-saya", http.StatusSeeOther)
}

func (c *OrderController) ExportOrdersPDF(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	buffer, err := c.service.ExportOrdersToPDF()
	if err != nil {
		http.Error(w, "Gagal membuat file PDF", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename=laporan_pesanan.pdf")
	w.Header().Set("Content-Length", strconv.Itoa(buffer.Len()))

	_, err = buffer.WriteTo(w)
	if err != nil {
		http.Error(w, "Gagal mengirim file", http.StatusInternalServerError)
	}
}
