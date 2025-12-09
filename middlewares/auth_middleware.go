package middlewares

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/julienschmidt/httprouter"

	"github.com/nabilulilalbab/rental-mobil/config"
)

// Definisikan tipe custom untuk kunci konteks agar tidak bentrok
type contextKey string

const (
	claimsKey contextKey = "claims"
)

// AuthMiddleware memeriksa token JWT yang valid
func AuthMiddleware(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		// 1. Ambil cookie token
		cookie, err := r.Cookie("token")
		if err != nil {
			// Jika tidak ada cookie, arahkan ke login
			log.Println("Middleware: Cookie token tidak ditemukan")
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// 2. Parse dan validasi token
		tokenString := cookie.Value
		claims := &config.JWTClaims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return config.JWT_KEY, nil
		})

		if err != nil || !token.Valid {
			// Jika token tidak valid (error atau kedaluwarsa), hapus cookie dan arahkan ke login
			log.Printf("Middleware: Token tidak valid: %v", err)
			http.SetCookie(w, &http.Cookie{
				Name: "token", Value: "", Expires: time.Unix(0, 0), Path: "/", HttpOnly: true,
			})
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// 3. Simpan claims ke dalam context request
		ctx := context.WithValue(r.Context(), config.ClaimsKey, claims)

		// Lanjutkan ke handler berikutnya dengan request yang sudah dimodifikasi
		next(w, r.WithContext(ctx), ps)
	}
}

// AdminMiddleware memeriksa apakah role user adalah 'admin'
func AdminMiddleware(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		// Ambil claims dari context (yang sudah diisi oleh AuthMiddleware)
		claims, ok := r.Context().Value(config.ClaimsKey).(*config.JWTClaims)
		if !ok || claims.Role != "admin" {
			log.Printf("Akses ditolak untuk user: %s, role: %s", claims.Name, claims.Role)
			http.Error(w, "Akses ditolak: Anda bukan admin.", http.StatusForbidden)
			return
		}

		// Jika admin, lanjutkan ke handler berikutnya
		next(w, r, ps)
	}
}

// GuestMiddleware adalah kebalikan dari AuthMiddleware.
// Jika user sudah login (punya token valid), arahkan ke halaman utama.
// Jika tidak, biarkan user melihat halaman (misal: halaman login/register).
func GuestMiddleware(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		cookie, err := r.Cookie("token")
		if err == nil { // Jika cookie ditemukan
			tokenString := cookie.Value
			_, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				return config.JWT_KEY, nil
			})

			if err == nil { // Dan tokennya valid
				// Arahkan ke halaman utama karena sudah login
				http.Redirect(w, r, "/", http.StatusSeeOther)
				return
			}
		}

		// Jika tidak ada cookie atau token tidak valid, lanjutkan ke halaman login/register
		next(w, r, ps)
	}
}

func OptionalAuthMiddleware(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		cookie, err := r.Cookie("token")
		// Jika tidak ada cookie atau ada error, lanjutkan saja sebagai tamu
		if err != nil {
			next(w, r, ps)
			return
		}

		tokenString := cookie.Value
		claims := &config.JWTClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return config.JWT_KEY, nil
		})

		// Jika token valid, tambahkan claims ke context
		if err == nil && token.Valid {
			ctx := context.WithValue(r.Context(), config.ClaimsKey, claims)
			next(w, r.WithContext(ctx), ps)
			return
		}

		// Jika token tidak valid, lanjutkan saja sebagai tamu
		next(w, r, ps)
	}
}
