package jwt

import (
	"context"
	"errors"
	"github.com/Anniext/Arkitektur/code"
	"github.com/Anniext/Arkitektur/utils"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtToken *Jwt

type Jwt struct {
	SigningKey []byte
}

func NewJwt() *Jwt {
	if jwtToken == nil {
		jwtToken = &Jwt{}
	}

	return jwtToken
}

// GenerateJwtToken method    生成jwt密钥
func (j *Jwt) GenerateJwtToken(claimsMap jwt.MapClaims) (string, code.ErrCode) {
	if claimsMap == nil {
		claimsMap = make(jwt.MapClaims)
	}
	claimsMap["iat"] = time.Now().Unix()
	claimsMap["nbf"] = time.Now().Unix()
	if _, ok := claimsMap["exp"]; ok == false {
		claimsMap["exp"] = time.Now().Add(time.Minute * 15).Unix()
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claimsMap)
	tokenString, err := token.SignedString(j.getSecretKey())
	if err != nil {
		return "", code.ErrCodeJwtGenerateErr
	}

	return tokenString, code.NoErrCode
}

// ParseJwtToken method    解析jwt token
func (j *Jwt) ParseJwtToken(tokenString string) (map[string]interface{}, code.ErrCode) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return j.getSecretKey(), nil
	})
	if err != nil {
		var ve *jwt.ValidationError
		if errors.As(err, &ve) {
			if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, code.ErrCodeJwtTokenIsExpired
			} else if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, code.ErrCodeJwtNotEvenAToken
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, code.ErrCodeJwtTokenNotInvalid
			} else {
				return nil, code.ErrCodeJwtTokenNotActiveYet
			}
		}
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, code.NoErrCode
	}

	return nil, code.ErrCodeJwtTokenNotActiveYet
}

// loadSecretKey method    加载jwt的secretKey
func (j *Jwt) loadSecretKey() []byte {
	cnf := GetDefaultJwtConfig()

	return []byte(cnf.JwtSigningKey)
}

// getSecretKey method    获取jwt的secretKey
func (j *Jwt) getSecretKey() []byte {

	if len(j.SigningKey) == 0 {
		j.SigningKey = j.loadSecretKey()
	}

	return j.SigningKey
}

// GetTokenData function    获取token
func GetTokenData[T utils.MapSupportedTypes](ctx context.Context, key string) T {
	claimsMap, ok := ctx.Value("claims").(map[string]any)
	if !ok {
		var zero T
		return zero
	}

	return utils.GetMapSpecificValue[T](claimsMap, key)
}
