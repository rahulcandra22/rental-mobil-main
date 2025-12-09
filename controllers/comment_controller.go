package controllers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"

	"github.com/nabilulilalbab/rental-mobil/config"
	"github.com/nabilulilalbab/rental-mobil/models"
	"github.com/nabilulilalbab/rental-mobil/services"
)

type CommentController struct {
	service services.CommentService
}

func NewCommentController(service services.CommentService) *CommentController {
	return &CommentController{service}
}

func (c *CommentController) Create(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	claims, _ := r.Context().Value(config.ClaimsKey).(*config.JWTClaims)
	r.ParseForm()

	carID, _ := strconv.Atoi(r.PostFormValue("car_id"))

	comment := models.Comment{
		UserID:  claims.UserID,
		CarID:   uint(carID),
		Content: r.PostFormValue("content"),
	}

	_, err := c.service.CreateComment(comment)
	if err != nil {
		log.Printf("Gagal menyimpan komentar: %v", err)
	}

	// Redirect kembali ke halaman detail mobil
	http.Redirect(w, r, "/mobil/detail/"+strconv.Itoa(carID), http.StatusSeeOther)
}
