package utils

import (
	"bytes"
	"strconv"

	"github.com/xuri/excelize/v2"

	"github.com/nabilulilalbab/rental-mobil/models"
)

func ExportCarsToExcel(cars []models.Car) (*bytes.Buffer, error) {
	f := excelize.NewFile()
	sheetName := "Daftart Mobil"
	f.NewSheet(sheetName)
	f.DeleteSheet("Sheet1") // Hapus sheet default

	// Set header tabel
	headers := []string{"ID", "Nama Mobil", "Kapasitas", "Transmisi", "Harga/Hari", "Status"}
	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheetName, cell, header)
	}

	// Isi data mobil
	for i, car := range cars {
		row := i + 2 // Mulai dari baris ke-2
		status := "Disewa"
		if car.IsAvailable {
			status = "Tersedia"
		}

		f.SetCellValue(sheetName, "A"+strconv.Itoa(row), car.ID)
		f.SetCellValue(sheetName, "B"+strconv.Itoa(row), car.Name)
		f.SetCellValue(sheetName, "C"+strconv.Itoa(row), car.Capacity)
		f.SetCellValue(sheetName, "D"+strconv.Itoa(row), car.Transmission)
		f.SetCellValue(sheetName, "E"+strconv.Itoa(row), car.PricePerDay)
		f.SetCellValue(sheetName, "F"+strconv.Itoa(row), status)
	}

	// Tulis ke buffer di memori
	buffer, err := f.WriteToBuffer()
	if err != nil {
		return nil, err
	}

	return buffer, nil
}

func ExportOrdersToExcel(orders []models.Order) (*bytes.Buffer, error) {
	f := excelize.NewFile()
	sheetName := "Laporan Pesanan"
	f.NewSheet(sheetName)
	f.DeleteSheet("Sheet1")

	// Set header tabel
	headers := []string{"Order ID", "Tgl Pesan", "Nama Pengguna", "Email Pengguna", "Nama Mobil", "Tgl Ambil", "Tgl Kembali", "Total Harga", "Status"}
	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheetName, cell, header)
	}

	// Isi data pesanan
	for i, order := range orders {
		row := i + 2 // Mulai dari baris ke-2

		f.SetCellValue(sheetName, "A"+strconv.Itoa(row), order.ID)
		f.SetCellValue(sheetName, "B"+strconv.Itoa(row), order.CreatedAt.Format("2 Jan 2006"))
		f.SetCellValue(sheetName, "C"+strconv.Itoa(row), order.User.Name) // Mengambil dari data preloaded
		f.SetCellValue(sheetName, "D"+strconv.Itoa(row), order.User.Email)
		f.SetCellValue(sheetName, "E"+strconv.Itoa(row), order.Car.Name) // Mengambil dari data preloaded
		f.SetCellValue(sheetName, "F"+strconv.Itoa(row), order.PickupDate.Format("2 Jan 2006"))
		f.SetCellValue(sheetName, "G"+strconv.Itoa(row), order.ReturnDate.Format("2 Jan 2006"))
		f.SetCellValue(sheetName, "H"+strconv.Itoa(row), order.TotalPrice)
		f.SetCellValue(sheetName, "I"+strconv.Itoa(row), order.Status)
	}

	buffer, err := f.WriteToBuffer()
	if err != nil {
		return nil, err
	}

	return buffer, nil
}
