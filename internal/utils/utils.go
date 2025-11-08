package utils

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"mime/multipart"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

func ParseAuthorization(authorization string) (string, error) {
	parts := strings.Split(authorization, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", errors.New("authorization invalide")
	}

	return parts[1], nil
}

func GetRealMimeType(file *multipart.FileHeader) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	buf := make([]byte, 512)
	_, err = src.Read(buf)
	if err != nil {
		return "", err
	}

	mimeType := http.DetectContentType(buf)
	return mimeType, nil
}

func Base64Encode(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

func Base64Decode(encodeData string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(encodeData)
}

func ArrayExistInt(array []int, target int) bool {
	for _, v := range array {
		if v == target {
			return true
		}
	}
	return false
}

func GenerateUUID() string {
	return strings.ReplaceAll(uuid.New().String(), "-", "")
}

func GenerateTempUserID(ip, userAgent string) string {
	h := sha1.New()
	h.Write([]byte(ip + userAgent))
	return hex.EncodeToString(h.Sum(nil))
}

func Retry(count int, delayMill int, fn func() error) error {
	var err error
	for i := 0; i < count; i++ {
		if err = fn(); err == nil {
			return nil
		}
		time.Sleep(time.Duration(delayMill) * time.Millisecond)
	}
	return err
}
