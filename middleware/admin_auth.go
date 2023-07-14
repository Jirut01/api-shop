package middleware

import (
	"app-backend/setting"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/patcharp/golib/v2/server"
)

func AdminAuth(skipper *server.SkipperPath) fiber.Handler {
	return basicauth.New(basicauth.Config{
		Users: map[string]string{
			setting.GetCfg().AdminUsername: setting.GetCfg().AdminPassword,
		},
	})

}