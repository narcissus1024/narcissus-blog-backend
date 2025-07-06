package service

import (
	"errors"
	"fmt"
	"mime/multipart"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/narcissus1949/narcissus-blog/cmd/config"
	"github.com/narcissus1949/narcissus-blog/internal/encrypt"
	cerr "github.com/narcissus1949/narcissus-blog/internal/error"
	"github.com/narcissus1949/narcissus-blog/internal/utils"
	"github.com/narcissus1949/narcissus-blog/pkg/dto"
	"github.com/narcissus1949/narcissus-blog/pkg/vo"
	"go.uber.org/zap"
)

var (
	CommonServiceInstance = new(commoneService)
	allowedMimeTypes      = map[string]string{
		"image/jpeg": "jpeg",
		"image/png":  "png",
		"image/gif":  "gif",
	}
	allowedImageExt              = []string{"jpg", "jpeg", "png", "gif"}
	compressImageThreshold int64 = 2 << 20 // 2MB
)

type commoneService struct {
}

func (s *commoneService) UploadImage(ctx *gin.Context, file *multipart.FileHeader) (vo.UploadImageVo, *cerr.Error) {
	var resp vo.UploadImageVo

	if err := checkImageFile(file); err != nil {
		zap.L().Error("Failed to check image file", zap.Error(err))
		return resp, cerr.NewParamError()
	}

	// 图片转为webp格式
	f, openErr := file.Open()
	if openErr != nil {
		zap.L().Error("Failed to open image file", zap.Error(openErr))
		return resp, cerr.NewSysError()
	}
	defer f.Close()

	webpImageByte, convertImgErr := utils.ImageBytes2WebpBytes(f, 85)
	if convertImgErr != nil {
		zap.L().Error("Failed to convert image to webp", zap.Error(convertImgErr))
		return resp, cerr.NewSysError()
	}

	// 存储图片
	imgRelativePath := path.Join(
		time.Now().Format("2006"),
		utils.FileName2RandomName(utils.ConvertFileNameExt(file.Filename, "webp")),
	)
	savePath := path.Join(config.Config.App.ImgDataDir, imgRelativePath)

	if err := utils.SaveFileBytes(webpImageByte, savePath); err != nil {
		zap.L().Error("Failed to save image file", zap.Error(err))
		return resp, cerr.NewSysError()
	}

	resp.ImgUrl = config.Config.App.ImgProxyURL + "/" + imgRelativePath
	return resp, nil
}

func (s *commoneService) GetRASPublicKey(ctx *gin.Context) (vo.RASPublicKeyVo, *cerr.Error) {
	var resp vo.RASPublicKeyVo
	publicKey, err := encrypt.RSAReadPublicKey(config.Config.App.PublicKeyDir)
	if err != nil {
		zap.L().Error("Failed to read public key", zap.Error(err))
		return resp, cerr.NewSysError()
	}
	publicKeyPEM, err := encrypt.RSAPublicKey2Mem(publicKey)
	if err != nil {
		zap.L().Error("Failed to convert public key to PEM", zap.Error(err))
		return resp, cerr.NewSysError()
	}
	resp.PublicKey = string(publicKeyPEM)
	return resp, nil
}

func (s *commoneService) PublicKeyEncrypt(req dto.PublicKeyEncrypDto) (vo.PublicKeyEncryptVo, *cerr.Error) {
	var resp vo.PublicKeyEncryptVo
	encryptedData, err := encrypt.RSAEncryptWithBase64([]byte(req.Data), config.Config.App.PublicKeyDir)
	if err != nil {
		zap.L().Error("Failed to encrypt data with public key", zap.Error(err))
		return resp, cerr.NewSysError()
	}
	resp.EncryptedData = string(encryptedData)
	return resp, nil
}

func checkImageFile(file *multipart.FileHeader) error {
	if file == nil || file.Filename == "" || file.Size == 0 {
		return errors.New("image is empty")
	}

	// 检查文件扩展名
	if err := checkImageExt(file.Filename); err != nil {
		return err
	}

	// 检查文件mimeType
	if err := checkImageMimeType(file); err != nil {
		return err
	}

	return nil
}

func checkImageExt(filename string) error {
	ext := filepath.Ext(filename)
	if len(ext) == 0 {
		return errors.New("does not support image extension")
	}
	ext = strings.TrimPrefix(ext, ".")
	extVliade := false
	for i := range allowedImageExt {
		if ext == allowedImageExt[i] {
			extVliade = true
			break
		}
	}
	if !extVliade {
		return fmt.Errorf("does not support image extension: %s", ext)
	}
	return nil
}

func checkImageMimeType(file *multipart.FileHeader) error {
	mimeType, getMimeTypeErr := utils.GetRealMimeType(file)
	if getMimeTypeErr != nil {
		return getMimeTypeErr
	}
	if _, ok := allowedMimeTypes[mimeType]; !ok {
		return fmt.Errorf("does not support image format: %s", mimeType)
	}
	return nil
}
