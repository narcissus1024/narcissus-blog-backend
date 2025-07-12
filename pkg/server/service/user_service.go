package service

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/narcissus1949/narcissus-blog/cmd/config"
	"github.com/narcissus1949/narcissus-blog/internal/database/cache"
	"github.com/narcissus1949/narcissus-blog/internal/encrypt"
	cerr "github.com/narcissus1949/narcissus-blog/internal/error"
	"github.com/narcissus1949/narcissus-blog/internal/jwt"
	"github.com/narcissus1949/narcissus-blog/internal/model"
	"github.com/narcissus1949/narcissus-blog/internal/utils"
	"github.com/narcissus1949/narcissus-blog/pkg/dto"
	"github.com/narcissus1949/narcissus-blog/pkg/server/dao"
	"github.com/narcissus1949/narcissus-blog/pkg/vo"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

var UserServiceInstance = new(userService)

type userService struct {
}

func (s *userService) SignIn(signIn dto.SignIn) error {
	if err := signIn.Validate(); err != nil {
		return cerr.NewParamError(err.Error())
	}

	existUser, queryErr := dao.UserDaoInstance.QueryByUsername(signIn.Username)
	if queryErr != nil {
		zap.L().Error("Failed to query user by username", zap.Error(queryErr), zap.String("username", signIn.Username))
		return queryErr
	}

	if existUser != nil {
		return cerr.New(cerr.ERROR_USER_ALREADY_EXIST)
	}

	// rsa解密
	decryptPasswd, decryptErr := encrypt.RSADecryptWithBase64([]byte(signIn.Password), config.Config.App.PrivateKeyDir)
	if decryptErr != nil {
		zap.L().Error("Failed to decrypt password", zap.Error(decryptErr), zap.String("username", signIn.Username))
		return decryptErr
	}

	hashedPassword, hashErr := bcrypt.GenerateFromPassword([]byte(decryptPasswd), bcrypt.DefaultCost)
	if hashErr != nil {
		zap.L().Error("Failed to hash password", zap.Error(hashErr), zap.String("username", signIn.Username))
		return hashErr
	}

	now := time.Now()
	user := model.User{
		Username:    signIn.Username,
		Nickname:    signIn.Nickname,
		Password:    string(hashedPassword),
		Email:       signIn.Email,
		PhoneNumber: signIn.PhoneNumber,
		AvatarPath:  signIn.AvatarPath,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := dao.UserDaoInstance.InsertUser(&user); err != nil {
		zap.L().Error("Failed to insert user", zap.Error(err), zap.String("username", signIn.Username))
		return err
	}
	return nil
}

func (s *userService) Login(loginRequest dto.LoginDto) (*vo.LoginVo, error) {
	// 用户验证
	user, queryErr := dao.UserDaoInstance.QueryByUsername(loginRequest.Username)
	if queryErr != nil {
		zap.L().Error("Failed to query user by username", zap.Error(queryErr), zap.String("username", loginRequest.Username))
		return nil, queryErr
	}

	if user == nil {
		return nil, cerr.New(cerr.ERROR_USER_NOT_EXIST)
	}

	// rsa解密
	decryptPasswd, decryptErr := encrypt.RSADecryptWithBase64([]byte(loginRequest.Password), config.Config.App.PrivateKeyDir)
	if decryptErr != nil {
		zap.L().Error("Failed to decrypt password", zap.Error(decryptErr), zap.String("username", loginRequest.Username))
		return nil, decryptErr
	}
	// 密码验证
	if compareErr := bcrypt.CompareHashAndPassword([]byte(user.Password), decryptPasswd); compareErr != nil {
		return nil, cerr.New(cerr.ERROR_USER_PASSWORD_WRONG)
	}

	// 生成token
	accessToken, genAccessTokenErr := jwt.GenAccessToken(jwt.DefaultIssuse, user.ID)
	if genAccessTokenErr != nil {
		zap.L().Error("Failed to generate access token", zap.Error(genAccessTokenErr), zap.String("username", loginRequest.Username))
		return nil, genAccessTokenErr
	}
	refreshToken, genRefreshTokenErr := jwt.GenRefreshToken(jwt.DefaultIssuse, user.ID)
	if genRefreshTokenErr != nil {
		zap.L().Error("Failed to generate refresh token", zap.Error(genRefreshTokenErr), zap.String("username", loginRequest.Username))
		return nil, genRefreshTokenErr
	}
	resp := &vo.LoginVo{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return resp, nil
}

func (s *userService) Logout(ctx *gin.Context, logoutRequest dto.LogoutDto) error {
	redisClient := cache.Client

	accessTokenRemainExpireTime, getAccessTokenRemainExpireTimeErr := jwt.GetRemainExpireTime(logoutRequest.AccessToken)
	if getAccessTokenRemainExpireTimeErr != nil {
		zap.L().Error("Failed to get access token remain expire time", zap.Error(getAccessTokenRemainExpireTimeErr))
		return cerr.New(cerr.ERROR_USER_TOKEN_INVALIDE)
	}
	userIDFromAccessToken, getUserIDFromAccessTokenErr := jwt.GetUserID(logoutRequest.AccessToken)
	if getUserIDFromAccessTokenErr != nil && !jwt.IsExpireErr(getUserIDFromAccessTokenErr) {
		zap.L().Error("Failed to get user id from access token", zap.Error(getUserIDFromAccessTokenErr))
		return cerr.New(cerr.ERROR_USER_TOKEN_INVALIDE)
	}

	refreshTokenRemainExpireTime, getRefreshTokenRemainExpireTimeErr := jwt.GetRemainExpireTime(logoutRequest.RefreshToken)
	if getRefreshTokenRemainExpireTimeErr != nil {
		zap.L().Error("Failed to get refresh token remain expire time", zap.Error(getRefreshTokenRemainExpireTimeErr))
		return cerr.New(cerr.ERROR_USER_TOKEN_INVALIDE)
	}
	userIDFromRefreshToken, getUserIDFromRefreshTokenErr := jwt.GetUserID(logoutRequest.RefreshToken)
	if getUserIDFromRefreshTokenErr != nil && !jwt.IsExpireErr(getUserIDFromRefreshTokenErr) {
		zap.L().Error("Failed to get user id from refresh token", zap.Error(getUserIDFromRefreshTokenErr))
		return cerr.New(cerr.ERROR_USER_TOKEN_INVALIDE)
	}

	if userIDFromAccessToken != userIDFromRefreshToken {
		zap.L().Error("User id from access token and refresh token not equal")
		return cerr.New(cerr.ERROR_USER_TOKEN_INVALIDE)
	}

	if accessTokenRemainExpireTime > 0 {
		status := redisClient.Set(ctx, utils.ACCESS_TOKEN_BLACKLIST+logoutRequest.AccessToken,
			userIDFromAccessToken,
			time.Duration(accessTokenRemainExpireTime)*time.Second)
		if status.Err() != nil {
			zap.L().Error("Failed to set access token blacklist", zap.Error(status.Err()))
			return status.Err()
		}
	}

	if refreshTokenRemainExpireTime > 0 {
		status := redisClient.Set(ctx, utils.REFRESH_TOKEN_BLACKLIST+logoutRequest.RefreshToken,
			userIDFromRefreshToken,
			time.Duration(refreshTokenRemainExpireTime)*time.Second)
		if status.Err() != nil {
			zap.L().Error("Failed to set refresh token blacklist", zap.Error(status.Err()))
			return status.Err()
		}
	}

	return nil
}

func (s *userService) RefreshToken(ctx *gin.Context, refreshTokenRequest dto.RefreshTokenDto) (*vo.LoginVo, error) {
	claims, parseTokenErr := jwt.ParseToken(refreshTokenRequest.RefreshToken)
	if parseTokenErr != nil {
		zap.L().Error("Failed to parse refresh token", zap.Error(parseTokenErr))
		return nil, parseTokenErr
	}
	expireTime, getExpireTimeErr := claims.GetExpirationTime()
	if getExpireTimeErr != nil {
		zap.L().Error("Failed to get expire time from refresh token", zap.Error(getExpireTimeErr))
		return nil, getExpireTimeErr
	}

	resp := &vo.LoginVo{}
	// 刷新token
	accessToken, genAccessTokenErr := jwt.GenAccessToken(jwt.DefaultIssuse, claims.UserID)
	if genAccessTokenErr != nil {
		zap.L().Error("Failed to generate access token", zap.Error(genAccessTokenErr))
		return nil, genAccessTokenErr
	}
	resp.AccessToken = accessToken

	// 若refresh token距离过期时间小于jwt.RefreshTokenInterval，则同时刷新refresh token
	remainExpire := expireTime.Unix() - time.Now().Unix()
	if remainExpire <= int64(jwt.RefreshTokenInterval) {
		if remainExpire > 0 {
			_, setBlacklistErr := cache.Client.Set(ctx, utils.REFRESH_TOKEN_BLACKLIST+refreshTokenRequest.RefreshToken,
				claims.UserID,
				time.Duration(expireTime.Unix())*time.Second).Result()
			if setBlacklistErr != nil {
				zap.L().Error("Failed to set refresh token blacklist", zap.Error(setBlacklistErr))
				return nil, setBlacklistErr
			}
		}

		refreshToken, genRefreshTokenErr := jwt.GenRefreshToken(jwt.DefaultIssuse, claims.UserID)
		if genRefreshTokenErr != nil {
			zap.L().Error("Failed to generate refresh token", zap.Error(genRefreshTokenErr))
			return nil, genRefreshTokenErr
		}
		resp.RefreshToken = refreshToken
	}

	return resp, nil
}
