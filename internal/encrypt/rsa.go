package encrypt

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"
	"path/filepath"

	"github.com/narcissus1949/narcissus-blog/internal/utils"
)

const (
	PrivateKeyName    = "private.pem"
	PublicKeyName     = "public.pem"
	PEMPublicKeyType  = "RSA PUBLIC KEY"
	PEMPrivateKeyType = "RSA PRIVATE KEY"
)

func RSAGenKey(bits int, fileDir string) error {
	// 生成私钥文件
	// GenerateKey函数使用随机数据生成器random生成一对具有指定字位数的RSA密钥
	// 参数1: Reader是一个全局、共享的密码用强随机数生成器
	// 参数2: 秘钥的位数 - bit
	privateKey, generateErr := rsa.GenerateKey(rand.Reader, bits)
	if generateErr != nil {
		return generateErr
	}
	// MarshalPKCS1PrivateKey将rsa私钥序列化为ASN.1 PKCS#1 DER编码
	derStream := x509.MarshalPKCS1PrivateKey(privateKey)
	// Block代表PEM编码的结构, 对其进行设置
	block := pem.Block{
		Type:  PEMPrivateKeyType,
		Bytes: derStream,
	}
	// 创建文件
	if _, statErr := os.Stat(fileDir); statErr != nil {
		if os.IsNotExist(statErr) {
			if makedirERR := os.MkdirAll(fileDir, 0755); makedirERR != nil {
				return makedirERR
			}
		} else {
			return statErr
		}
	}

	privFile, createErr := os.Create(filepath.Join(fileDir, PrivateKeyName))
	if createErr != nil {
		return createErr
	}
	// 关闭文件
	defer privFile.Close()
	// 使用pem编码, 并将数据写入文件中
	if encodeErr := pem.Encode(privFile, &block); encodeErr != nil {
		return encodeErr
	}

	// 生成公钥文件
	publicKey := privateKey.PublicKey
	derPkix := x509.MarshalPKCS1PublicKey(&publicKey)
	block = pem.Block{
		Type:  PEMPublicKeyType,
		Bytes: derPkix,
	}
	pubFile, createPublicErr := os.Create(filepath.Join(fileDir, PublicKeyName))
	if createPublicErr != nil {
		return createPublicErr
	}
	// 编码公钥, 写入文件
	if encodeErr := pem.Encode(pubFile, &block); encodeErr != nil {
		return encodeErr
	}
	defer pubFile.Close()

	return nil
}

func RSAEncrypt(src []byte, fileDir string) ([]byte, error) {
	pubKey, readErr := RSAReadPublicKey(fileDir)
	if readErr != nil {
		return nil, readErr
	}

	// 公钥加密
	result, encryptErr := rsa.EncryptPKCS1v15(rand.Reader, pubKey, src)
	if encryptErr != nil {
		return nil, encryptErr
	}
	return result, nil
}

func RSAEncryptWithBase64(src []byte, fileDir string) (string, error) {
	result, encryptErr := RSAEncrypt(src, fileDir)
	if encryptErr != nil {
		return "", encryptErr
	}

	return utils.Base64Encode(result), nil
}

func RSADecrypt(src []byte, fileDir string) ([]byte, error) {
	privateKey, readErr := RSAReadPrivateKey(fileDir)
	if readErr != nil {
		return nil, readErr
	}
	// 私钥解密
	result, decryptErr := rsa.DecryptPKCS1v15(rand.Reader, privateKey, src)
	if decryptErr != nil {
		return nil, decryptErr
	}

	return result, nil
}

func RSADecryptWithBase64(src []byte, fileDir string) ([]byte, error) {
	decodeSrc, decodeErr := utils.Base64Decode(string(src))
	if decodeErr != nil {
		return nil, decodeErr
	}
	return RSADecrypt(decodeSrc, fileDir)
}

func RSAReadPublicKey(fileDir string) (*rsa.PublicKey, error) {
	allText, readErr := os.ReadFile(filepath.Join(fileDir, PublicKeyName))
	if readErr != nil {
		return nil, readErr
	}
	// 从数据中查找到下一个PEM格式的块
	block, _ := pem.Decode(allText)
	if block == nil {
		return nil, errors.New("public key decode empty")
	}
	// 解析一个pem格式的公钥
	publicKey, parseErr := x509.ParsePKCS1PublicKey(block.Bytes)
	if parseErr != nil {
		return nil, parseErr
	}
	return publicKey, nil
}

func RSAReadPrivateKey(fileDir string) (*rsa.PrivateKey, error) {
	allText, readErr := os.ReadFile(filepath.Join(fileDir, PrivateKeyName))
	if readErr != nil {
		return nil, readErr
	}
	// 从数据中查找到下一个PEM格式的块
	block, _ := pem.Decode(allText)
	if block == nil {
		return nil, errors.New("private key decode empty")
	}
	// 解析一个pem格式的私钥
	privateKey, parseErr := x509.ParsePKCS1PrivateKey(block.Bytes)
	if parseErr != nil {
		return nil, parseErr
	}
	return privateKey, nil
}

func RSAPublicKey2Mem(publicKey *rsa.PublicKey) ([]byte, error) {
	bytes := x509.MarshalPKCS1PublicKey(publicKey)
	publicKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  PEMPublicKeyType,
		Bytes: bytes,
	})

	return publicKeyPEM, nil
}
