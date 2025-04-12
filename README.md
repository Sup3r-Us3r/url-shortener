# 🔗 Encrypted URL Shortener (Go)

This is a simple URL shortener project built with Go. It uses AES encryption to securely store URLs and generates a unique short identifier for redirection.

---

## 📦 Features

- ✅ Shortens URLs with `http` or `https` support
- 🔐 Encrypts stored URLs using AES-CTR
- 🔁 Redirects to the original URL using the short ID
- 📦 In-memory storage with `map[string]string`
- 🔒 Secure short ID generation using `crypto/rand`

---

## 🚀 Getting Started

### 1. Clone the repository

```bash
$ git clone https://github.com/Sup3r-Us3r/url-shortener.git
$ cd url-shortener
```

### 2. Run the project

```bash
$ go run main.go
```

Server will be running at:

```
http://localhost:8080
```

---

## 📌 API Endpoints

### `GET /shorten?url=<your_url>`

Shortens a valid URL (must start with `http://` or `https://`).

#### Query Parameters:

- `url`: The original URL you want to shorten

#### Example:

```bash
$ curl "http://localhost:8080/shorten?url=https://example.com"
```

#### Response:

```json
{
  "originalUrl": "https://example.com",
  "shortUrl": "http://localhost:8080/ABC123"
}
```

---

### `GET /<shortId>`

Redirects to the original URL associated with the given short ID.

#### Example:

```bash
$ curl -v http://localhost:8080/ABC123
```

---

## 🔐 Encryption

The URLs are encrypted using AES (CTR mode) with a 24-byte key:

```go
secretKey = []byte("superSecretKey1234567890")
```

> ⚠️ **Note:** In production environments, store the `secretKey` securely using environment variables or a secrets manager.

---

## ⚙️ Tech Stack

- [Go](https://golang.org/) — Main language
- `crypto/aes`, `crypto/cipher`, `crypto/rand` — Encryption
- `net/http` — Native HTTP server
- `encoding/json`, `encoding/hex` — Serialization and encoding

---

## 📝 License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
