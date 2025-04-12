package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"strings"
	"sync"
)

var (
	urlStore    = make(map[string]string)
	mu          sync.Mutex
	secretKey   = []byte("superSecretKey1234567890")
	lettersRune = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
)

func encrypt(originalUrl string) string {
	block, err := aes.NewCipher(secretKey)
	if err != nil {
		log.Fatal(err)
	}

	plainText := []byte(originalUrl)
	cipherText := make([]byte, aes.BlockSize+len(plainText))

	iv := cipherText[:aes.BlockSize]

	if _, err := rand.Read(iv); err != nil {
		log.Fatal(err)
	}

	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], plainText)

	return hex.EncodeToString(cipherText)
}

func decrypt(encryptedUrl string) string {
	block, err := aes.NewCipher(secretKey)
	if err != nil {
		log.Fatal(err)
	}

	cipherText, err := hex.DecodeString(encryptedUrl)
	if err != nil {
		log.Fatal(err)
	}

	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(cipherText, cipherText)

	return string(cipherText)
}

func generateShortId() string {
	b := make([]rune, 6)
	for i := range b {
		number, err := rand.Int(rand.Reader, big.NewInt(int64(len(lettersRune))))
		if err != nil {
			log.Fatal(err)
		}

		b[i] = lettersRune[number.Int64()]
	}

	return string(b)
}

func shortenUrlHandler(w http.ResponseWriter, r *http.Request) {
	originalUrl := r.URL.Query().Get("url")

	if originalUrl == "" {
		http.Error(w, `The query param "url" is required`, http.StatusBadRequest)
		return
	}

	if !(strings.HasPrefix(originalUrl, "https://") || strings.HasPrefix(originalUrl, "http://")) {
		http.Error(w, "The url must have the value https:// or http://", http.StatusBadRequest)
		return
	}

	encryptedUrl := encrypt(originalUrl)
	shortId := generateShortId()

	mu.Lock()
	urlStore[shortId] = encryptedUrl
	mu.Unlock()

	shortUrl := fmt.Sprintf("http://localhost:8080/%s", shortId)

	type response struct {
		OriginalUrl string `json:"originalUrl"`
		ShortUrl    string `json:"shortUrl"`
	}

	responseData := response{
		OriginalUrl: originalUrl,
		ShortUrl:    shortUrl,
	}

	json.NewEncoder(w).Encode(responseData)
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	shortId := r.URL.Path[1:]

	mu.Lock()
	encryptedUrl, ok := urlStore[shortId]
	mu.Unlock()

	if !ok {
		http.Error(w, "URL not found", http.StatusNotFound)
		return
	}

	decryptedUrl := decrypt(encryptedUrl)
	http.Redirect(w, r, decryptedUrl, http.StatusFound)
}

func main() {
	http.HandleFunc("/shorten", shortenUrlHandler)
	http.HandleFunc("/", redirectHandler)

	fmt.Println("SERVER IS RUNNING AT :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
