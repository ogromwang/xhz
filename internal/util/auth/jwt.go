package auth

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
	"time"
	"xiaohuazhu/internal/config"
	"xiaohuazhu/internal/model"
)

func GenerateToken(dto *model.AccountDTO) (string, error) {
	expireTime := time.Now().Add(time.Duration(config.AllConfig.Application.Auth.JwtExpireHour) * time.Hour)

	claims := model.JwtClaims{
		ID:       int64(dto.Id),
		Username: dto.Username,
		Icon:     dto.Icon,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    dto.Username,
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(config.AllConfig.Application.Auth.JwtSigned))
}

func ParseToken(token string) (*model.AccountDTO, error) {
	tokenClaims, err := jwt.ParseWithClaims(
		token,
		&model.JwtClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(config.AllConfig.Application.Auth.JwtSigned), nil
		})

	if err != nil {
		logrus.Errorf("[jwt|ParseToken] err: [%+v]", err)
		return nil, err
	}
	if tokenClaims != nil {
		if _, ok := tokenClaims.Claims.(*model.JwtClaims); ok && tokenClaims.Valid {
			jwtClaims := tokenClaims.Claims.(*model.JwtClaims)
			return &model.AccountDTO{
				Id:       uint(jwtClaims.ID),
				Username: jwtClaims.Username,
				Icon:     jwtClaims.Icon,
			}, nil
		}
	}

	return nil, err
}
