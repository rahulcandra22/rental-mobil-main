package controllers

import (
	"net/http"

	"github.com/nabilulilalbab/rental-mobil/config"
)

// NewTemplateData adalah helper untuk membuat map data default untuk template
// yang berisi informasi login pengguna.
func NewTemplateData(r *http.Request) map[string]any {
	data := make(map[string]any)
	data["IsLoggedIn"] = false // Default

	// Coba ambil claims dari context (yang diisi oleh AuthMiddleware)
	if claims, ok := r.Context().Value(config.ClaimsKey).(*config.JWTClaims); ok {
		data["IsLoggedIn"] = true
		data["UserName"] = claims.Name
		data["IsAdmin"] = (claims.Role == "admin")
	}
	return data
}
