package utils

import (
	"bytes"
	"fmt"

	"github.com/jung-kurt/gofpdf"

	"github.com/nabilulilalbab/rental-mobil/models"
)

func ExportOrdersToPDF(orders []models.Order) (*bytes.Buffer, error) {
	pdf := gofpdf.New("L", "mm", "A4", "") // L for Landscape, P for Portrait
	pdf.AddPage()

	// Header Dokumen
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "Laporan Semua Pesanan")
	pdf.Ln(20) // Line break

	// Header Tabel
	pdf.SetFont("Arial", "B", 10)
	pdf.SetFillColor(240, 240, 240) // Warna abu-abu untuk header
	headers := []string{"ID", "Pengguna", "Mobil", "Tgl Ambil", "Tgl Kembali", "Total Harga", "Status"}
	widths := []float64{15, 50, 60, 30, 30, 40, 30} // Lebar setiap kolom

	for i, header := range headers {
		pdf.CellFormat(widths[i], 7, header, "1", 0, "C", true, 0, "")
	}
	pdf.Ln(-1) // Pindah ke baris baru

	// Body Tabel
	pdf.SetFont("Arial", "", 10)
	pdf.SetFillColor(255, 255, 255)

	for _, order := range orders {
		pdf.CellFormat(widths[0], 7, fmt.Sprintf("%d", order.ID), "1", 0, "C", false, 0, "")
		pdf.CellFormat(widths[1], 7, order.User.Name, "1", 0, "L", false, 0, "")
		pdf.CellFormat(widths[2], 7, order.Car.Name, "1", 0, "L", false, 0, "")
		pdf.CellFormat(widths[3], 7, order.PickupDate.Format("2 Jan 2006"), "1", 0, "C", false, 0, "")
		pdf.CellFormat(widths[4], 7, order.ReturnDate.Format("2 Jan 2006"), "1", 0, "C", false, 0, "")
		pdf.CellFormat(widths[5], 7, fmt.Sprintf("Rp %.0f", order.TotalPrice), "1", 0, "R", false, 0, "")
		pdf.CellFormat(widths[6], 7, order.Status, "1", 0, "C", false, 0, "")
		pdf.Ln(-1)
	}

	// Output ke buffer
	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return nil, err
	}

	return &buf, nil
}
