package middleware

import (
	"context"
	"errors"
	"time"

	"app-backend/setting"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/patcharp/golib/v2/crypto"
	"github.com/patcharp/golib/v2/server"
	"github.com/patcharp/golib/v2/util/httputil"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

type RequestGenerateToken struct {
	Username    string `json:"username"`
	FirstNameTh string `json:"first_name_th"`
	LastNameTh  string `json:"last_name_th"`
}

type AuthClaims struct {
	Id       uuid.UUID `json:"id,omitempty"`
	Subject  string    `json:"subject"`
	IssuedAt time.Time `json:"issued_at"`

	ExpiresAt float64 `json:"expires_at"`

	Username    string `json:"username"`
	FirstNameTh string `json:"first_name_th"`
	LastNameTh  string `json:"last_name_th"`
}

func AuthSystem(skipper *server.SkipperPath) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		if skipper != nil && skipper.Test(ctx) {
			return ctx.Next()
		}
		token := httputil.GetTokenFromCookie(ctx, "user_token")
		if token == "" {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":     fiber.StatusUnauthorized,
				"message":    "authorized token not found",
				"message_th": "ยืนยันตัวตนไม่สำเร็จ",
				"error":      "unauthorized",
			})
		}

		claims, err := VerifyJwt(token)
		if err != nil {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":     fiber.StatusUnauthorized,
				"message":    "invalid session",
				"message_th": "ยืนยันตัวตนไม่สำเร็จ",
				"error":      "unauthorized",
			})
		}

		// Check expires token
		if float64(time.Now().Unix()) > claims.ExpiresAt {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":     fiber.StatusUnauthorized,
				"message":    "JWT Expire",
				"message_th": "token หมดอายุ",
				"error":      "unauthorized",
			})
		}

		request := context.Background()
		authCtx := context.WithValue(request, "user_login_system", claims)
		ctx.SetUserContext(authCtx)
		return ctx.Next()

	}
}

func GenerateToken(id, subject string, model RequestGenerateToken, duration *time.Time) (string, error) {
	issAt := time.Now()
	claims := &jwt.MapClaims{
		"id":        id,
		"issued_at": issAt.Unix(),
		"subject":   subject,

		"expires_at": duration.Unix(),

		//model user
		"first_name_th": model.FirstNameTh,
		"last_name_th":  model.LastNameTh,
		"username":      model.Username,
	}

	return crypto.EncodeJWTAccessToken(claims, setting.GetCfg().PrivateKey)
}

func VerifyJwt(tokenString string) (*AuthClaims, error) {
	result := AuthClaims{}
	token, err := DecodeToken(tokenString)

	if err != nil {
		return nil, errors.New("invalid token claims")
	}

	jwtClaim := token.(jwt.MapClaims)
	result.Id = uuid.FromStringOrNil(jwtClaim["id"].(string))
	result.Subject = jwtClaim["subject"].(string)
	result.IssuedAt = time.Unix(int64(jwtClaim["issued_at"].(float64)), 0)

	result.Username = jwtClaim["username"].(string)
	result.LastNameTh = jwtClaim["last_name_th"].(string)
	result.FirstNameTh = jwtClaim["first_name_th"].(string)
	result.ExpiresAt = jwtClaim["expires_at"].(float64)

	return &result, nil
}

func DecodeToken(tokenString string) (jwt.Claims, error) {
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return &setting.GetCfg().PrivateKey.PublicKey, nil
	})
	if err != nil {
		logrus.Error("parse claims err :", err)
		return nil, err
	}
	if token.Valid {
		return token.Claims, nil
	}
	return nil, errors.New("decode token error")
}


func GetAuthUser(ctx *fiber.Ctx) (*AuthClaims, error) {
	s, ok := ctx.UserContext().Value("user_login_system").(*AuthClaims)

	if !ok {
		return nil, errors.New("unauthorized claims")
	}
	return s, nil
}
