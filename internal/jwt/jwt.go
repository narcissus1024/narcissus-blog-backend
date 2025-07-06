package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// todo 写入配置文件？
var (
	secret               = "gflksadh8y9saslk"
	accessTokenExpire    = 2                // h
	refreshTokenExpire   = 24 * 10          // h
	RefreshTokenInterval = 60 * 60 * 24 * 1 // ms
	DefaultIssuse        = "jwt"
)

type MyClaims struct {
	UserID int
	jwt.RegisteredClaims
}

func GenToken(issuer string, expireHour int, userId int) (string, error) {
	c := MyClaims{
		UserID: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    issuer,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expireHour) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	return token.SignedString([]byte(secret))
}

func GenAccessToken(issuer string, userId int) (string, error) {
	return GenToken(issuer, accessTokenExpire, userId)
}

func GenRefreshToken(issuer string, userId int) (string, error) {
	return GenToken(issuer, refreshTokenExpire, userId)
}

func ParseToken(tokenString string) (*MyClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("token error")
}

func IsExpireErr(err error) bool {
	return errors.Is(err, jwt.ErrTokenExpired)
}

func IsTokenExpire(tokenString string) (bool, error) {
	_, err := ParseToken(tokenString)
	if err != nil {
		if IsExpireErr(err) {
			return true, nil
		}
		return false, err
	}

	return false, nil
}

// 若token已经过期，返回0
func GetRemainExpireTime(tokenString string) (int64, error) {
	claims, parseTokenErr := ParseToken(tokenString)
	if parseTokenErr != nil {
		if IsExpireErr(parseTokenErr) {
			return 0, nil
		}
		return 0, parseTokenErr
	}
	expiresAt, getExpirationTimeErr := claims.GetExpirationTime()
	if getExpirationTimeErr != nil {
		return 0, getExpirationTimeErr
	}
	remainExpireTime := expiresAt.Unix() - time.Now().Unix()
	if remainExpireTime < 0 {
		remainExpireTime = 0
	}
	return remainExpireTime, nil
}

// func RefreshToken(aToken, rToken string, expireHour int) (newAToken, newRToken string, err error) {
// 	// refresh token无效直接返回
// 	if _, err = jwt.Parse(rToken, func(t *jwt.Token) (interface{}, error) {
// 		return []byte(secret), nil
// 	}); err != nil {
// 		return
// 	}

// 	// 从旧access token中解析出claims数据
// 	var claims MyClaims
// 	_, err = jwt.ParseWithClaims(aToken, &claims, func(t *jwt.Token) (interface{}, error) {
// 		return []byte(secret), nil
// 	})
// 	// 当access token是过期错误 并且 refresh token没有过期时就创建一个新的access token
// 	if errors.Is(err, jwt.ErrTokenExpired) {
// 		token, _ := GenToken(claims.Issuer, expireHour, claims.UserID)
// 		return token, "", nil
// 	}

// 	return
// }

func GetUserID(tokenString string) (int, error) {
	claims, parseTokenErr := ParseToken(tokenString)
	if parseTokenErr != nil {
		return -1, parseTokenErr
	}
	return claims.UserID, nil
}
