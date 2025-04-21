# INSIS_REST-API

Mata Kuliah Integrasi Sistem 2025

Proyek 1 REST API - GO-Urban

## Kelompok I

|      **Nama**       |  **NRP**   |
| :-----------------: | :--------: |
|    Abhirama T.H     | 5027231061 |
|        Hasan        | 5027231073 |
| Nabiel Nizar Anwari | 5027231087 |



# GO-Urban: Simulasi Kota Cerdas dengan gRPC dan Go

Selamat datang di **GO-Urban**, sebuah proyek simulasi pengelolaan kota cerdas yang dibangun dengan teknologi **Go** dan **gRPC**! Proyek ini memungkinkan Anda mengontrol lalu lintas, memantau kualitas udara, dan mengelola unit darurat seperti ambulans melalui antarmuka web yang interaktif. Dengan konteks lokal Indonesia, GO-Urban menawarkan pengalaman realistis dan menarik untuk memahami teknologi *smart city* dan komunikasi gRPC.

\
*Gambar: Antarmuka web GO-Urban untuk mengelola kota cerdas.*

## Apa Itu GO-Urban?

GO-Urban adalah simulasi sistem *smart city* yang menunjukkan bagaimana teknologi dapat digunakan untuk mengelola kota secara efisien. Proyek ini menggunakan **gRPC** untuk komunikasi antara server dan client, serta **WebSocket** untuk streaming data real-time. Fitur utamanya meliputi:

- **Status Lalu Lintas (Unary RPC)**: Cek kemacetan dan jumlah kendaraan di jalan seperti "Jalan Sudirman".
- **Kualitas Udara (Server Streaming RPC)**: Pantau polusi di zona seperti "Jakarta Pusat" secara real-time.
- **Atur Lampu Lalu Lintas (Client Streaming RPC)**: Kirim perintah untuk mengatur lampu di persimpangan seperti "Simpang Senayan" dengan durasi detik (misalnya, 10, 9, ..., 1 detik).
- **Kontrol Darurat (Bidirectional Streaming RPC)**: Kelola unit darurat seperti ambulans, dari "Pusat" ke "Perjalanan" hingga tiba di tujuan, dengan waktu tempuh 1 menit.

Semua fitur diakses melalui antarmuka web yang ramah pengguna, dengan feedback dalam bahasa Indonesia untuk konteks lokal.

## Tujuan Proyek

GO-Urban dirancang untuk:

- Menyediakan simulasi *smart city* yang interaktif dan realistis.
- Mendemonstrasikan pola komunikasi gRPC: Unary, Server Streaming, Client Streaming, dan Bidirectional Streaming.
- Memberikan pengalaman belajar tentang Go, gRPC, dan WebSocket dalam aplikasi nyata.
- Menawarkan antarmuka yang mudah digunakan untuk masyarakat Indonesia,

## Struktur Direktori

Berikut adalah struktur direktori proyek GO-Urban:

```
GO-Urban/
├── smartcity/
│   ├── smartcity.pb.go         # File proto yang dihasilkan (definisi gRPC)
│   └── smartcity_grpc.pb.go    # File proto yang dihasilkan (implementasi gRPC)
├── static/
│   └── index.html              # Antarmuka web (HTML, CSS, JavaScript)
├── smartcity.proto             # Definisi layanan gRPC
├── server.go                   # Server gRPC untuk menangani logika kota cerdas
├── webserver.go                # Server web untuk antarmuka dan komunikasi WebSocket
└── go.mod                      # Modul dependensi Go
```

- `smartcity/`: Berisi file yang dihasilkan dari `smartcity.proto` untuk definisi dan implementasi gRPC.
- `static/index.html`: Halaman web utama dengan antarmuka untuk menguji semua fitur.
- `smartcity.proto`: Mendefinisikan layanan gRPC seperti `GetTrafficStatus`, `StreamAirQuality`, `SetTrafficLights`, dan `EmergencyControl`.
- `server.go`: Menjalankan server gRPC di port `50051`, menangani logika seperti pergerakan unit darurat dan pengaturan lampu lalu lintas.
- `webserver.go`: Menjalankan server web di port `8081`, menghubungkan antarmuka web dengan server gRPC melalui WebSocket.
- `go.mod`: Mengelola dependensi seperti `gin-gonic/gin` dan `gorilla/websocket`.

## Prasyarat

Untuk menjalankan GO-Urban, Anda memerlukan:

- **Go** (versi 1.16 atau lebih baru): [Unduh Disini](https://go.dev/doc/install)
- **protoc**: Kompilator Protocol Buffers untuk menghasilkan file gRPC. Petunjuk instalasi.
- Browser modern (Chrome, Firefox, dll.) untuk mengakses antarmuka web.

## Cara Instalasi

1. **Kloning Repositori**:

   ```bash
   git clone https://github.com/<username>/GO-Urban.git
   cd GO-Urban
   ```

2. **Instal Dependensi**:

   ```bash
   go mod tidy
   go get github.com/gin-gonic/gin
   go get github.com/gorilla/websocket
   ```

3. **Hasilkan File gRPC** (jika `smartcity.pb.go` belum ada):

   ```bash
   protoc --go_out=. --go-grpc_out=. smartcity.proto
   ```

## Cara Menjalankan

1. **Jalankan Server gRPC**:

   - Buka terminal di direktori `GO-Urban`.
   - Jalankan:

     ```bash
     go run server.go
     ```
   - Output:

     ```
     Server dimulai di port 50051
     ```

2. **Jalankan Server Web**:

   - Buka terminal lain.
   - Pastikan port 8081 bebas:

     ```bash
     netstat -aon | findstr :8081
     taskkill /PID <PID> /F
     ```
   - Jalankan:

     ```bash
     go run webserver.go
     ```
   - Output:

     ```
     Web server started on :8081
     ```

3. **Akses Antarmuka Web**:

   - Buka browser dan kunjungi: `http://localhost:8081`.
   - Anda akan melihat antarmuka dengan empat bagian: Unary RPC, Server Streaming, Client Streaming, dan Bidirectional Streaming.

## Cara Penggunaan

Berikut adalah panduan untuk menggunakan fitur-fitur GO-Urban melalui antarmuka web:

### 1. Dapatkan Status Lalu Lintas (Unary RPC)

- **Deskripsi**: Cek status kemacetan dan jumlah kendaraan di jalan tertentu.
- **Cara Pakai**:
  - Masukkan **ID Jalan** (misalnya, "Jalan Sudirman").
  - Klik "Dapatkan Status".
  - Contoh output:

    ```
    Status Lalu Lintas: Jalan Sudirman - Sedang, Kendaraan: 45
    ```

### 2. Streaming Kualitas Udara (Server Streaming RPC)

- **Deskripsi**: Pantau data kualitas udara di zona tertentu secara real-time.
- **Cara Pakai**:
  - Masukkan **ID Zona** (misalnya, "Jakarta Pusat").
  - Klik "Mulai Streaming".
  - Data polusi akan muncul setiap detik, misalnya:

    ```
    Kualitas Udara di Jakarta Pusat: 65.23 pada 2025-04-20T10:00:00Z
    ```
  - Klik "Hentikan Streaming" untuk berhenti.

### 3. Atur Lampu Lalu Lintas (Client Streaming RPC)

- **Deskripsi**: Atur lampu lalu lintas di persimpangan dengan durasi detik, dikirim berulang hingga habis.
- **Cara Pakai**:
  - Klik "Tambah Perintah" untuk membuka modal.
  - Masukkan **ID Persimpangan** (misalnya, "Simpang Senayan") dan **Durasi** (misalnya, `10`).
  - Tambah perintah lain jika perlu, lalu klik "Kirim Perintah".
  - Perintah dikirim per detik (10, 9, ..., 1), misalnya:

    ```
    Mengirim perintah: Simpang Senayan - 10 second
    Mengirim perintah: Simpang Senayan - 9 second
    ...
    Mengirim perintah: Simpang Senayan - 1 second
    Semua perintah telah diproses
    ```

### 4. Kontrol Darurat (Bidirectional Streaming RPC)

- **Deskripsi**: Kelola unit darurat seperti ambulans dengan perintah dan feedback real-time.
- **Cara Pakai**:
  - Klik "Mulai Kontrol Darurat" untuk membuka modal.
  - Masukkan **ID Unit** (misalnya, "Ambulans 01") dan **Perintah** (misalnya, "Kirim Bantuan ke Lokasi Bencana").
  - Klik "Kirim Perintah".
  - Unit akan bergerak dari "Pusat" ke "Perjalanan" hingga tiba di tujuan (1 menit), dengan feedback seperti:

    ```
    Feedback: Ambulans 01 - Unit memulai perjalanan ke Lokasi Bencana, perkiraan waktu: 60s di Pusat
    Feedback: Ambulans 01 - Unit dalam perjalanan menuju Lokasi Bencana, waktu tersisa: 59s di Perjalanan
    ...
    Feedback: Ambulans 01 - Unit telah tiba di tujuan di Lokasi Bencana
    ```
  - Perintah tambahan:
    - `check_status`: Cek status unit (misalnya, "Unit dalam perjalanan menuju Lokasi Bencana, waktu tersisa: 45s").
    - `priority_mode`: Percepat waktu tempuh (misalnya, dari 60 detik jadi 30 detik).
  - Klik "Hentikan Streaming" untuk berhenti.

## Teknologi yang Digunakan

- **Go**: Bahasa pemrograman untuk server gRPC dan web.
- **gRPC**: Framework untuk komunikasi RPC yang efisien.
- **WebSocket**: Digunakan untuk streaming data real-time di antarmuka web.
- **Gin**: Framework web Go untuk menangani rute dan WebSocket.
- **Bootstrap**: Untuk antarmuka web yang responsif dan menarik.
- **Protocol Buffers**: Mendefinisikan layanan dan pesan gRPC.

## Penjelasan Fitur Utama

### Pergerakan Unit Darurat

- Unit darurat (misalnya, ambulans) bergerak melalui tiga fase: **Pusat -&gt; Perjalanan -&gt; Tujuan**.
- Waktu tempuh diatur 1 menit (60 detik), dengan pembaruan status setiap detik.
- Pengguna dapat mengirim perintah seperti `priority_mode` untuk mempercepat perjalanan.

### Pengaturan Lampu Lalu Lintas

- Perintah dikirim berulang per detik berdasarkan durasi input (misalnya, 10, 9, ..., 1 detik).
- Mendukung beberapa persimpangan, diproses secara berurutan.
- Feedback real-time menunjukkan setiap perintah yang dikirim.

### Konteks Lokal

- Nama jalan dan persimpangan menggunakan Bahasa Indonesia (misalnya, "Jalan Sudirman", "Simpang Senayan").
- Semua feedback dalam bahasa Indonesia untuk kemudahan pengguna lokal.


## Lisensi

Proyek ini dilisensikan di bawah Kelompok I (Nabiel,Hasan,Abhi) Mata Kuliah Integrasi Sistem Departemen Teknologi Informasi License. Silakan gunakan, modifikasi, dan distribusikan sesuai kebutuhan.

---

**GO-Urban**: Mengelola kota cerdas dengan teknologi modern, satu perintah pada satu waktu! 🚦🚑
