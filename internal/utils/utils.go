package utils

import (
	"encoding/base64"
	"errors"
	"mime/multipart"
	"net/http"
	"strings"
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
