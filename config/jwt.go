package config

import "github.com/golang-jwt/jwt/v5"

// Kunci ini HARUS SANGAT RAHASIA dan idealnya diambil dari environment variable
var JWT_KEY = []byte("kunci-rahasia-yang-sangat-panjang-dan-sulit-ditebak-456")

// Claims adalah data yang kita simpan di dalam token
type JWTClaims struct {
	UserID uint   `json:"user_id"`
	Name   string `json:"name"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

type ContextKey string

const ClaimsKey ContextKey = "claims"
