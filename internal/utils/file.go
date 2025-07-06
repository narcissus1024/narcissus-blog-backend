package utils

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/chai2010/webp"
	"github.com/disintegration/imaging"
	"github.com/google/uuid"
)

func ImageBytes2WebpBytes(input io.Reader, quality float32) ([]byte, error) {
	img, err := imaging.Decode(input)
	if err != nil {
		return nil, err
	}

	webpBytes, err := webp.EncodeRGBA(img, quality)
	if err != nil {
		return nil, err
	}

	return webpBytes, nil
}

func SaveFileBytes(data []byte, imagePath string) error {
	if err := os.MkdirAll(filepath.Dir(imagePath), 0750); err != nil {
		return err
	}
	return os.WriteFile(imagePath, data, 0644)
}

func FileName2RandomName(fileName string) string {
	if len(strings.TrimSpace(fileName)) == 0 {
		return fileName
	}
	randomStr := strings.ReplaceAll(uuid.New().String(), "-", "")
	split := strings.Split(fileName, ".")
	if len(split) < 2 {
		return randomStr
	}
	return randomStr + "." + split[len(split)-1]
}

func ConvertFileNameExt(fileName string, ext string) string {
	dir, baseFile := filepath.Split(fileName)
	split := strings.Split(baseFile, ".")
	if len(split) < 2 {
		return filepath.Join(dir, baseFile+"."+ext)
	}
	return filepath.Join(dir, split[0]+"."+ext)
}

func GetRootDir() string {
	ex, _ := os.Executable()
	exPath := filepath.Dir(ex)
	realPath, _ := filepath.EvalSymlinks(exPath)
	return realPath
}
