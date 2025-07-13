package dto

import (
	"errors"
	"regexp"
	"strings"

	"github.com/narcissus1949/narcissus-blog/cmd/blog/app/config"
	"github.com/narcissus1949/narcissus-blog/internal/encrypt"
	"github.com/narcissus1949/narcissus-blog/internal/jwt"
)

type SignIn struct {
	Username    string `json:"username" binding:"required,no_spacing,gte=5,lte=20"`
	Nickname    string `json:"nickname" binding:"required,no_spacing,gte=5,lte=10"`
	Password    string `json:"password" binding:"required"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	AvatarPath  string `json:"avatar_path"`
}

func (r *SignIn) Validate() error {
	// 英文、字母、下划线、5-20
	userNamePattern := `^[a-zA-Z0-9_]{5,20}$`
	// 英文、字母、下划线、中文、5-10
	nickNamePattern := `^[\w\p{Han}]{5,10}$`
	// 字母、数字、特殊字符、8-64
	passwdPattern := `^[a-zA-Z0-9_!@#$%^&.*]{8,64}$`

	// 账号
	if match, err := regexp.MatchString(userNamePattern, r.Username); err != nil {
		return err
	} else if !match {
		return errors.New("username invalid")
	}

	// 昵称
	if match, err := regexp.MatchString(nickNamePattern, r.Nickname); err != nil {
		return err
	} else if !match {
		return errors.New("nickname invalid")
	}

	// 密码
	if len(strings.TrimSpace(r.Password)) <= 0 {
		return errors.New("password invalid")
	}
	passwdDecrypt, decryptErr := encrypt.RSADecryptWithBase64([]byte(r.Password), config.Config.App.PrivateKeyDir)
	if decryptErr != nil {
		return decryptErr
	}

	if match, err := regexp.MatchString(passwdPattern, string(passwdDecrypt)); err != nil {
		return err
	} else if !match {
		return errors.New("password invalid")
	}

	return nil
}

type LoginDto struct {
	Username         string `json:"username" binding:"required,no_spacing,gte=5,lte=20"`
	Password         string `json:"password" binding:"required,gte=5"`
	VerificationCode string `json:"verification_code"`
}

type LogoutDto struct {
	AccessToken  string `json:"access_token" binding:"required"`
	RefreshToken string `json:"refresh_token" binding:"required"`
}

func (r *LogoutDto) Validate() error {
	if len(strings.TrimSpace(r.AccessToken)) <= 0 || len(strings.TrimSpace(r.RefreshToken)) <= 0 {
		return errors.New("access token or refresh token invalid")
	}
	return nil
}

type RefreshTokenDto struct {
	AccessToken  string `json:"access_token" binding:"required"`
	RefreshToken string `json:"refresh_token" binding:"required"`
}

func (r *RefreshTokenDto) Validate() error {
	if len(strings.TrimSpace(r.AccessToken)) <= 0 || len(strings.TrimSpace(r.RefreshToken)) <= 0 {
		return errors.New("access token or refresh token invalid")
	}

	if expire, err := jwt.IsTokenExpire(r.AccessToken); err != nil {
		return err
	} else if !expire {
		return errors.New("access token does not expired")
	}

	if expire, err := jwt.IsTokenExpire(r.RefreshToken); err != nil {
		return err
	} else if expire {
		return errors.New("refresh token has expired")
	}
	return nil
}
