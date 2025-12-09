package utils

import (
	"fmt"
	"strings"

	"github.com/sigurn/crc16"
	"github.com/skip2/go-qrcode"
)

// GenerateDynamicQRIS mengonversi payload QRIS statis menjadi dinamis dengan nominal tertentu.
// Ini adalah hasil rewrite dari script PHP Anda.
func GenerateDynamicQRIS(staticPayload string, amount uint) (string, error) {
	// 1. Hapus 4 karakter CRC dari payload statis
	payload := staticPayload[:len(staticPayload)-4]

	// 2. Ganti tipe QR menjadi dinamis (dari 11 ke 12)
	payload = strings.Replace(payload, "010211", "010212", 1)

	// 3. Tambahkan nominal transaksi (Tag 54)
	amountStr := fmt.Sprintf("%d", amount)
	transactionAmount := fmt.Sprintf("54%02d%s", len(amountStr), amountStr)

	// 4. Sisipkan nominal sebelum tag negara (Tag 58)
	parts := strings.Split(payload, "5802ID")
	payload = parts[0] + transactionAmount + "5802ID" + parts[1]

	// 5. Hitung ulang CRC16
	table := crc16.MakeTable(crc16.CRC16_CCITT_FALSE)
	crc := crc16.Checksum([]byte(payload), table)
	crcStr := fmt.Sprintf("%04X", crc) // Format ke Hex 4 digit

	return payload + crcStr, nil
}

// GenerateQRCodeImage membuat gambar QR code dari payload dan mengembalikannya sebagai byte PNG.
func GenerateQRCodeImage(payload string) ([]byte, error) {
	// Menghasilkan QR code dalam bentuk byte slice (PNG)
	png, err := qrcode.Encode(payload, qrcode.Medium, 256)
	if err != nil {
		return nil, err
	}
	return png, nil
}
