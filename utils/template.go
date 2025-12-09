package utils

import (
	"html/template"
	"log"
)

func ParseTemplates() *template.Template {
	log.Println("Parsing templates...")

	// tmplt, err := template.ParseGlob("../templates/*.html")
	// if err != nil {
	// 	log.Fatalf("Gagal mem-parsing template layout: %v", err)
	// }

	tmplt, err := template.ParseGlob("templates/car/*.html")
	if err != nil {
		log.Fatalf("Gagal mem-parsing template produk: %v", err)
	}

	// TAMBAHKAN INI: Parsing folder partials
	_, err = tmplt.ParseGlob("templates/partials/*.html")
	if err != nil {
		log.Fatalf("Gagal mem-parsing template partials: %v", err)
	}
	_, err = tmplt.ParseGlob("templates/auth/*.html")
	if err != nil {
		log.Fatalf("Gagal mem-parsing template partials: %v", err)
	}
	_, err = tmplt.ParseGlob("templates/*.html")
	if err != nil {
		log.Fatalf("Gagal mem-parsing template partials: %v", err)
	}

	log.Println("Parsing templates selesai.")
	return tmplt
}
