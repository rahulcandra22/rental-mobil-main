// File: controllers/auth_controller.go
package controllers

import (
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/julienschmidt/httprouter"

	"github.com/nabilulilalbab/rental-mobil/config"
	"github.com/nabilulilalbab/rental-mobil/models"
	"github.com/nabilulilalbab/rental-mobil/services"
)

type AuthController struct {
	service  services.UserService
	template *template.Template
}

func NewAuthController(service services.UserService, tmpl *template.Template) *AuthController {
	return &AuthController{service, tmpl}
}

// Menampilkan form (tidak ada perubahan)
func (c *AuthController) ShowRegisterForm(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	data := NewTemplateData(r) // Gunakan helper
	c.template.ExecuteTemplate(w, "register.html", data)
}

func (c *AuthController) ShowLoginForm(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	data := NewTemplateData(r) // Gunakan helper
	c.template.ExecuteTemplate(w, "login.html", data)
}

// Register (tidak ada perubahan)
func (c *AuthController) Register(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	r.ParseForm()
	user := models.User{
		Name:     r.PostFormValue("name"),
		Email:    r.PostFormValue("email"),
		Password: r.PostFormValue("password"),
	}
	_, err := c.service.Register(user)
	if err != nil {
		log.Printf("Gagal register: %v", err)
		http.Redirect(w, r, "/register", http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// Login sekarang akan membuat JWT dan menyimpannya di cookie
func (c *AuthController) Login(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	r.ParseForm()
	user, err := c.service.Login(r.PostFormValue("email"), r.PostFormValue("password"))
	if err != nil {
		log.Printf("Gagal login: %v", err)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// 1. Tentukan waktu kedaluwarsa token
	expTime := time.Now().Add(time.Hour * 24)

	// 2. Buat claims
	claims := &config.JWTClaims{
		UserID: user.ID,
		Name:   user.Name,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expTime),
		},
	}

	// 3. Buat token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 4. Tandatangani token dengan secret key
	tokenString, err := token.SignedString(config.JWT_KEY)
	if err != nil {
		http.Error(w, "Gagal membuat token", http.StatusInternalServerError)
		return
	}

	// 5. Set token sebagai httpOnly cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    tokenString,
		Expires:  expTime,
		Path:     "/",
		HttpOnly: true,
	})

	// 6. Redirect berdasarkan role
	if user.Role == "admin" {
		http.Redirect(w, r, "/mobil", http.StatusSeeOther)
	} else {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

// Logout akan menghapus cookie token
func (c *AuthController) Logout(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Set cookie dengan waktu kedaluwarsa di masa lalu untuk menghapusnya
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "",
		Expires:  time.Unix(0, 0),
		Path:     "/",
		HttpOnly: true,
	})
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
