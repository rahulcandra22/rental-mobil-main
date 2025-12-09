# Aplikasi Rental Mobil

A comprehensive web application for car rental management, built with Go and leveraging a clean, modern UI with Tailwind CSS.

## Fitur

*   Autentikasi pengguna (Registrasi, Login, Logout)
*   Manajemen mobil (Tambah, Edit, Hapus mobil - Khusus Admin)
*   Penjelajahan dan pencarian mobil
*   Pemesanan mobil dan manajemen pesanan
*   Integrasi pembayaran QRIS
*   Riwayat pesanan untuk pengguna
*   Panel admin untuk manajemen pesanan
*   Ekspor data pesanan ke Excel dan PDF
*   Desain responsif dengan pengalih mode gelap

## Teknologi yang Digunakan

*   **Backend:** Go
*   **Database:** SQLite (via GORM)
*   **Kerangka Kerja Web:** `julienschmidt/httprouter`
*   **Autentikasi:** JWT (`golang-jwt/jwt/v5`)
*   **Styling:** Tailwind CSS
*   **Pembuatan Kode QR:** `skip2/go-qrcode`
*   **Pembuatan PDF:** `jung-kurt/gofpdf`
*   **Ekspor Excel:** `xuri/excelize/v2`
*   **Pustaka Go Lainnya:** `golang.org/x/crypto` (untuk hashing kata sandi), `mattn/go-sqlite3`, `jinzhu/inflection`, `jinzhu/now`, `sigurn/crc16`, `tiendc/go-deepcopy`, `xuri/efp`, `xuri/nfp`, `golang.org/x/net`, `golang.org/x/text`.

## Struktur Proyek

Proyek ini mengikuti arsitektur modular dan berlapis untuk memastikan pemeliharaan, skalabilitas, dan pemisahan kekhawatiran yang jelas.

```
rental-mobil/
├── cmd/
│   └── main.go
├── config/
│   ├── database.go
│   └── jwt.go
├── controllers/
│   ├── auth_controller.go
│   ├── base_controller.go
│   ├── car_controller.go
│   ├── comment_controller.go
│   ├── home_controller.go
│   └── order_controller.go
├── middlewares/
│   └── auth_middleware.go
├── models/
│   ├── car_image.go
│   └── car.go
├── repositories/
│   ├── car_repository.go
│   ├── comment_repository.go
│   ├── order_repository.go
│   └── user_repository.go
├── routes/
│   └── routes.go
├── services/
│   ├── car_service.go
│   ├── comment_service.go
│   ├── order_service.go
│   └── user_service.go
├── static/
│   ├── css/
│   │   ├── input.css
│   │   └── output.css
│   ├── js/
│   │   └── theme.js
│   ├── uploads/
│   │   └── cars/
│   │       └── ... (Gambar mobil yang diunggah)
│   ├── package-lock.json
│   ├── package.json
│   └── tailwind.config.js
├── templates/
│   ├── auth/
│   │   ├── login.html
│   │   └── register.html
│   ├── car/
│   │   ├── add.html
│   │   ├── edit.html
│   │   └── index.html
│   ├── partials/
│   │   ├── footer.html
│   │   ├── navbar.html
│   │   └── test.html
│   ├── admin_orders.html
│   ├── booking.html
│   ├── detail.html
│   ├── home.html
│   └── payment.html
├── tests/
├── utils/
│   ├── excel_exporter.go
│   ├── pdf_exporter.go
│   ├── qris.go
│   └── template.go
├── go.mod
├── go.sum
├── README.md
└── rental.db
```

### Rincian Direktori

#### `cmd/`
Direktori ini berisi titik masuk utama aplikasi.
*   `main.go`: File utama yang menginisialisasi aplikasi, menyiapkan koneksi database, mendaftarkan rute, dan memulai server HTTP.

#### `config/`
Direktori ini menyimpan file-file terkait konfigurasi.
*   `database.go`: Mengelola koneksi dan migrasi database. Biasanya menginisialisasi GORM dan terhubung ke database SQLite (`rental.db`).
*   `jwt.go`: Menangani konfigurasi JSON Web Token (JWT), termasuk manajemen kunci rahasia dan logika pembuatan/validasi token.

#### `controllers/`
Pengontrol bertanggung jawab untuk menangani permintaan HTTP yang masuk, memprosesnya, dan mengembalikan respons yang sesuai. Mereka berinteraksi dengan layanan untuk melakukan logika bisnis.
*   `auth_controller.go`: Mengelola autentikasi pengguna, termasuk registrasi, login, dan logout.
*   `base_controller.go`: Menyediakan fungsionalitas umum atau struktur dasar untuk pengontrol lain.
*   `car_controller.go`: Menangani operasi terkait manajemen mobil (misalnya, menambah, mengedit, menghapus mobil, melihat daftar mobil).
*   `comment_controller.go`: Mengelola komentar pada mobil.
*   `home_controller.go`: Menangani halaman arahan utama dan fungsionalitas pencarian mobil.
*   `order_controller.go`: Mengelola pemesanan mobil dan operasi terkait pesanan.

#### `middlewares/`
Middleware adalah fungsi yang memproses permintaan HTTP sebelum mencapai penangan rute yang sebenarnya. Mereka digunakan untuk tugas-tugas seperti autentikasi, logging, dll.
*   `auth_middleware.go`: Mengimplementasikan middleware autentikasi untuk melindungi rute yang memerlukan pengguna yang masuk atau peran tertentu (misalnya, admin).

#### `models/`
Direktori ini mendefinisikan struktur data (struct) yang mewakili entitas dalam domain aplikasi, seringkali memetakan langsung ke tabel database.
*   `car_image.go`: Mendefinisikan struktur untuk gambar mobil, biasanya ditautkan ke model `Car`.
*   `car.go`: Mendefinisikan model `Car`, mewakili detail mobil seperti nama, kapasitas, transmisi, harga, dan ketersediaan.

#### `repositories/`
Repositori mengabstraksi lapisan data, menyediakan metode untuk berinteraksi dengan database. Mereka merangkum logika penyimpanan dan pengambilan data.
*   `car_repository.go`: Menyediakan metode untuk operasi CRUD (Create, Read, Update, Delete) pada data `Car` di database.
*   `comment_repository.go`: Menyediakan metode untuk mengelola data `Comment`.
*   `order_repository.go`: Menyediakan metode untuk mengelola data `Order`.
*   `user_repository.go`: Menyediakan metode untuk mengelola data `User`.

#### `routes/`
Direktori ini mendefinisikan rute URL aplikasi dan memetakannya ke fungsi pengontrol yang sesuai.
*   `routes.go`: Memusatkan semua definisi rute, mengaitkan jalur URL dengan metode pengontrol tertentu dan menerapkan middleware yang diperlukan.

#### `services/`
Layanan merangkum logika bisnis aplikasi. Mereka mengatur interaksi antara repositori dan melakukan operasi kompleks.
*   `car_service.go`: Berisi logika bisnis terkait mobil, seperti mencari, memfilter, dan mengelola ketersediaan mobil.
*   `comment_service.go`: Berisi logika bisnis untuk komentar.
*   `order_service.go`: Berisi logika bisnis untuk pemesanan mobil, termasuk perhitungan harga, pembaruan status, dan integrasi pembayaran.
*   `user_service.go`: Berisi logika bisnis untuk manajemen pengguna, seperti registrasi, login, dan manajemen profil.

#### `static/`
Direktori ini menyajikan aset statis seperti CSS, JavaScript, dan file yang diunggah.
*   `css/`: Berisi file CSS.
    *   `input.css`: File CSS sumber, kemungkinan digunakan dengan Tailwind CSS untuk pengembangan.
    *   `output.css`: File CSS yang dikompilasi dan diminifikasi yang dihasilkan oleh Tailwind CSS, digunakan dalam produksi.
*   `js/`: Berisi file JavaScript.
    *   `theme.js`: JavaScript untuk menangani fungsionalitas pengalih tema gelap/terang.
*   `uploads/`: Direktori untuk konten yang diunggah pengguna.
    *   `cars/`: Khusus untuk gambar mobil yang diunggah.
*   `package-lock.json`: Mencatat versi pasti dependensi npm.
*   `package.json`: Mendefinisikan metadata dan dependensi proyek npm (misalnya, Tailwind CSS).
*   `tailwind.config.js`: File konfigurasi untuk Tailwind CSS, mendefinisikan tema kustom, plugin, dll.

#### `templates/`
Direktori ini menyimpan semua file template HTML yang dirender oleh aplikasi Go.
*   `auth/`: Template terkait autentikasi pengguna.
    *   `login.html`: Formulir login.
    *   `register.html`: Formulir registrasi.
*   `car/`: Template untuk car management (admin side).
    *   `add.html`: Formulir untuk menambah mobil baru.
    *   `edit.html`: Formulir untuk mengedit mobil yang sudah ada.
    *   `index.html`: Daftar mobil (tampilan admin).
*   `partials/`: Cuplikan HTML yang dapat digunakan kembali yang disertakan dalam beberapa template.
    *   `footer.html`: Bagian footer aplikasi.
    *   `navbar.html`: Bilah navigasi aplikasi.
    *   `test.html`: Template placeholder atau pengujian.
*   `admin_orders.html`: Tampilan admin untuk mengelola semua pesanan penyewaan mobil.
*   `booking.html`: Formulir bagi pengguna untuk memesan mobil.
*   `detail.html`: Tampilan detail satu mobil.
*   `home.html`: Halaman arahan utama yang menampilkan mobil yang tersedia.
*   `payment.html`: Halaman untuk pemrosesan pembayaran QRIS.

#### `tests/`
Direktori ini ditujukan untuk pengujian unit dan integrasi. (Saat ini kosong, menunjukkan area untuk pengembangan di masa mendatang).

#### `utils/`
Direktori ini berisi fungsi utilitas dan modul pembantu yang digunakan di seluruh aplikasi.
*   `excel_exporter.go`: Utilitas untuk mengekspor data (misalnya, daftar pesanan) ke format Excel.
*   `pdf_exporter.go`: Utilitas untuk mengekspor data (misalnya, detail pesanan) ke format PDF.
*   `qris.go`: Berisi logika untuk pembuatan dan pemrosesan pembayaran QRIS (Quick Response Code Indonesian Standard).
*   `template.go`: Fungsi pembantu untuk merender template HTML.

### File Root

*   `go.mod`: Mendefinisikan jalur modul dan mencantumkan dependensi langsung dan tidak langsung dari proyek Go.
*   `go.sum`: Berisi checksum kriptografi dari dependensi modul, memastikan integritas.
*   `README.md`: File dokumentasi ini.
*   `rental.db`: File database SQLite tempat semua data aplikasi disimpan.

## Memulai

### Prasyarat

*   Go (versi 1.24.4 atau lebih tinggi)
*   Node.js dan npm (untuk Tailwind CSS)

### Instalasi

1.  **Kloning repositori:**
    ```bash
    git clone https://github.com/nabilulilalbab/rental-mobil.git
    cd rental-mobil
    ```
2.  **Instal dependensi Go:**
    ```bash
    go mod tidy
    ```
3.  **Instal dependensi Node.js (untuk Tailwind CSS):**
    ```bash
    cd static
    npm install
    cd ..
    ```
4.  **Bangun Tailwind CSS:**
    ```bash
    npx tailwindcss -i ./static/css/input.css -o ./static/css/output.css --minify
    ```

### Menjalankan Aplikasi

1.  **Jalankan backend Go:**
    ```bash
    go run cmd/main.go
    ```
2.  Buka browser web Anda dan navigasikan ke `http://localhost:8080` (atau port apa pun tempat aplikasi dimulai).

## Berkontribusi

Kontribusi dipersilakan! Silakan ikuti langkah-langkah berikut:

1.  Fork repositori.
2.  Buat cabang baru (`git checkout -b feature/nama-fitur-anda`).
3.  Buat perubahan Anda.
4.  Tulis tes untuk perubahan Anda (jika berlaku).
5.  Pastikan semua tes lulus (`go test ./...`).
6.  Commit perubahan Anda (`git commit -m 'feat: Tambah fitur baru'`).
7.  Push ke cabang (`git push origin feature/nama-fitur-anda`).
8.  Buka Pull Request.

## Lisensi

Proyek ini dilisensikan di bawah Lisensi MIT. Lihat file `LICENSE` untuk detailnya.
